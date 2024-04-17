package main

import (
	"net/http"
)

type UserRepository interface {
	GetUser(userId int) string //todo which case is expected? snake case vs camel case?
}

type BookRepository interface {
	GetBook(isbn string) Book
}

func StartServer(cfg *config) error {
	const port = "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/collections", cfg.handlerAddBook)
	mux.HandleFunc("GET /api/books/{isbn}", cfg.handlerGetBookByISBN)

	return nil

}
