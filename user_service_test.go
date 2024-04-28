package main

import (
	"bytes"
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Tests creation of a new user.
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

// Tests to create link request
func TestRequestLink(t *testing.T) {
	//given
	s, err := setupTest()
	if err != nil {
		return
	}

	user1, _ := s.CreateUser("Alpha")
	user2, _ := s.CreateUser("Bravo")

	if err != nil {
		t.Fatalf("Unable to prepare users needed for the test %v", err)
	}

	//when
	requestBody := LinkRequest{SenderId: user1.ID, ReceiverId: user2.ID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/links", bytes.NewReader(jsonBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got LinkResponse
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}
	assertStatus(t, response.Code, http.StatusOK)
	assert.Equal(t, got.SenderId, user1.ID)
	assert.Equal(t, got.ReceiverId, user2.ID)
	savedLink, err := s.LinkRepository.Get(user1.ID, user2.ID)
	assert.Equal(t, savedLink.SenderId, user1.ID)
	assert.Equal(t, savedLink.ReceiverId, user2.ID)
}

// Tests to create link request
func TestGetLinks(t *testing.T) {
	//given
	s, err := setupTest()
	if err != nil {
		return
	}

	user1, _ := s.CreateUser("Alpha")
	user2, _ := s.CreateUser("Bravo")

	if err != nil {
		t.Fatalf("Unable to prepare users needed for the test %v", err)
	}

	//when
	requestBody := LinkRequest{SenderId: user1.ID, ReceiverId: user2.ID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/links", bytes.NewReader(jsonBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got LinkResponse
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}
	assertStatus(t, response.Code, http.StatusOK)
	assert.Equal(t, got.SenderId, user1.ID)
	assert.Equal(t, got.ReceiverId, user2.ID)
	assert.Equal(t, got.isLinked, false)
	savedLink, err := s.LinkRepository.Get(user1.ID, user2.ID)
	assert.Equal(t, savedLink.SenderId, user1.ID)
	assert.Equal(t, savedLink.ReceiverId, user2.ID)
}
