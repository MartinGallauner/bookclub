package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Books []Book `gorm:"foreignKey:ISBN"`
}

func (cfg *config) AddBook(isbn string, userId int) (Book, error) {
	//todo check if book already exists, if yes load existing.
	//The book table could become a cache in the future

	book, err := cfg.Client.FetchBook(isbn)
	if err != nil {
		return Book{}, err
	}
	err = cfg.Database.Table("books").Save(&book).Error

	return book, nil
}
