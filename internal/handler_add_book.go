package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *BookclubServer) handlerAddBook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := AddBookRequest{}
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	book, err := cfg.AddBookToCollection(body.ISBN, body.UserId)

	if err != nil {
		respondWithError(w, 400, "Unable to add the requested book")
		return
	}
	respondWithJSON(w, 200, book) //todo reconsider response body
	return
}

type AddBookRequest struct { //todo what is the best location for this struct definition?
	UserId uint   `json:"user_id"` //todo this needs to be changed later when auth is in place
	ISBN   string `json:"isbn"`
}