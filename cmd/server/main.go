package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/clary-work01/hot-post-api/internal/config"
	"github.com/clary-work01/hot-post-api/internal/infra/cache"
	"github.com/clary-work01/hot-post-api/internal/infra/db"
)

func main() {
	cfg := config.Load()

	initCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := db.NewPostgres(initCtx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect postgres failed", "error", err)
		os.Exit(1)
	}
	slog.Info("connect postgres success")
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	rdb, err := cache.NewRedis(initCtx, cfg.RedisAddr)
	if err != nil {
		slog.Error("connect redis failed", "error", err)
		os.Exit(1)
	}
	slog.Info("connect redis success")
	defer rdb.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

	slog.Info("server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
