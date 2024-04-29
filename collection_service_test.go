package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Tests if a saved book can be added to the collection of an existing user.
func TestAddBookToUser(t *testing.T) {
	//given
	s, err := setupTest()

	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	s.BookRepository.Save(mockBook)
	mockUser := &User{Name: "Test User"}
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

// Add Book to an unknown user
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

	userWithBook, err := s.CreateUser("Book Owner")
	book := &Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	s.BookRepository.Save(*book)
	_, err = s.AddBookToCollection(book.ISBN, userWithBook.ID)

	userWithoutBooks, err := s.CreateUser("Reader")

	s.LinkUsers(userWithBook.ID, userWithoutBooks.ID)
	s.LinkUsers(userWithoutBooks.ID, userWithBook.ID)

	link, err := s.LinkRepository.Get(userWithBook.ID, userWithoutBooks.ID)
	fmt.Println(link)

	if err != nil {
		t.Fatalf("Unable to setup test.")
	}

	//when
	requestBody := SearchRequest{UserId: uint(1), ISBN: "1234567890"}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/search", bytes.NewReader(jsonBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got SearchResponse
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}

	assert.Equal(t, got.Isbn, "1234567890", "Did we search for the wrong book?")
	assert.Equal(t, len(got.Users), 1, ". Expected a different number of book owners.")
	assertStatus(t, response.Code, http.StatusOK)
}
