package main

import (
	"log"
	"net/http"
	"time"
)

type config struct {
	Client         Client
	BookRepository BookRepository
	UserRepository UserRepository
}

func main() {
	db, err := SetupDatabase("host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Vienna")
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(5 * time.Second)
	cfg := &config{
		Client:         client,
		BookRepository: &PostgresBookRepository{Database: db},
		UserRepository: &PostgresUserRepository{Database: db},
	}
	handler := http.HandlerFunc(cfg.handlerAddBook)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
