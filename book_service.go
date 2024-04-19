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

// Adds book to user's collection
func (cfg *config) AddBook(isbn string, userId int) (Book, error) {
	var user User
	user, err := cfg.UserRepository.Get(userId) // todo return error when user not exists, sentinel error?
	if err != nil {
		return Book{}, err
	}

	var book Book
	book = cfg.BookRepository.GetBook(isbn)
	if book.ISBN == "" {
		book, err = cfg.Client.FetchBook(isbn)
		if err != nil {
			return Book{}, err //todo return sentinel error instead empty book
		}
	}
	user.Books = append(user.Books, book)
	err = cfg.UserRepository.Save(user)
	if err != nil {
		return Book{}, err //todo return sentinel error instead empty book
	}
	return book, nil
}
