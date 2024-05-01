package internal

import (
	"net/http"
)

// That handler fetches data directly from the OpenLibrary API. Most likely I will delete that soon.
func (cfg *BookclubServer) handlerGetBookByISBN(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")
	book, err := cfg.Client.FetchBook(isbn)
	if err != nil {
		respondWithError(w, 400, "Unable to fetch the requested book")
		return
	}
	respondWithJSON(w, 200, book)
	return
}
