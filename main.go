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

func main() {
	db, err := SetupDatabase("host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Vienna")
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(5 * time.Second)

	server := NewBookclubServer(client, &PostgresBookRepository{Database: db}, &PostgresUserRepository{Database: db})
	StartServer(server)
}
