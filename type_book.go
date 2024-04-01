package main

type Book struct {
	ISBN        string
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	Pages       int      `json:"pagination"`
	PublishDate int      `json:"publish_date"`
	Subjects    []string `json:"subjects"`
}
