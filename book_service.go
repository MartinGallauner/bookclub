package main

import "fmt"

func AddBook(cfg config, isbn string, userId int) (Book, error) {
	book, err := cfg.Client.FetchBook(isbn)
	if err != nil {
		return Book{}, err
	}

	//todo save to database
	fmt.Println("saved to database")
	cfg.Database.Collection[userId] = append(cfg.Database.Collection[userId], book)

	return book, nil

}
