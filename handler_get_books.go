package main

import (
	"net/http"
)

func (cfg *config) handlerGetBookByISBN(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")
	book, err := cfg.Client.FetchBook(isbn)
	if err != nil {
		respondWithError(w, 400, "Unable to fetch the requested book")
		return
	}
	respondWithJSON(w, 200, book)
	return
}
