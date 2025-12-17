package internal

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"net/http"
)

type AuthService interface {
	CompleteUserAuth(w http.ResponseWriter, r *http.Request) (goth.User, error)
}

type GothicAuthService struct {
}

func (svc *GothicAuthService) CompleteUserAuth(w http.ResponseWriter, r *http.Request) (goth.User, error) {
	return gothic.CompleteUserAuth(w, r)
}
