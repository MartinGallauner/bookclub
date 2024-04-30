package main

import (
	"log"
	"net/http"
)

type UserRepository interface {
	Get(id uint) (User, error)
	Save(user *User) error

	SearchBook(isbn string) ([]User, error)
}

type BookRepository interface {
	GetBook(isbn string) (Book, error)
	Save(book Book) error
}

type LinkRepository interface {
	//Returns specific Link between two users
	Get(senderId uint, receiverId uint) (Link, error)

	//Returns all links concerned with the user
	GetById(userId string) ([]Link, error)

	//Returns all links concerned with the user that are accepted
	GetAcceptedById(userId uint) ([]Link, error)

	//Persists link.
	Save(link *Link) error
}

func StartServer(cfg *BookclubServer) {
	log.Print("Starting bookclub on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", cfg.Handler))
	//todo not sure if I should return an error here
}

type BookclubServer struct {
	Client         Client
	BookRepository BookRepository
	UserRepository UserRepository
	LinkRepository LinkRepository
	http.Handler
}

func NewBookclubServer(client Client, repository BookRepository, userRepository UserRepository, linkRepository LinkRepository) *BookclubServer {
	s := new(BookclubServer)
	s.BookRepository = repository
	s.UserRepository = userRepository
	s.LinkRepository = linkRepository
	s.Client = client

	router := http.NewServeMux()
	router.Handle("/api/search", http.HandlerFunc(s.handlerSearch))
	router.Handle("/api/collections", http.HandlerFunc(s.handlerAddBook))
	router.Handle("/api/users", http.HandlerFunc(s.handlerCreateUser))
	router.Handle("/api/links/{id}", http.HandlerFunc(s.handlerGetLinks))
	router.Handle("/api/links", http.HandlerFunc(s.handlerCreateLink))
	router.Handle("/auth/{provider}/callback", http.HandlerFunc(s.handlerCallback))
	router.Handle("/auth/{provider}/logout", http.HandlerFunc(s.handlerLogout))
	router.Handle("/auth/{provider}", http.HandlerFunc(s.handlerLogin))

	s.Handler = router
	return s
}
