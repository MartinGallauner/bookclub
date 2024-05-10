package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Search a book in the database //TODO: filter for connected users
func (cfg *BookclubServer) handlerSearch(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := AddBookRequest{}
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	//TODO: remove ID from user response
	users, err := cfg.SearchBookInNetwork(body.UserId, body.ISBN)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Book is not available in the users network.")
	}

	var responseBody []UserResponse
	for _, user := range users {
		responseBody = append(responseBody, UserResponse{Name: user.Name, Email: user.Email})
	}

	searchResponse := SearchResponse{body.ISBN, responseBody}

	respondWithJSON(w, http.StatusOK, searchResponse)
	return
}

type SearchResponse struct {
	Isbn  string         `json:"ISBN"`
	Users []UserResponse `json:"users"`
}

type UserResponse struct {
	Name  string
	Email string
}
