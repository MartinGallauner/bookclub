package internal

import (
	"context"
	"fmt"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
	"net/http"
)

func (cfg *BookclubServer) handlerCallback(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(context.Background(), "provider", "google"))

	user, err := cfg.AuthService.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Println("logging in user: ", user.Email) //todo delete

	http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
}

func (cfg *BookclubServer) handlerLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loggin out user") //todo delete

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
		gothic.BeginAuthHandler(w, r) //todo add to interface
	}
	//check if user exists, if not, create
	persistedUser, err := cfg.UserRepository.GetByEmail(gothUser.Email)
	if err == gorm.ErrRecordNotFound {
		persistedUser, err = cfg.CreateUser(gothUser.Name, gothUser.Email)
		if err != nil {
			respondWithError(w, 400, "Unable to create new user.")
		}
	}
	jwt, err := cfg.JwtService.CreateToken("bookclub-access", int(persistedUser.ID))
	if err != nil {
		//todo logging
		respondWithError(w, 400, "Unable to login user.")
	}
	//todo return login response
	loginResponse := LoginResponse{Name: persistedUser.Name, Email: persistedUser.Email, Jwt: jwt}
	respondWithJSON(w, 200, loginResponse)
}

type LoginResponse struct {
	Name  string
	Email string
	Jwt   string //todo reconsider naming
}
