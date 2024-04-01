package main

import (
	"log"
	"net/http"
	"time"
)

type config struct {
	Client Client
}

func main() {
	client := NewClient(5 * time.Second)
	cfg := &config{
		Client: client,
	}

	const port = "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/collection/{isbn}", cfg.handlerAddBook)
	mux.HandleFunc("GET /api/books/{isbn}", cfg.handlerGetBookByISBN)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Starting bookclub on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
