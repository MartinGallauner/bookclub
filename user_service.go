package main

import "gorm.io/gorm"

func (cfg *BookclubServer) CreateUser(name string) (User, error) {
	user := &User{Name: name}
	err := cfg.UserRepository.Save(user)
	if err != nil {
		return User{}, nil
	}
	return *user, nil
}

func (cfg *BookclubServer) LinkUsers(senderId uint, receiverId uint) (Link, error) {
	return Link{SenderId: senderId, ReceiverId: receiverId}, nil

}

type Link struct {
	SenderId   uint
	ReceiverId uint
	gorm.Model
}
