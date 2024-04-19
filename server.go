package main

import (
	"net/http"
)

type UserRepository interface {
	Get(id int) (User, error) //todo which case is expected? snake case vs camel case?
	Save(user User) error
}

type BookRepository interface {
	GetBook(isbn string) Book
	Save(book Book) error
}

func StartServer(cfg *config) error {
	const port = "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/collections", cfg.handlerAddBook)
	mux.HandleFunc("GET /api/books/{isbn}", cfg.handlerGetBookByISBN)

	return nil

}
