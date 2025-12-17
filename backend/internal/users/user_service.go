package users

import (
	"github.com/martingallauner/bookclub/internal"
	"github.com/martingallauner/bookclub/internal/client"
	"github.com/martingallauner/bookclub/internal/repository"
	"gorm.io/gorm"
	"log"
	"time"
)

type Service struct {
	userRepository repository.UserRepository
	bookRepository repository.BookRepository
	linkRepository repository.LinkRepository
}

// TODO: check whats needed
func New(userRepository repository.UserRepository, bookRepository repository.BookRepository, linkRepository repository.LinkRepository, client client.Client) *Service {
	return &Service{userRepository: userRepository, bookRepository: bookRepository, linkRepository: linkRepository}
}

func (srv *Service) CreateUser(name, email string) (internal.User, error) {
	user := &internal.User{Name: name, Email: email}
	err := srv.userRepository.Save(user)
	if err != nil {
		return internal.User{}, nil
	}
	return *user, nil
}

// Creates a link request betweem two users
func (srv *Service) LinkUsers(senderId uint, receiverId uint) (internal.Link, error) {
	//ask if request already exists, if yes exit
	existingLink, err := srv.linkRepository.Get(senderId, receiverId)

	if existingLink.SenderId == senderId && existingLink.ReceiverId == receiverId { //if already exists, return that right away
		return existingLink, err
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Creating new link request for users with id %v and %v", senderId, receiverId)
		} else {
			return internal.Link{}, err
		}
	}

	//if opposite request exists -> accept
	existingRequest, err := srv.linkRepository.Get(receiverId, senderId)
	if err != gorm.ErrRecordNotFound {
		existingRequest.AcceptedAt = time.Now()
		err := srv.linkRepository.Save(&existingRequest)
		if err != nil {
			return internal.Link{}, err
		}
		return existingRequest, nil
	}

	//if request not exists -> create new
	link := &internal.Link{SenderId: senderId, ReceiverId: receiverId}
	err = srv.linkRepository.Save(link)
	if err != nil {
		return internal.Link{}, err
	}
	return internal.Link{SenderId: senderId, ReceiverId: receiverId}, nil
}

// Returns all link requests concerning the specified user
func (srv *Service) GetLinks(userId string) ([]internal.Link, error) {
	links, err := srv.linkRepository.GetById(userId)
	if err != nil {
		return nil, err
	}
	return links, nil
}
