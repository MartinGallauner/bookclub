package main

import (
	"bytes"
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Tests to create a new user.
func TestCreateNewUser(t *testing.T) {
	//given
	s, err := setupTest()
	if err != nil {
		return
	}
	//when

	requestBody := CreateUserRequest{Name: "Mocki"}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(jsonBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got User
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}
	assert.Equal(t, got.Name, "Mocki")
	assert.Equal(t, got.ID, uint(1))

	assertStatus(t, response.Code, http.StatusOK)
}
