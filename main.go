package main

import (
	"github.com/joho/godotenv"
	internal "github.com/martingallauner/bookclub/internal"
	"github.com/martingallauner/bookclub/internal/auth"
	. "github.com/martingallauner/bookclub/internal/client"
	repository "github.com/martingallauner/bookclub/internal/repository"
	server "github.com/martingallauner/bookclub/internal/server"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	auth.NewAuth()

	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DBNAME")
	port := os.Getenv("POSTGRES_PORT")
	sslmode := os.Getenv("POSTGRES_SSLMODE")
	timezone := os.Getenv("TIMEZONE")

	db, err := internal.SetupDatabase(host, user, password, dbname, port, sslmode, timezone)
	if err != nil {
		log.Fatal(err)
	}
	client := NewClient(5 * time.Second)

	server := server.NewBookclubServer(
		client,
		&repository.PostgresBookRepository{Database: db},
		&repository.PostgresUserRepository{Database: db},
		&repository.PostgresLinkRepository{Database: db},
		&internal.GothicAuthService{},
		&internal.JwtServiceImpl{})
	err = server.StartServer(server)
	if err != nil {
		log.Fatal(err)
	}
}
