package internal

import (
	"context"
	"fmt"
	"github.com/markbates/goth/gothic"
	"net/http"
)

func (cfg *BookclubServer) handlerCallback(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(context.Background(), "provider", "google"))

	user, err := cfg.AuthService.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Println("logging in user: ", user.Email)

	http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
}

func (cfg *BookclubServer) handlerLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loggin out user")

	gothic.Logout(w, r)
	w.Header().Set("Location", "http://localhost:5173")
	w.WriteHeader(http.StatusTemporaryRedirect)

	//http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
}

func (cfg *BookclubServer) handlerLogin(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(context.Background(), "provider", "google")) //todo I don't fully understand that tbh
	// try to get the user without re-authenticating
	gothUser, err := cfg.AuthService.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
	}
	//check if user exists, if not, create
	persisted, err := cfg.UserRepository.GetByEmail(gothUser.Email) //todo validate
	fmt.Println(persisted)

	respondWithJSON(w, 200, gothUser) //todo return jwt
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}
