# book-club
Connect your physical bookshelf with your friends!


## How to run
- Checkout the project and using the terminal, navigate to its folder
- Run `go build && ./bookclub`

## API

* Fetch book from OpenLibrary API ` GET http://localhost:8080/api/books/{isbn}`
* Add book to the logged in user `POST http://localhost:8080/api/collection/{isbn}`



