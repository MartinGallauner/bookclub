package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/martingallauner/bookclub/docs"
	internal "github.com/martingallauner/bookclub/internal"
	client "github.com/martingallauner/bookclub/internal/client"
	"github.com/martingallauner/bookclub/internal/collections"
	"github.com/martingallauner/bookclub/internal/users"
	repository "github.com/martingallauner/bookclub/internal/repository"
	httpSwagger "github.com/swaggo/http-swagger"
)

// StartServer starts the server :)
func StartServer(cfg *BookclubServer) error {
	log.Print("Starting bookclub on port: 8080") //TODO:: make parameters configurable

	s := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		Handler:           cfg.Handler,
	}
	return s.ListenAndServe()
}

type BookclubServer struct {
	Client         client.Client
	BookRepository repository.BookRepository
	UserRepository repository.UserRepository
	LinkRepository repository.LinkRepository
	AuthService    internal.AuthService
	JwtService     internal.JwtService
	http.Handler
	CollectionService *collections.Service
	UsersService *users.Service

}

func New(client client.Client, 
	bookRepository repository.BookRepository, 
	userRepository repository.UserRepository, 
	linkRepository repository.LinkRepository, 
	authService internal.AuthService, 
	jwtService internal.JwtService, 
	collectionService *collections.Service,
	usersService *users.Service) *BookclubServer {
	s := new(BookclubServer)
	s.BookRepository = bookRepository
	s.UserRepository = userRepository
	s.LinkRepository = linkRepository
	s.Client = client
	s.AuthService = authService
	s.JwtService = jwtService
	s.CollectionService = collectionService
	s.UsersService = usersService
	router := http.NewServeMux() //TODO: add jwtMiddleware to all concerned handler
	router.Handle("/api/search", http.HandlerFunc(s.handlerSearch))
	router.Handle("/api/collections", http.HandlerFunc(s.handlerAddBook))
	//router.Handle("/api/collections", http.HandlerFunc(jwtMiddleware(s.handlerAddBook)))
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

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)
		if tokenString == "" {
			RespondWithError(w, http.StatusUnauthorized, "Missing JWT")
			return
		}
		claims, err := validateToken(tokenString)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid JWT")
			return
		}
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func extractToken(r *http.Request) string {
	authorizationHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
	return tokenString
}

func validateToken(tokenString string) (jwt.RegisteredClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}
	if claims.Issuer != "bookclub_access" {
		return jwt.RegisteredClaims{}, errors.New(fmt.Sprintf("Invalid issuer %s", claims.Issuer))
	}
	if !token.Valid {
		return jwt.RegisteredClaims{}, errors.New("Invalid token")
	}
	return claims, nil
}
