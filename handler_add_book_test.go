package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type StubBookRepository struct {
	books map[string]Book
}

func (r *StubBookRepository) GetBook(isbn string) Book {
	return r.books[isbn]
}

func (r *StubBookRepository) Save(book Book) error {
	len := len(r.books) //todo this stub is crap
	r.books[strconv.Itoa(len+1)] = book
	return nil
}

type StubUserRepository struct {
	users map[int]User
}

func (r *StubUserRepository) Get(id int) (User, error) {
	return r.users[id], nil
}

func (r *StubUserRepository) Save(user User) error {
	len := len(r.users) //todo this stub is crap
	r.users[len+1] = user
	return nil
}

func TestPOSTBookToCollection(t *testing.T) {
	//given
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

	//when
	t.Run("Add book to user's book collection", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/api/collections", strings.NewReader(`{"user_id": 1, "isbn": "1234"}`))
		response := httptest.NewRecorder()

		cfg.handlerAddBook(response, request)

		want := "{\"ISBN\":\"1234\",\"url\":\"https://openlibrary.org/books/OL28151326M/The_Wednesday_surprise\",\"title\":\"The Wednesday surprise\"}"
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
