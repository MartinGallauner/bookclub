package main

import (
	"gorm.io/driver/postgres"
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

type GormBookRepository struct { //todo I'm not sure about the name
	Database *gorm.DB
}

func (r *GormBookRepository) GetBook(isbn string) Book {
	var book Book
	r.Database.Table("books").Find(&book, isbn)
	return book
}

type GormUserRepository struct {
	Database *gorm.DB
}

func (g *GormUserRepository) Get(id int) User {
	var user User
	g.Database.Table("users").Find(&user, id)
	return user
}

func main() {
	db := SetupDatabase()

	client := NewClient(5 * time.Second)
	cfg := &config{
		Client:         client,
		Database:       db,
		BookRepository: &GormBookRepository{Database: db},
		UserRepository: &GormUserRepository{Database: db},
	}
	handler := http.HandlerFunc(cfg.handlerAddBook)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func SetupDatabase() *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Vienna"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&User{}, &Book{}, &UserBooks{})
	err = db.SetupJoinTable(&User{}, "Books", &UserBooks{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
