package main

// Adds book to user's collection
func (cfg *BookclubServer) AddBookToCollection(isbn string, userId uint) (Book, error) {
	user, err := cfg.UserRepository.Get(userId)
	if err != nil {
		return Book{}, err
	}

	var book Book
	book, _ = cfg.BookRepository.GetBook(isbn) //todo handle error
	if book.ISBN == "" {
		book, err = cfg.Client.FetchBook(isbn)
		if err != nil {
			return Book{}, err
		}
	}
	user.Books = append(user.Books, book)
	err = cfg.UserRepository.Save(&user)
	if err != nil {
		return Book{}, err
	}
	return book, nil
}

func (cfg *BookclubServer) SearchBookInNetwork(userId uint, isbn string) ([]User, error) {

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

	var collection = make(map[uint]User)
	for _, user := range users {
		for _, link := range links {
			if user.ID == link.SenderId || user.ID == link.ReceiverId {
				collection[user.ID] = user
			}
		}
	}

	result := make([]User, 0, len(collection))
	for _, value := range collection {
		result = append(result, value)
	}

	return result, nil
}

type SearchRequest struct {
	UserId uint   `json:"user_id"`
	ISBN   string `json:"isbn"`
}
