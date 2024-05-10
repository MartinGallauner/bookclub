package internal

import (
	"net/http"
)

func (cfg *BookclubServer) handlerGetLinks(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id") //TODO: in the future I want to read the userId from the token

	links, err := cfg.GetLinks(userId)
	var result []LinkResponse
	for _, link := range links {
		result = append(result, mapLinkResponse(link))
	}

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to create user link")
		return
	}
	respondWithJSON(w, http.StatusOK, result)
	return
}
