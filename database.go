package main

// This is just a stupid placeholder. The key is the userId, mapping to a slice of books
type DB struct {
	Collection map[int][]Book `json:"collections"`
}
