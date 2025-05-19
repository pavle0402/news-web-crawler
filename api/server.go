package api

import (
	"log"
	"net/http"
	"crawler/handlers"
)


func StartServer() {
	log.Println("Starting server on :8080")

	http.HandleFunc("/api/start-crawler", handlers.CrawlerHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}