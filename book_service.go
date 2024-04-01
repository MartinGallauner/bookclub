package main

type Mock struct {
	ID   uint
	Name string
}

func AddBook(cfg config, isbn string, userId int) (Book, error) {
	book, err := cfg.Client.FetchBook(isbn)
	if err != nil {
		return Book{}, err
	}
	err = cfg.Database.Table("books").Create(&book).Error
	return book, nil
}
