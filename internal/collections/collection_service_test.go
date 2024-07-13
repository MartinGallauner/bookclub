package collections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/martingallauner/bookclub/internal"
	"github.com/martingallauner/bookclub/internal/client"
	"github.com/martingallauner/bookclub/internal/server"
	"github.com/martingallauner/bookclub/internal/repository"
	"github.com/martingallauner/bookclub/test"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"context"
	"log"
	"time"
)

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) { //TODO: consider removing that helper and assert inline
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func setupTest() (*server.BookclubServer, error) {
	container, err := CreatePostgresContainer()
	if err != nil {
		log.Fatal(err)
	}

	db, err := internal.SetupDatabaseWithDSN(container.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	client := client.Client{}
	bookRepository := &repository.PostgresBookRepository{Database: db}
	userRepository := &repository.PostgresUserRepository{Database: db}
	linkRepository := &repository.PostgresLinkRepository{Database: db}
	collectionService := New(userRepository, bookRepository, linkRepository, client)

	s := server.New(client, bookRepository, userRepository, linkRepository, nil, nil, collectionService, nil)
	return s, err
}

var (
	ctx = context.Background()
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
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

// Tests if a saved book can be added to the collection of an existing user.
func TestAddBookToUser(t *testing.T) {
	//given
	s, err := setupTest()

	mockBook := internal.Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	s.BookRepository.Save(mockBook)
	mockUser := &internal.User{Name: "Test User"}
	mockUser.ID = 1
	err = s.UserRepository.Save(mockUser)
	if err != nil {
		return
	}
	//when
	request, _ := http.NewRequest(http.MethodPost, "/api/collections", strings.NewReader(`{"user_id": 1, "isbn": "1234567890"}`))
	request.Header.Add("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	assertStatus(t, response.Code, http.StatusOK)
	user, err := s.UserRepository.Get(1)

	if len(user.Books) == 0 {
		t.Errorf("User has no books")
		t.FailNow()
	}
	addedBook := user.Books[0]
	assert.Equal(t, addedBook, mockBook, "Added book should match existing book")
}

// Add Book to an unknown user
func TestAddBookToUnknownUser(t *testing.T) {
	//given
	s, err := setupTest()
	mockBook := internal.Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	s.BookRepository.Save(mockBook)

	if err != nil {
		return
	}

	//when
	request, _ := http.NewRequest(http.MethodPost, "/api/collections", strings.NewReader(`{"user_id": 99, "isbn": "1234567890"}`))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	assertStatus(t, response.Code, http.StatusBadRequest)
}


// Tests search function when one user has the book
func TestSearchBookInNetwork(t *testing.T) {
	//given
	s, err := setupTest()

	userWithBook, err := s.UsersService.CreateUser("Book Owner", "owner@gmail.com")
	book := &internal.Book{ISBN: "1234567890", URL: "https://...", Title: "Test Book"}
	s.BookRepository.Save(*book)
	_, err = s.CollectionService.AddBookToCollection(book.ISBN, userWithBook.ID)

	userWithoutBooks, err := s.UsersService.CreateUser("Reader", "reader@gmail.com")

	s.UsersService.LinkUsers(userWithBook.ID, userWithoutBooks.ID)
	s.UsersService.LinkUsers(userWithoutBooks.ID, userWithBook.ID)

	link, err := s.LinkRepository.Get(userWithBook.ID, userWithoutBooks.ID)
	fmt.Println(link)

	if err != nil {
		t.Fatalf("Unable to setup test.")
	}

	//when
	
	requestBody := server.CollectionService.SearchRequest{UserId: uint(1), ISBN: "1234567890"}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/search", bytes.NewReader(jsonBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got server.SearchResponse
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}

	assert.Equal(t, got.Isbn, "1234567890", "Did we search for the wrong book?")
	assert.Equal(t, len(got.Users), 1, ". Expected a different number of book owners.")
	assertStatus(t, response.Code, http.StatusOK)
}
