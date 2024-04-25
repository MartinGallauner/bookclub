package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchNetwork(t *testing.T) {

	server := NewBookclubServer(Client{}, &StubBookRepository{books: make(map[string]Book)}, &StubUserRepository{users: make(map[int]User)})

	/*t.Run("returns 200 on /search", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/api/search/1234", nil)
		response := httptest.NewRecorder()
		book := Book{ISBN: "1234", Title: "title", URL: "url"}
		user := User{Name: "User", Books: []Book{book}}
		err := server.UserRepository.Save(user)
		if err != nil {
			return
		}

		server.ServeHTTP(response, request)

		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error("Unable to read response")
		}
		var searchResponse SearchResponse
		err = json.Unmarshal(body, &searchResponse)

		assertStatus(t, response.Code, http.StatusOK)
		assert.Equal(t, searchResponse.Isbn, "1234")
		assert.Equal(t, searchResponse.Users[0].Name, "User")
	})*/

	t.Run("returns 404 on /search when book doesn't exist", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/api/search/1234", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})

}
