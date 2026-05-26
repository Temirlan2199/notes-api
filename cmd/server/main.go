package main

import (
	"log"
	"net/http"

	"notes-api/internal/config"
	"notes-api/internal/handler"
	"notes-api/internal/repository"
	"notes-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

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

	log.Printf("Сервер запускается на http://localhost%s", cfg.Addr())
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
