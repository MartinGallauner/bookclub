package main

import "gorm.io/gorm"

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
	err := r.Database.Table("users").Preload("Books").First(&user, id).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) Save(user User) error {
	err := r.Database.Table("users").Save(&user).Error
	return err
}

func (r *PostgresUserRepository) SearchBook(isbn string) ([]User, error) {
	var users []User
	err := r.Database.Preload("Books", "isbn = ?", isbn).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
