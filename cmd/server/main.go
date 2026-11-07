package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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

func main() {
	if err := run(); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	slog.Info("load config success")

	logger := logx.New(cfg.IsProduction())
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ---- Postgres ----
	gormDB, err := db.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}
	defer sqlDB.Close()
	slog.Info("connect postgres success")

	// ---- Redis ----
	rdb, err := cache.NewRedis(ctx, cfg.RedisAddr)
	if err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}
	defer rdb.Close()
	slog.Info("connect redis success")

	// ---- Routes  ----
	tokenMgr := token.NewManager(cfg.JWTSecret, auth.AccessTTL)

	authRepo := auth.NewRepo(gormDB)
	userRepo := user.NewRepo(gormDB)
	authSvc := auth.NewService(authRepo, userRepo, tokenMgr)

	r := api.NewRouter(logger, tokenMgr, api.Handlers{
		Auth: auth.NewHandler(authSvc, cfg.IsProduction()),
		User: user.NewHandler(user.NewService(userRepo)),
		Post: post.NewHandler(post.NewService(post.NewRepo(gormDB))),
	})

	// ---- Start Server ----
	slog.Info("server starting", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
