package main

import (
	"net/http"
)

// That handler fetches data directly from the OpenLibrary API. Most likely I will delete that soon.
func (cfg *BookclubServer) handlerSearch(w http.ResponseWriter, r *http.Request) {
	//book, err := cfg.Client.FetchBook(isbn)
	//if err != nil {
	//	respondWithError(w, 400, "Unable to fetch the requested book")
	//	return
	//}
	respondWithJSON(w, 200, nil)
	return
}
