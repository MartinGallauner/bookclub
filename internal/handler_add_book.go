package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// handlerAddBook handles requests to add books to user collection
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
	respondWithJSON(w, 200, book) //TODO: reconsider response body
	return
}

type AddBookRequest struct { //TODO: what is the best location for this struct definition?
	UserId uint   `json:"user_id"` //TODO: this needs to be changed later when auth is in place
	ISBN   string `json:"isbn"`
}
