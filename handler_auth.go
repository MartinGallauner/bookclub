package main

import (
	"context"
	"fmt"
	"github.com/markbates/goth/gothic"
	"html/template"
	"net/http"
)

func (cfg *BookclubServer) handlerCallback(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(context.Background(), "provider", "google")) //todo I don't understand that tbh

	user, err := gothic.CompleteUserAuth(w, r)
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
	r = r.WithContext(context.WithValue(context.Background(), "provider", "google")) //todo I don't understand that tbh
	// try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(w, gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`
