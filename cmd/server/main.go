package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	slog.Info("server starting", "port", ":8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
