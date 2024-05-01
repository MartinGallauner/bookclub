# book-club
Connect your physical bookshelf with your friends!


## How to run
- Checkout the project and using the terminal, navigate to its folder.
- Run `docker compose up -d` to startup the database.
- Run `go build && ./bookclub`.

## API

* You can find the swagger documentation at `http://localhost:8080/swagger/index.html`
* Fetch book from OpenLibrary API ` GET http://localhost:8080/api/books/{isbn}`
* Add book to the logged in user `POST http://localhost:8080/api/collection/{isbn}`



## Trade offs
- I regret adding Gorm as the ORM library.
- The search implementation is not okay.

