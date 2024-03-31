package main

import (
	"encoding/json"
	"io"
	"net/http"
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

func (c *Client) FetchBook(isbn string) (Book, error) {
	url := baseURL + "books?jscmd=data&format=json&bibkeys=ISBN:" + isbn

	req, error := http.NewRequest("GET", url, nil)
	if error != nil {
		return Book{}, error
	}

	resp, error := c.httpClient.Do(req)
	if error != nil {
		return Book{}, error
	}
	defer resp.Body.Close()

	dat, error := io.ReadAll(resp.Body)
	if error != nil {
		return Book{}, error
	}

	response := Book{}
	error = json.Unmarshal(dat, &response)
	if error != nil {
		return Book{}, error
	}
	return response, nil
}
