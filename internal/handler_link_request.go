package internal

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
	linkResponse := mapLinkResponse(link)
	if err != nil {
		respondWithError(w, 400, "Unable to create user link")
		return
	}
	respondWithJSON(w, 200, linkResponse)
	return

}

func mapLinkResponse(link Link) LinkResponse {
	linkResponse := LinkResponse{SenderId: link.SenderId, ReceiverId: link.ReceiverId, IsLinked: false}
	if link.DeletedAt.Before(link.AcceptedAt) {
		linkResponse.IsLinked = true
	}
	return linkResponse
}

type LinkRequest struct {
	SenderId   uint `json:"sender_id"`
	ReceiverId uint `json:"receiver_id"`
}

type LinkResponse struct {
	SenderId   uint `json:"sender_id"`
	ReceiverId uint `json:"receiver_id"`
	IsLinked   bool `json:"is_linked"`
}
