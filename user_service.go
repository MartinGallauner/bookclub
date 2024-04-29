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

// Creates a link request betweem two users
func (cfg *BookclubServer) LinkUsers(senderId uint, receiverId uint) (Link, error) {
	//todo what if request already exists?
	//todo accept request when sender and receiver are inverted
	link := &Link{SenderId: senderId, ReceiverId: receiverId}
	err := cfg.LinkRepository.Save(link)
	if err != nil {
		return Link{}, err
	}
	return Link{SenderId: senderId, ReceiverId: receiverId}, nil
}

// Returns all link requests concerning the specified user
func (cfg *BookclubServer) GetLinks(userId string) ([]Link, error) {
	links, err := cfg.LinkRepository.GetById(userId)
	if err != nil {
		return nil, err
	}
	return links, nil
}

type Link struct {
	SenderId   uint `gorm:"primaryKey"`
	ReceiverId uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	AcceptedAt time.Time
	DeletedAt  time.Time
}
