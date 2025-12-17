package server

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
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if request.Name == "" {
		RespondWithError(w, http.StatusBadRequest, "Empty username is not accepted.")
		return
	}

	user, err := cfg.UsersService.CreateUser(request.Name, request.Email)

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Unable to create the user")
		return
	}
	RespondWithJSON(w, http.StatusOK, user)
	return
}

type CreateUserRequest struct {
	Name  string
	Email string
}
