package main

type Book struct {
	ISBN  string `gorm:"primaryKey"`
	URL   string `json:"url"`
	Title string `json:"title"`
}
