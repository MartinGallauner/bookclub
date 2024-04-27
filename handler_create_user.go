package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *BookclubServer) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := CreateUserRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	user, err := cfg.CreateUser(request.Name)

	if err != nil {
		respondWithError(w, 400, "Unable to add the requested book")
		return
	}
	respondWithJSON(w, 200, user)
	return
}

type CreateUserRequest struct {
	Name string
}
