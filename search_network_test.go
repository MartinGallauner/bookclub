package main

import (
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchNetwork(t *testing.T) {

	server := NewBookclubServer(Client{}, &StubBookRepository{}, &StubUserRepository{})

	t.Run("returns 200 on /search", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/search/1234", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error("Unable to read response")
		}
		var searchResponse SearchResponse
		err = json.Unmarshal(body, &searchResponse)

		//assertResponseBody(t, response.Body.String(), want)

		assertStatus(t, response.Code, http.StatusOK)
		assert.Equal(t, searchResponse.Isbn, "1234")
		assert.Equal(t, searchResponse.Users[0].Name, "Tester")

	})

}
