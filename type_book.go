package main

type Authors struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}
type Identifiers struct {
	Isbn13      []string `json:"isbn_13"`
	Openlibrary []string `json:"openlibrary"`
}
type Publishers struct {
	Name string `json:"name"`
}
type Subjects struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Book struct {
	URL         string       `json:"url"`
	Key         string       `json:"key"`
	Title       string       `json:"title"`
	Authors     []Authors    `json:"authors"`
	Pagination  string       `json:"pagination"`
	Weight      string       `json:"weight"`
	Identifiers Identifiers  `json:"identifiers"`
	Publishers  []Publishers `json:"publishers"`
	PublishDate string       `json:"publish_date"`
	Subjects    []Subjects   `json:"subjects"`
}
