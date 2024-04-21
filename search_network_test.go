package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchNetwork(t *testing.T) {

	server := NewBookclubServer(Client{}, &StubBookRepository{}, &StubUserRepository{})

	t.Run("returns 200 on /search", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/search", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := "{\"ISBN\":\"1234\",\"url\":\"url\",\"title\":\"title\"}"

		assertResponseBody(t, response.Body.String(), want)

		assertStatus(t, response.Code, http.StatusOK)

	})

}
