package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/clary-work01/hot-post-api/internal/api"
	"github.com/clary-work01/hot-post-api/internal/config"
	"github.com/clary-work01/hot-post-api/internal/domain/post"
	"github.com/clary-work01/hot-post-api/internal/infra/cache"
	"github.com/clary-work01/hot-post-api/internal/infra/db"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ---- Postgres ----
	gormDB, err := db.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect postgres failed: %w", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %w", err)
	}
	defer sqlDB.Close()
	slog.Info("connect postgres success")

	// ---- Redis ----
	rdb, err := cache.NewRedis(ctx, cfg.RedisAddr)
	if err != nil {
		return fmt.Errorf("connect redis failed: %w", err)
	}
	defer rdb.Close()
	slog.Info("connect redis success")

	// ---- Routes  ----
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

	postSvc := post.NewService(post.NewRepo(gormDB))

	postHandler := api.NewPostHandler(postSvc)
	postHandler.RegisterRoutes(mux)

	// ---- Start Server ----
	slog.Info("server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}
