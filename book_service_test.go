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

// Tests if a saved book can be linked to an existing user.
func TestLinkBookToUser(t *testing.T) {
	container, err := CreatePostgresContainer()
	if err != nil {
		log.Fatal(err)
	}

	db, err := SetupDatabase(container.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &config{
		BookRepository: &PostgresBookRepository{Database: db},
		UserRepository: &PostgresUserRepository{Database: db},
	}

	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	cfg.BookRepository.Save(mockBook)
	mockUser := User{Name: "Test User"}
	mockUser.ID = 1
	err = cfg.UserRepository.Save(mockUser)
	if err != nil {
		return
	}

	book, _ := cfg.AddBook(mockBook.ISBN, 1)

	assert.Equal(t, mockBook, book, "Added book should match existing book")
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

	cfg := &config{
		BookRepository: &PostgresBookRepository{Database: db},
		UserRepository: &PostgresUserRepository{Database: db},
	}

	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	cfg.BookRepository.Save(mockBook)

	if err != nil {
		return
	}

	_, err = cfg.AddBook(mockBook.ISBN, 1)

	assert.Equal(t, err.Error(), "record not found", "Expecting record not found error")
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
