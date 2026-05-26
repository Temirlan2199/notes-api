package main

import (
	"log/slog"
	"net/http"
	"os"

	"notes-api/internal/config"
	"notes-api/internal/handler"
	"notes-api/internal/middleware"
	"notes-api/internal/repository"
	"notes-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		os.Stderr.WriteString("config error: " + err.Error() + "\n")
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

	// Оборачиваем mux в цепочку middleware
	// Порядок: Recovery (внешний) → Logging → mux
	// Recovery должен быть внешним, чтобы ловить паники из Logging тоже
	var rootHandler http.Handler = mux
	rootHandler = middleware.Logging(rootHandler)
	rootHandler = middleware.Recovery(rootHandler)

	server := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      rootHandler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	logger.Info("server starting",
		"addr", cfg.Addr(),
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
