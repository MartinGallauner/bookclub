package main

import (
	"net/http"
)

func (cfg *BookclubServer) handlerGetLinks(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id") //todo in the future I want to read the userId from the token

	link, err := cfg.GetLinks(userId)

	if err != nil {
		respondWithError(w, 400, "Unable to create user link")
		return
	}
	respondWithJSON(w, 200, link)
	return
}
