package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clary-work01/hot-post-api/internal/api"
	"github.com/clary-work01/hot-post-api/internal/domain/auth"
	"github.com/clary-work01/hot-post-api/internal/domain/post"
	"github.com/clary-work01/hot-post-api/internal/domain/user"
	"github.com/clary-work01/hot-post-api/internal/infra/cache"
	"github.com/clary-work01/hot-post-api/internal/infra/db"
	"github.com/clary-work01/hot-post-api/internal/platform/config"
	"github.com/clary-work01/hot-post-api/internal/platform/logx"
	"github.com/clary-work01/hot-post-api/internal/platform/token"
)

const (
	startupTimeout  = 5 * time.Second
	shutdownTimeout = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// 監聽關機訊號
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	slog.Info("load config success")

	logger := logx.New(cfg.IsProduction())
	slog.SetDefault(logger)

	startCtx, cancelStart := context.WithTimeout(ctx, 5*time.Second)
	defer cancelStart()

	// ---- Postgres ----
	gormDB, err := db.NewPostgres(startCtx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}
	slog.Info("connect postgres success")

	defer func() {
		sqlDB.Close()
		slog.Info("postgres closed")
	}()

	// ---- Redis ----
	rdb, err := cache.NewRedis(startCtx, cfg.RedisAddr)
	if err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}
	slog.Info("connect redis success")

	defer func() {
		rdb.Close()
		slog.Info("redis closed")
	}()

	// ---- Routes  ----
	tokenMgr := token.NewManager(cfg.JWTSecret, auth.AccessTTL)

	authRepo := auth.NewRepo(gormDB)
	userRepo := user.NewRepo(gormDB)
	postRepo := post.NewRepo(gormDB)

	r := api.NewRouter(logger, tokenMgr, api.Handlers{
		Auth: auth.NewHandler(auth.NewService(authRepo, userRepo, tokenMgr), cfg.IsProduction()),
		User: user.NewHandler(user.NewService(userRepo)),
		Post: post.NewHandler(post.NewService(postRepo)),
	})

	// ---- Start Server ----
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
		// 防 Slowloris 攻擊（故意很慢地送 header 佔住連線）
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		// handler 加寫回應的上限，要大於你最慢的正常請求
		WriteTimeout: 15 * time.Second,
		// keep-alive 連線閒置多久關掉
		IdleTimeout: 60 * time.Second,
	}

	// 在 goroutine 裡啟動 server
	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// 等待：訊號 or 啟動失敗
	select {
	case err := <-errCh:
		return fmt.Errorf("listen and serve: %w", err)
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	// 恢復預設行為：再按一次 Ctrl+C 會直接強制結束
	stop()

	// 優雅關閉
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	slog.Info("server stopped gracefully")

	// return 之後 defer 依序關閉 Redis → Postgres
	return nil
}
