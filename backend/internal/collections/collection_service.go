package collections

import (
	"github.com/martingallauner/bookclub/internal"
	"github.com/martingallauner/bookclub/internal/client"
	"github.com/martingallauner/bookclub/internal/repository"
)

type Service struct {
	userRepository repository.UserRepository
	bookRepository repository.BookRepository
	linkRepository repository.LinkRepository
	client         client.Client
}

func New(userRepository repository.UserRepository, bookRepository repository.BookRepository, linkRepository repository.LinkRepository, client client.Client) *Service {
	return &Service{userRepository: userRepository, bookRepository: bookRepository, linkRepository: linkRepository}
}

// Adds book to user's collection
func (srv *Service) AddBookToCollection(isbn string, userId uint) (internal.Book, error) {
	user, err := srv.userRepository.Get(userId)
	if err != nil {
		return internal.Book{}, err
	}

	var book internal.Book
	book, err = srv.bookRepository.GetBook(isbn)
	if err != nil {
		return internal.Book{}, err
	}
	if book.ISBN == "" {
		book, err = srv.client.FetchBook(isbn)
		if err != nil {
			return internal.Book{}, err
		}
	}
	user.Books = append(user.Books, book)
	err = srv.userRepository.Save(&user)
	if err != nil {
		return internal.Book{}, err
	}
	return book, nil
}

// Searches for the given book within the network of the given user.
func (srv *Service) SearchBookInNetwork(userId uint, isbn string) ([]internal.User, error) {

	//get linked users
	links, err := srv.linkRepository.GetAcceptedById(userId)
	if err != nil {
		return nil, err
	}

	// filter user for searched book
	users, err := srv.userRepository.SearchBook(isbn)
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
