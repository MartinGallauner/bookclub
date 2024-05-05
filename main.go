package main

import (
	"github.com/joho/godotenv"
	"github.com/martingallauner/bookclub/internal"
	"github.com/martingallauner/bookclub/internal/auth"
	"log"
	"os"
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
	//client := internal.NewClient(5 * time.Second)

	internal.StartFreshServer(&internal.PostgresBookRepository{Database: db}, &internal.PostgresUserRepository{Database: db}, &internal.PostgresLinkRepository{Database: db})

	//server := internal.NewBookclubServer(client, &internal.PostgresBookRepository{Database: db}, &internal.PostgresUserRepository{Database: db}, &internal.PostgresLinkRepository{Database: db})
	//internal.StartServer(server)
}
