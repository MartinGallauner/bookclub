package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserRepository interface {
	GetUser(userId int) string //todo which case is expected? snake case vs camel case?
}

type BookRepository interface {
	GetBook(isbn string) Book
}

type BookclubServer struct {
	bookRepository BookRepository
}

func StartServer(cfg *config) {
	const port = "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/collections", cfg.handlerAddBook)
	mux.HandleFunc("GET /api/books/{isbn}", cfg.handlerGetBookByISBN)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting bookclub on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func (srv *BookclubServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := AddBookRequest{}
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	fmt.Fprint(w, srv.bookRepository.GetBook("1234"))
}
