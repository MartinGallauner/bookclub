package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type StubBookRepository struct {
	books map[string]Book
}

func (r *StubBookRepository) GetBook(isbn string) Book {
	return r.books[isbn]
}

type StubUserRepository struct {
	users map[int]User
}

func (r *StubUserRepository) Get(id int) User {
	return r.users[id]
}

func TestPOSTBookToCollection(t *testing.T) {
	bookRepository := StubBookRepository{
		map[string]Book{
			"1234": {ISBN: "1234", URL: "url", Title: "title"},
		},
	}
	userRepository := StubUserRepository{
		map[int]User{
			1: {Name: "Anna"},
		},
	}

	cfg := config{BookRepository: &bookRepository, UserRepository: &userRepository}

	t.Run("Add book to user's book collection", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/api/collections", strings.NewReader(`{"user_id": 1234, "isbn": "url"}`))
		response := httptest.NewRecorder()

		cfg.handlerAddBook(response, request)

		want := "{1234 url title}"
		assertResponseBody(t, response.Body.String(), want)
		assertStatus(t, response.Code, http.StatusOK)
	})
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
