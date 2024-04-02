package main

import (
	"net/http"
)

func (cfg *config) handlerAddBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")
	book, err := cfg.AddBook(isbn, 6)

	if err != nil {
		respondWithError(w, 400, "Unable to add the requested book")
		return
	}
	respondWithJSON(w, 200, book)
	return
}
