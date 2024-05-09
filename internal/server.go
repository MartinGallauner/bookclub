package internal

import (
	_ "github.com/martingallauner/bookclub/docs"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

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
	AuthService    AuthService
	JwtService     JwtService
	http.Handler
}

type UseCases interface {
	AddBookToCollection(isbn string, userId uint) (Book, error)
	SearchBookInNetwork(userId uint, isbn string) ([]User, error)
	CreateUser(name string) (User, error)
	LinkUsers(senderId uint, receiverId uint) (Link, error)
}

func NewBookclubServer(client Client, repository BookRepository, userRepository UserRepository, linkRepository LinkRepository, authService AuthService, jwtService JwtService) *BookclubServer {
	s := new(BookclubServer)
	s.BookRepository = repository
	s.UserRepository = userRepository
	s.LinkRepository = linkRepository
	s.Client = client
	s.AuthService = authService
	s.JwtService = jwtService
	router := http.NewServeMux()
	router.Handle("/api/search", http.HandlerFunc(s.handlerSearch))
	router.Handle("/api/collections", http.HandlerFunc(s.handlerAddBook))
	router.Handle("/api/books/{isbn}", http.HandlerFunc(s.handlerGetBookByISBN))
	router.Handle("/api/users", http.HandlerFunc(s.handlerCreateUser))
	router.Handle("/api/links/{id}", http.HandlerFunc(s.handlerGetLinks))
	router.Handle("/api/links", http.HandlerFunc(s.handlerCreateLink))
	router.Handle("/api/auth/{provider}/callback", http.HandlerFunc(s.handlerCallback))
	router.Handle("/api/auth/{provider}/logout", http.HandlerFunc(s.handlerLogout))
	router.Handle("/api/auth/{provider}", http.HandlerFunc(s.handlerLogin))
	router.Handle("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
	router.Handle("/api/healthz", http.HandlerFunc(handlerReadiness))

	s.Handler = router
	return s
}
