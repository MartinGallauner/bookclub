package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Search a book in the database //todo filter for connected users
func (cfg *BookclubServer) handlerSearch(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := AddBookRequest{}
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	users, err := cfg.SearchBookInNetwork(body.UserId, body.ISBN)
	if err != nil {
		respondWithError(w, 404, "Book is not available in the users network.")
	}

	searchResponse := SearchResponse{body.ISBN, users}

	respondWithJSON(w, 200, searchResponse)
	return
}

type SearchResponse struct {
	Isbn  string `json:"ISBN"`
	Users []User `json:"users"`
}
