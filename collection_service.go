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
func (cfg *BookclubServer) AddBookToCollection(isbn string, userId int) (Book, error) {
	user, err := cfg.UserRepository.Get(userId)
	if err != nil {
		return Book{}, err
	}

	var book Book
	book, _ = cfg.BookRepository.GetBook(isbn) //todo handle error
	if book.ISBN == "" {
		book, err = cfg.Client.FetchBook(isbn)
		if err != nil {
			return Book{}, err
		}
	}
	user.Books = append(user.Books, book)
	err = cfg.UserRepository.Save(&user)
	if err != nil {
		return Book{}, err
	}
	return book, nil
}

func (cfg *BookclubServer) SearchBookInNetwork(userId int, isbn string) ([]User, error) {
	users, err := cfg.UserRepository.SearchBook(isbn)
	if err != nil {
		return nil, err
	}
	return users, nil
}

type SearchRequest struct {
	UserId uint   `json:"user_id"`
	ISBN   string `json:"isbn"`
}
