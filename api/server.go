package api

import (
	"crawler/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StartServer() {
	log.Println("Starting server on :8080")

	r := chi.NewRouter()

	r.Post("/api/start-crawler", handlers.CrawlerHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
