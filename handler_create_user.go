package main

import (
	"net/http"
)

func (cfg *BookclubServer) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	//decoder := json.NewDecoder(r.Body)
	//body := AddBookRequest{}
	//err := decoder.Decode(&body)
	//if err != nil {
	//	respondWithError(w, 400, fmt.Sprintf("Error decoding parameters: %s", err))
	//	return
	//}
	////book, err := cfg.AddBookToCollection(body.ISBN, body.UserId)
	//
	//if err != nil {
	//	respondWithError(w, 400, "Unable to add the requested book")
	//	return
	//}
	respondWithJSON(w, 200, User{Name: "Mocki"})
	return
}

type CreateUserRequest struct {
	Name string
}
