package main

import (
	"log"
	"net/http"
)

type UserRepository interface {
	Get(id int) (User, error)
	Save(user User) error

	Find(isbn string) ([]User, error)
}

type BookRepository interface {
	GetBook(isbn string) Book
	Save(book Book) error
}

func StartServer(cfg *BookclubServer) {
	log.Print("Starting bookclub on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", cfg.Handler))
	//todo not sure if I should return an error here
}

func NewBookclubServer(client Client, repository BookRepository, userRepository UserRepository) *BookclubServer {
	s := new(BookclubServer)
	s.BookRepository = repository
	s.UserRepository = userRepository
	s.Client = client

	router := http.NewServeMux()
	router.Handle("/api/search/{isbn}", http.HandlerFunc(s.handlerSearch))
	router.Handle("/api/collections", http.HandlerFunc(s.handlerAddBook))
	//router.Handle("/api/books/{isbn}", http.HandlerFunc(s.addToCollectionHandler))

	s.Handler = router
	return s
}
