package main

import (
	"context"
	"github.com/magiconair/properties/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
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

func TestAddBookExistingBook(t *testing.T) {
	container, err := CreatePostgresContainer()
	if err != nil {
		log.Fatal(err)
	}

	db, err := SetupDatabase(container.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	mockUser := User{Name: "John Doe"}
	mockUser.ID = 1
	db.Table("books").Save(mockBook)

	cfg := &config{
		Database:       db,
		BookRepository: &PostgresBookRepository{Database: db},
		UserRepository: &PostgresUserRepository{Database: db},
	}

	book, _ := cfg.AddBook(mockBook.ISBN, int(mockUser.ID))

	assert.Equal(t, mockBook, book, "Added book should match existing book")
}

// todo extract setup
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
