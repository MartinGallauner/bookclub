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

type Link struct {
	SenderId   uint `gorm:"primaryKey"` //todo the concept of sender/receiver id is crap
	ReceiverId uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	AcceptedAt time.Time
	DeletedAt  time.Time
}
