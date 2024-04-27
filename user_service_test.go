package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Tests if a saved book can be linked to an existing user.
func TestCreateNewUser(t *testing.T) {
	//given
	s, err := setupTest()
	if err != nil {
		return
	}
	//when
	request, _ := http.NewRequest(http.MethodPost, "/api/users", nil)
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	assertStatus(t, response.Code, http.StatusOK)
}
