package server

import (
	"net/http"
)

// handlerGetBookByISBN fetches data directly from the OpenLibrary API. Most likely I will delete that soon.
func (cfg *BookclubServer) handlerGetBookByISBN(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")
	book, err := cfg.Client.FetchBook(isbn)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Unable to fetch the requested book")
		return
	}
	RespondWithJSON(w, http.StatusOK, book)
	return
}
