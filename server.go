package main

import (
	"log"
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

func StartServer(cfg *BookclubServer) {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/collections", cfg.handlerAddBook)
	router.HandleFunc("GET /api/books/{isbn}", cfg.handlerGetBookByISBN)

	log.Print("Starting bookclub on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
	//todo not sure if I should return an error here
}
