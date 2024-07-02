package main

import (
	"github.com/joho/godotenv"
	internal "github.com/martingallauner/bookclub/internal"
	"github.com/martingallauner/bookclub/internal/auth"
	. "github.com/martingallauner/bookclub/internal/client"
	repository "github.com/martingallauner/bookclub/internal/repository"
	bcServer "github.com/martingallauner/bookclub/internal/server"
	"log"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	auth.NewAuth()


	dbConfig, err := internal.ReadDatabaseConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := internal.SetupDatabase(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	client := NewClient(5 * time.Second)

	server := bcServer.NewBookclubServer(
		client,
		&repository.PostgresBookRepository{Database: db},
		&repository.PostgresUserRepository{Database: db},
		&repository.PostgresLinkRepository{Database: db},
		&internal.GothicAuthService{},
		&internal.JwtServiceImpl{})
	err = bcServer.StartServer(server)
	if err != nil {
		log.Fatal(err)
	}
}
