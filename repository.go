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
