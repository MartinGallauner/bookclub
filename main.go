package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type config struct {
	Client   Client
	Database *gorm.DB
}

func main() {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Vienna"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&User{}, &Book{})
	//err = db.SetupJoinTable(&User{}, "Books", &UserBooks{})
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(5 * time.Second)
	cfg := &config{
		Client:   client,
		Database: db,
	}

	const port = "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/collection/{isbn}", cfg.handlerAddBook)
	mux.HandleFunc("GET /api/books/{isbn}", cfg.handlerGetBookByISBN)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Starting bookclub on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
