package main

import (
	"fmt"
	"net/http"
)

func (cfg *config) handlerGetBookByISBN(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called handler")

	isbn := r.PathValue("isbn")

	response, err := cfg.Client.FetchBook(isbn)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	fmt.Println(response)
	w.WriteHeader(200)
	return
}
