package main

import (
	"net/http"
)

// Search a book in the database //todo filter for connected users
func (cfg *BookclubServer) handlerSearch(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")

	users, err := cfg.SearchBookInNetwork(isbn)
	if err != nil {
		respondWithError(w, 404, "Book not available in network.")
	}

	searchResponse := SearchResponse{isbn, users}

	respondWithJSON(w, 200, searchResponse)
	return
}

type SearchResponse struct {
	Isbn  string `json:"ISBN"`
	Users []User `json:"users"`
}
