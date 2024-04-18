package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *config) handlerAddBook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := AddBookRequest{}
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	book, err := cfg.AddBook(body.ISBN, body.UserId)

	if err != nil {
		respondWithError(w, 400, "Unable to add the requested book")
		return
	}
	respondWithJSON(w, 200, book) //todo reconsider response body
	return
}

type AddBookRequest struct { //todo what is the best location for this struct definition?
	UserId int    `json:"user_id"` //todo this needs to be changed later when auth is in place
	ISBN   string `json:"isbn"`
}
