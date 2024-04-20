package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchNetwork(t *testing.T) {

	server := NewBookclubServer(Client{}, &StubBookRepository{}, &StubUserRepository{})

	t.Run("returns 200 on /search", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/search", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

	})

}
