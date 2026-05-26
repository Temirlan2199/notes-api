package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"notes-api/internal/config"
	"notes-api/internal/handler"
	"notes-api/internal/repository"
	"notes-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		// Конфиг не загрузился — slog ещё нет, используем stderr напрямую
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	logger := newLogger(cfg)
	slog.SetDefault(logger)

	repo := repository.NewMemoryRepository()
	noteSvc := service.NewNoteService(repo)
	noteHandler := handler.NewNoteHandler(noteSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /notes", noteHandler.Create)
	mux.HandleFunc("GET /notes", noteHandler.List)
	mux.HandleFunc("GET /notes/count", noteHandler.Count)
	mux.HandleFunc("GET /notes/{id}", noteHandler.GetByID)
	mux.HandleFunc("PUT /notes/{id}", noteHandler.Update)
	mux.HandleFunc("DELETE /notes/{id}", noteHandler.Delete)

	server := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	logger.Info("server starting",
		"addr", cfg.Addr(),
		"read_timeout", cfg.ReadTimeout,
		"write_timeout", cfg.WriteTimeout,
		"log_level", cfg.LogLevel.String(),
		"log_format", cfg.LogFormat,
	)

	if err := server.ListenAndServe(); err != nil {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}

func newLogger(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: cfg.LogLevel}

	var handler slog.Handler
	if cfg.LogFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
