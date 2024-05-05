package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// @Description test test test
// @Summary create a new uer
func (cfg *BookclubServer) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := CreateUserRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if request.Name == "" {
		respondWithError(w, 400, "Empty username is not accepted.")
		return
	}

	user, err := cfg.CreateUser(request.Name, request.Email)

	if err != nil {
		respondWithError(w, 400, "Unable to create the user")
		return
	}
	respondWithJSON(w, 200, user)
	return
}

type CreateUserRequest struct {
	Name  string
	Email string
}