package repository

import (
	"fmt"
	"gorm.io/gorm"
	. "github.com/martingallauner/bookclub/internal"
)

type UserRepository interface {
	Get(id uint) (User, error)
	GetByEmail(email string) (User, error)
	Save(user *User) error
	SearchBook(isbn string) ([]User, error)
}

type BookRepository interface {
	GetBook(isbn string) (Book, error)
	Save(book Book) error
}

type LinkRepository interface {
	//Returns specific Link between two users
	Get(senderId uint, receiverId uint) (Link, error)

	//Returns all links concerned with the user
	GetById(userId string) ([]Link, error)

	//Returns all links concerned with the user that are accepted
	GetAcceptedById(userId uint) ([]Link, error)

	//Persists link.
	Save(link *Link) error
}

type PostgresBookRepository struct {
	Database *gorm.DB
}

func (r *PostgresBookRepository) GetBook(isbn string) (Book, error) {
	var book Book
	err := r.Database.Table("books").Find(&book, isbn).Error
	if err != nil {
		return Book{}, nil
	}
	return book, nil
}

func (r *PostgresBookRepository) Save(book Book) error {
	err := r.Database.Table("books").Save(&book).Error
	return err
}

type PostgresUserRepository struct {
	Database *gorm.DB
}

func (r *PostgresUserRepository) Get(id uint) (User, error) {
	var user User
	err := r.Database.Table("users").Preload("Books").First(&user, id).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (User, error) {
	var user User
	err := r.Database.Table("users").Preload("Books").First(&user).Where("email = ?", email).Error
	if err != nil {
		return User{}, err
	} //TODO: return ErrNotFound
	return user, nil
}

func (r *PostgresUserRepository) Save(user *User) error {
	err := r.Database.Table("users").Save(&user).Error
	return err
}

func (r *PostgresUserRepository) SearchBook(isbn string) ([]User, error) {
	var users []User
	result := r.Database.Raw(fmt.Sprintf("SELECT * FROM users\nJOIN user_books ON users.id = user_books.user_id\nWHERE book_isbn = '%v'", isbn)).Scan(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

type PostgresLinkRepository struct {
	Database *gorm.DB
}

func (r *PostgresLinkRepository) Get(senderId uint, receiverId uint) (Link, error) {
	var link Link
	err := r.Database.Where("sender_id = ? AND receiver_id = ?", senderId, receiverId).First(&link).Error
	if err != nil {
		return Link{}, err
	}

	if link.SenderId == 0 && link.ReceiverId == 0 { //TODO: I feel like that check is bad
		return Link{}, gorm.ErrRecordNotFound
	}
	return link, nil
}

func (r *PostgresLinkRepository) GetById(userId string) ([]Link, error) {
	var links []Link

	// Build the query with OR condition
	result := r.Database.Where("sender_id = ? OR receiver_id = ?", userId, userId).Find(&links)

	if result.Error != nil {
		// handle error
		return nil, result.Error
	}
	return links, nil
}

func (r *PostgresLinkRepository) GetAcceptedById(userId uint) ([]Link, error) {
	var links []Link

	// Build the query with OR condition
	result := r.Database.Where("sender_id = ? OR receiver_id = ?", userId, userId).Where("accepted_at > created_at").Find(&links)

	if result.Error != nil {
		// handle error
		return nil, result.Error
	}
	return links, nil
}

func (r *PostgresLinkRepository) Save(link *Link) error {
	err := r.Database.Table("links").Save(&link).Error
	return err
}
