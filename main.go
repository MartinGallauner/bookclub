package main

import (
	//"github.com/joho/godotenv"
	"log"
	"time"

	internal "github.com/martingallauner/bookclub/internal"
	"github.com/martingallauner/bookclub/internal/auth"
	"github.com/martingallauner/bookclub/internal/client"
	"github.com/martingallauner/bookclub/internal/collections"
	repository "github.com/martingallauner/bookclub/internal/repository"
	bcServer "github.com/martingallauner/bookclub/internal/server"
	"github.com/martingallauner/bookclub/internal/users"
)

func main() {
	/* 	err := godotenv.Load()
	   	if err != nil {
	   		log.Fatal("Error loading .env file")
	   	} */

	auth.NewAuth()

	dbConfig, err := internal.ReadDatabaseConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := internal.SetupDatabase(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	client := client.NewClient(5 * time.Second)

	bookRepository := &repository.PostgresBookRepository{Database: db}
	userRepository := &repository.PostgresUserRepository{Database: db}
	linkRepository := &repository.PostgresLinkRepository{Database: db}

	collectionService := collections.New(userRepository, bookRepository, linkRepository, client)
	userService := users.New(userRepository, bookRepository, linkRepository, client)

	server := bcServer.New(
		client,
		bookRepository,
		userRepository,
		linkRepository,
		&internal.GothicAuthService{},
		&internal.JwtServiceImpl{},
		collectionService,
		userService)
	err = bcServer.StartServer(server)
	if err != nil {
		log.Fatal(err)
	}
}
