package main

import (
	"net/http"
)

func (cfg *BookclubServer) handlerLinkUser(w http.ResponseWriter, r *http.Request) {
	/*
		decoder := json.NewDecoder(r.Body)
		request := CreateUserRequest{}
		err := decoder.Decode(&request)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error decoding parameters: %s", err))
			return
		}*/

	var err error
	if err != nil {
		respondWithError(w, 400, "Unable to add the requested book")
		return
	}
	respondWithJSON(w, 200, nil)
	return
}

type CreateLinkRequest struct {
	senderId   uint
	receiverId uint
}
