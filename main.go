package main

import (
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type config struct {
	Client         Client
	Database       *gorm.DB
	BookRepository BookRepository
	UserRepository UserRepository
}

type PostgresBookRepository struct {
	Database *gorm.DB
}

func (r *PostgresBookRepository) GetBook(isbn string) Book {
	var book Book
	r.Database.Table("books").Find(&book, isbn)
	return book
}

type PostgresUserRepository struct {
	Database *gorm.DB
}

func (g *PostgresUserRepository) Get(id int) User {
	var user User
	g.Database.Table("users").Find(&user, id)
	return user
}

func (g *PostgresUserRepository) Save(user User) error {
	err := g.Database.Table("users").Save(&user).Error
	return err
}

func main() {
	db, err := SetupDatabase("host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Vienna")
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(5 * time.Second)
	cfg := &config{
		Client:         client,
		Database:       db,
		BookRepository: &PostgresBookRepository{Database: db},
		UserRepository: &PostgresUserRepository{Database: db},
	}
	handler := http.HandlerFunc(cfg.handlerAddBook)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
