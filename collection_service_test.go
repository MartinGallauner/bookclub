package main

import (
	"context"
	"github.com/magiconair/properties/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	ctx = context.Background()
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

// Tests if a saved book can be linked to an existing user.
func TestLinkBookToUser(t *testing.T) {
	//given
	container, err := CreatePostgresContainer()
	if err != nil {
		log.Fatal(err)
	}

	db, err := SetupDatabase(container.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &BookclubServer{
		BookRepository: &PostgresBookRepository{Database: db},
		UserRepository: &PostgresUserRepository{Database: db},
	}

	s := NewBookclubServer(Client{}, &PostgresBookRepository{Database: db}, &PostgresUserRepository{Database: db})

	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	cfg.BookRepository.Save(mockBook)
	mockUser := User{Name: "Test User"}
	mockUser.ID = 1
	err = cfg.UserRepository.Save(mockUser)
	if err != nil {
		return
	}
	//todo extract setup code
	//when
	request, _ := http.NewRequest(http.MethodPost, "/api/collections", strings.NewReader(`{"user_id": 1, "isbn": "1234567890"}`))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	assertStatus(t, response.Code, http.StatusOK)
	user, err := cfg.UserRepository.Get(1)

	if len(user.Books) == 0 {
		t.Errorf("User has no books")
		t.FailNow()
	}
	addedBook := user.Books[0]

	assert.Equal(t, addedBook, mockBook, "Added book should match existing book")
}

func TestLinkBookToUnknownUser(t *testing.T) {
	container, err := CreatePostgresContainer()
	if err != nil {
		log.Fatal(err)
	}

	db, err := SetupDatabase(container.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &BookclubServer{
		BookRepository: &PostgresBookRepository{Database: db},
		UserRepository: &PostgresUserRepository{Database: db},
	}

	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	cfg.BookRepository.Save(mockBook)

	if err != nil {
		return
	}
	//todo extract setup code

	_, err = cfg.AddBookToCollection(mockBook.ISBN, 1)

	assert.Equal(t, err.Error(), "record not found", "Expecting record not found error")
}

// Tests if a saved book can be linked to an existing user.
func TestSearchBookInNetwork(t *testing.T) {
	container, err := CreatePostgresContainer()
	if err != nil {
		log.Fatal(err)
	}

	db, err := SetupDatabase(container.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &BookclubServer{
		BookRepository: &PostgresBookRepository{Database: db},
		UserRepository: &PostgresUserRepository{Database: db},
	}

	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	mockUser := User{Name: "Test User", Books: []Book{mockBook}}
	mockUser.ID = 1

	err = cfg.UserRepository.Save(mockUser)
	if err != nil {
		return
	}
	users, err := cfg.SearchBookInNetwork("1234567890")
	assert.Equal(t, len(users), 1)
	assert.Equal(t, users[0].Name, "Test User")
	assert.Equal(t, len(users[0].Books), 1)
	assert.Equal(t, users[0].Books[0].Title, "Test Book")
	assert.Equal(t, users[0].Books[0].URL, "https://...")
	assert.Equal(t, users[0].Books[0].ISBN, "1234567890")

}

func CreatePostgresContainer() (*PostgresContainer, error) {
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithDatabase("bookclub-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &PostgresContainer{
		PostgresContainer: postgresContainer,
		ConnectionString:  connStr,
	}, nil
}
