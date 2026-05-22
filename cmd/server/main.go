package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	store := NewStore()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /notes", makeCreateNoteHandler(store))
	mux.HandleFunc("GET /notes", makeListNotesHandler(store))
	mux.HandleFunc("GET /notes/{id}", makeGetNoteHandler(store))
	mux.HandleFunc("PUT /notes/{id}", makeUpdateNoteHandler(store))
	mux.HandleFunc("DELETE /notes/{id}", makeDeleteNoteHandler(store))
	mux.HandleFunc("GET /notes/count", makeCreateNoteHandler(store))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Сервер запускается на http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
