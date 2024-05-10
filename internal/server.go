package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/martingallauner/bookclub/docs"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"strings"
)

func StartServer(cfg *BookclubServer) error {
	log.Print("Starting bookclub on port: 8080")
	return http.ListenAndServe(":8080", cfg.Handler)
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

func NewBookclubServer(client Client, repository BookRepository, userRepository UserRepository, linkRepository LinkRepository, authService AuthService, jwtService JwtService) *BookclubServer {
	s := new(BookclubServer)
	s.BookRepository = repository
	s.UserRepository = userRepository
	s.LinkRepository = linkRepository
	s.Client = client
	s.AuthService = authService
	s.JwtService = jwtService
	router := http.NewServeMux() //todo add jwtMiddleware to all concerned handler
	router.Handle("/api/search", http.HandlerFunc(s.handlerSearch))
	router.Handle("/api/collections", http.HandlerFunc(jwtMiddleware(s.handlerAddBook)))
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
		tokenString := extractToken(r) //todo extract token
		if tokenString == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing JWT")
			return
		}
		claims, err := validateToken(tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid JWT")
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
	jwtSecret := os.Getenv("JWT_SECRET") //todo load from .env
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
