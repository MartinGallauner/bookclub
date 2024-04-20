package main

import (
	"log"
	"net/http"
	"time"
)

type BookclubServer struct {
	Client         Client
	BookRepository BookRepository
	UserRepository UserRepository
	http.Handler
}

func NewBookclubServer(client Client, repository BookRepository, userRepository UserRepository) *BookclubServer {
	s := new(BookclubServer)
	s.BookRepository = repository
	s.UserRepository = userRepository

	router := http.NewServeMux()
	router.Handle("/api/collections", http.HandlerFunc(s.handlerAddBook))
	//router.Handle("/api/books/{isbn}", http.HandlerFunc(s.addToCollectionHandler))

	s.Handler = router
	return s
}

func main() {
	db, err := SetupDatabase("host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Vienna")
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(5 * time.Second)

	server := NewBookclubServer(client, &PostgresBookRepository{Database: db}, &PostgresUserRepository{Database: db})
	StartServer(server)
}
