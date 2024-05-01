package internal

import (
	"gorm.io/gorm"
	"log"
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
	//ask if request already exists, if yes exit
	existingLink, err := cfg.LinkRepository.Get(senderId, receiverId)

	if existingLink.SenderId == senderId && existingLink.ReceiverId == receiverId { //if already exists, return that right away
		return existingLink, err
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Creating new link request for users with id %v and %v", senderId, receiverId)
		} else {
			return Link{}, err
		}
	}

	//if opposite request exists -> accept
	existingRequest, err := cfg.LinkRepository.Get(receiverId, senderId)
	if err != gorm.ErrRecordNotFound {
		existingRequest.AcceptedAt = time.Now()
		err := cfg.LinkRepository.Save(&existingRequest)
		if err != nil {
			return Link{}, err
		}
		return existingRequest, nil
	}

	//if request not exists -> create new
	link := &Link{SenderId: senderId, ReceiverId: receiverId}
	err = cfg.LinkRepository.Save(link)
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
