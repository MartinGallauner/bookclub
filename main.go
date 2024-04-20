package main

import (
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type config struct {
	Client         Client
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

func (r *PostgresBookRepository) Save(book Book) error {
	err := r.Database.Table("books").Save(&book).Error
	return err
}

type PostgresUserRepository struct {
	Database *gorm.DB
}

func (r *PostgresUserRepository) Get(id int) (User, error) {
	var user User
	err := r.Database.Table("users").First(&user, id).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) Save(user User) error {
	err := r.Database.Table("users").Save(&user).Error
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
		BookRepository: &PostgresBookRepository{Database: db},
		UserRepository: &PostgresUserRepository{Database: db},
	}
	handler := http.HandlerFunc(cfg.handlerAddBook)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
