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

func TestPOSTBookToCollection(t *testing.T) {
	bookRepository := StubBookRepository{
		map[string]Book{
			"1234": {ISBN: "1234", URL: "url", Title: "title"},
		},
	}
	server := &BookclubServer{&bookRepository}

	t.Run("Add book to user's book collection", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/api/collections", strings.NewReader(`{"user_id": 1234, "isbn": "url"}`))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := "{1234 url title}"
		assertResponseBody(t, response.Body.String(), want)
	})
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
