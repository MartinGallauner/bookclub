package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (cfg *BookclubServer) handlerLinkUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := LinkRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	Link, err := cfg.LinkUsers(request.SenderId, request.ReceiverId)

	linkResponse := LinkResponse{SenderId: Link.SenderId, ReceiverId: Link.ReceiverId}

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
	CreatedAt  time.Time
	AcceptedAt time.Time
}
