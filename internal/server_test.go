package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_jwtMiddleware_missing_jwt(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	mw := jwtMiddleware(handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mw(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	var got responseError
	err := json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", resp.Body, err)
	}

	if got.Error != "Missing JWT" {
		t.Errorf("expected error message 'Missing JWT', got %q", got.Error)
	}
}

func Test_jwtMiddleware_invalid_jwt(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	mw := jwtMiddleware(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer invalid")
	w := httptest.NewRecorder()
	mw(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	var got responseError
	err := json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", resp.Body, err)
	}

	if got.Error != "Invalid JWT" {
		t.Errorf("expected error message 'Missing JWT', got %q", got.Error)
	}
}

type responseError struct {
	Error string `json:"error"`
}
