package main

import (
	"github.com/martingallauner/bookclub/internal"
	"github.com/martingallauner/bookclub/internal/auth"
	"log"
	"time"
)

func main() {

	auth.NewAuth()

	db, err := internal.SetupDatabase("host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Vienna")
	if err != nil {
		log.Fatal(err)
	}

	client := internal.NewClient(5 * time.Second)

	server := internal.NewBookclubServer(client, &internal.PostgresBookRepository{Database: db}, &internal.PostgresUserRepository{Database: db}, &internal.PostgresLinkRepository{Database: db})
	internal.StartServer(server)
}
