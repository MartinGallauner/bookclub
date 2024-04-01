package main

type Book struct {
	ID    uint
	ISBN  string
	URL   string `json:"url"`
	Title string `json:"title"`
}
