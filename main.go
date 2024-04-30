package main

import (
	"github.com/martingallauner/bookclub/internal/auth"
	"log"
	"time"
)

func main() {

	auth.NewAuth()

	db, err := SetupDatabase("host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Vienna")
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(5 * time.Second)

	server := NewBookclubServer(client, &PostgresBookRepository{Database: db}, &PostgresUserRepository{Database: db}, &PostgresLinkRepository{Database: db})
	StartServer(server)
}
