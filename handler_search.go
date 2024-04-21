package main

import (
	"net/http"
)

// Search a book in the database //todo filter for connected users
func (cfg *BookclubServer) handlerSearch(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")

	//book, err := cfg.Client.FetchBook(isbn)
	//if err != nil {
	//	respondWithError(w, 400, "Unable to fetch the requested book")
	//	return
	//}
	book := Book{ISBN: isbn, Title: "title", URL: "url"}
	user := User{Name: "Tester", Books: []Book{book}}
	response := SearchResponse{Isbn: isbn, Users: []User{user}}

	respondWithJSON(w, 200, response)
	return
}

type SearchResponse struct {
	Isbn  string `json:"ISBN"`
	Users []User `json:"users"`
}
