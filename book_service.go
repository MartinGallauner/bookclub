package main

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name  string
	Books []Book `gorm:"many2many:user_books"`
}

type Book struct {
	ISBN  string `gorm:"primaryKey"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

type UserBooks struct {
	UserId    int    `gorm:"primaryKey"`
	BookISBN  string `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (cfg *config) AddBook(isbn string, userId int) (Book, error) {
	var user User
	cfg.Database.Table("users").Find(&user, userId)

	var err error
	var book Book
	cfg.BookRepository.Database.Table("books").Find(&book, isbn) //todo abstract
	if book.ISBN == "" {
		book, err = cfg.Client.FetchBook(isbn)
		if err != nil {
			return Book{}, err //todo return sentinel error instead empty book
		}
	}
	user.Books = append(user.Books, book)
	err = cfg.Database.Table("users").Save(&user).Error
	if err != nil {
		return Book{}, err //todo return sentinel error instead empty book
	}
	return book, nil
}
