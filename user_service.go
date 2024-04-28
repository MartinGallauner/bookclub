package main

import (
	"time"
)

func (cfg *BookclubServer) CreateUser(name string) (User, error) {
	user := &User{Name: name}
	err := cfg.UserRepository.Save(user)
	if err != nil {
		return User{}, nil
	}
	return *user, nil
}

func (cfg *BookclubServer) LinkUsers(senderId uint, receiverId uint) (Link, error) {
	link := &Link{SenderId: senderId, ReceiverId: receiverId}
	err := cfg.LinkRepository.Save(link)
	if err != nil {
		return Link{}, err
	}
	return Link{SenderId: senderId, ReceiverId: receiverId}, nil
}

type Link struct {
	SenderId   uint `gorm:"primaryKey"`
	ReceiverId uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	AcceptedAt time.Time
}
