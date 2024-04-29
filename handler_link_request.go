package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *BookclubServer) handlerCreateLink(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := LinkRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	link, err := cfg.LinkUsers(request.SenderId, request.ReceiverId)
	linkResponse := LinkResponse{SenderId: link.SenderId, ReceiverId: link.ReceiverId}
	if link.DeletedAt.Before(link.AcceptedAt) {
		linkResponse.isLinked = true
	}
	if err != nil {
		respondWithError(w, 400, "Unable to create user link")
		return
	}
	respondWithJSON(w, 200, linkResponse)
	return

}

type LinkRequest struct {
	SenderId   uint
	ReceiverId uint
}

type LinkResponse struct {
	SenderId   uint
	ReceiverId uint
	isLinked   bool
}
