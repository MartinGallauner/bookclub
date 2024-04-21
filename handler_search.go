package main

import (
	"net/http"
)

// Search a book in the database //todo filter for connected users
func (cfg *BookclubServer) handlerSearch(w http.ResponseWriter, r *http.Request) {
	//book, err := cfg.Client.FetchBook(isbn)
	//if err != nil {
	//	respondWithError(w, 400, "Unable to fetch the requested book")
	//	return
	//}

	respondWithJSON(w, 200, Book{ISBN: "1234", Title: "title", URL: "url"})
	return
}
