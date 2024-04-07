package main

import (
	"context"
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testing"
)

var (
	ctx    = context.Background()
	dbName = "test_db"
)

func TestAddBookExistingBook(t *testing.T) {
	host, port, err := testWithPostgres(t)
	if err != nil {
		t.Fatal(err)
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s password=password dbname=%s sslmode=disable", host, port, "postgres", dbName)
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	mockBook := Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	mockUser := User{Name: "John Doe"}
	mockUser.ID = 1

	cfg := &config{
		Database: db,
	}

	book, _ := cfg.AddBook(mockBook.ISBN, int(mockUser.ID))

	assert.Equal(t, mockBook, book, "Added book should match existing book")
}

func testWithPostgres(t *testing.T) (host string, portNumber int, error error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       dbName,
		},
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatal("Could not start postgres")
	}

	host, err = postgres.Host(ctx)
	if err != nil {
		return "", 0, err
	}

	port, err := postgres.MappedPort(ctx, "5432")
	if err != nil {
		return "", 0, err
	}

	defer func() {
		if err := postgres.Terminate(ctx); err != nil {
			log.Fatal("Could not stop postgres:")
		}
	}()
	return host, port.Int(), nil
}
