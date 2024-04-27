package main

import (
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Tests if a saved book can be linked to an existing user.
func TestLinkBookToUser(t *testing.T) {
	//given
	s, err := setupTest()

	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	s.BookRepository.Save(mockBook)
	mockUser := User{Name: "Test User"}
	mockUser.ID = 1
	err = s.UserRepository.Save(mockUser)
	if err != nil {
		return
	}
	//when
	request, _ := http.NewRequest(http.MethodPost, "/api/collections", strings.NewReader(`{"user_id": 1, "isbn": "1234567890"}`))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	assertStatus(t, response.Code, http.StatusOK)
	user, err := s.UserRepository.Get(1)

	if len(user.Books) == 0 {
		t.Errorf("User has no books")
		t.FailNow()
	}
	addedBook := user.Books[0]
	assert.Equal(t, addedBook, mockBook, "Added book should match existing book")
}

func TestAddBookToUnknownUser(t *testing.T) {
	//given
	s, err := setupTest()
	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	s.BookRepository.Save(mockBook)

	if err != nil {
		return
	}

	//when
	request, _ := http.NewRequest(http.MethodPost, "/api/collections", strings.NewReader(`{"user_id": 99, "isbn": "1234567890"}`))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	assertStatus(t, response.Code, http.StatusBadRequest)
}

// Tests search function when one user has the book
func TestSearchBookInNetwork(t *testing.T) {
	//given
	s, err := setupTest()
	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	mockUser := User{Name: "Test User", Books: []Book{mockBook}}
	mockUser.ID = 1

	err = s.UserRepository.Save(mockUser)
	if err != nil {
		return
	}

	//when
	request, _ := http.NewRequest(http.MethodGet, "/api/search/1234567890", nil)
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got SearchResponse
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
	}

	assert.Equal(t, got.Isbn, "1234567890")
	assertStatus(t, response.Code, http.StatusOK)
}
