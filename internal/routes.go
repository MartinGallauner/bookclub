package internal

import "net/http"

func addRoutes(
	mux *http.ServeMux,
	bookRepository *BookRepository,
	userRepository *UserRepository,
	linkRepository *LinkRepository,
) {
	mux.Handle("/api/search", handleSearch2(linkRepository, userRepository))

}
