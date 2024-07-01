package server

import (
	"github.com/martingallauner/bookclub/internal"
	//"github.com/martingallauner/bookclub/internal/client"

)

// Adds book to user's collection
func (cfg *BookclubServer) AddBookToCollection(isbn string, userId uint) (internal.Book, error) {
	user, err := cfg.UserRepository.Get(userId)
	if err != nil {
		return internal.Book{}, err
	}

	var book internal.Book
	book, err = cfg.BookRepository.GetBook(isbn) 
	if err != nil {
		return internal.Book{}, err
	}
	if book.ISBN == "" {
		book, err = cfg.Client.FetchBook(isbn)
		if err != nil {
			return internal.Book{}, err
		}
	}
	user.Books = append(user.Books, book)
	err = cfg.UserRepository.Save(&user)
	if err != nil {
		return internal.Book{}, err
	}
	return book, nil
}

// Searches for the given book within the network of the given user.
func (cfg *BookclubServer) SearchBookInNetwork(userId uint, isbn string) ([]internal.User, error) {

	//get linked users
	links, err := cfg.LinkRepository.GetAcceptedById(userId)
	if err != nil {
		return nil, err
	}

	// filter user for searched book
	users, err := cfg.UserRepository.SearchBook(isbn)
	if err != nil {
		return nil, err
	}

	var collection = make(map[uint]internal.User)
	for _, user := range users {
		for _, link := range links {
			if user.ID == link.SenderId || user.ID == link.ReceiverId {
				collection[user.ID] = user
			}
		}
	}

	result := make([]internal.User, 0, len(collection))
	for _, value := range collection {
		result = append(result, value)
	}

	return result, nil
}

type SearchRequest struct {
	UserId uint   `json:"user_id"`
	ISBN   string `json:"isbn"`
}