package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

const (
	baseURL = "https://openlibrary.org/api/"
)

// Fetches a single book by it's ISBN from the OpenLibraryAPI
func (c *Client) FetchBook(isbn string) (Book, error) {
	url := baseURL + "books?jscmd=data&format=json&bibkeys=ISBN:" + isbn

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Book{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Book{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Book{}, err
	}

	var response map[string]map[string]interface{}
	err = json.Unmarshal(dat, &response)
	if err != nil {
		return Book{}, err
	}

	book, err := mapBookResponse(response)
	if err != nil {
		fmt.Errorf("Failed mapping response to book, %w", err)
	}

	return book, nil
}

// mapBookResponse reads the ugly response from the OpenLibraryAPI and maps it to a simple book entity.
// todo clean that crap up
func mapBookResponse(response map[string]map[string]interface{}) (Book, error) {
	var book Book
	for isbn, bookData := range response {
		book.URL = bookData["url"].(string)
		book.ISBN = strings.TrimPrefix(isbn, "ISBN:")
		book.Title = bookData["title"].(string)
	}
	return book, nil
}
