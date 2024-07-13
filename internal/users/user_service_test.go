package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/martingallauner/bookclub/internal"
	"github.com/martingallauner/bookclub/internal/server"
	"github.com/martingallauner/bookclub/internal/client"
	"github.com/martingallauner/bookclub/internal/repository"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"net/http"
	"net/http/httptest"
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

// helper method to run for each test //TODO: please don't start a new container for each test
func SetupTest() (*server.BookclubServer, error) {
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
	usersService := New(userRepository, bookRepository, linkRepository, client)

	s := server.New(client, bookRepository, userRepository, linkRepository, nil, nil, nil, usersService)
	return s, err
}

// Tests creation of a new user.
func TestCreateNewUser(t *testing.T) {
	//given
	s, err := SetupTest()
	if err != nil {
		return
	}

	//when
	requestBody := server.CreateUserRequest{Name: "Mocki", Email: "mock@gmail.com"}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(jsonBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got internal.User
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}
	assert.Equal(t, got.Name, "Mocki")
	assert.Equal(t, got.Email, "mock@gmail.com")
	assert.Equal(t, got.ID, uint(1))

	AssertStatus(t, response.Code, http.StatusOK)
}

// Tests to create link request
func TestRequestLink(t *testing.T) {
	//given
	s, err := SetupTest()
	if err != nil {
		return
	}

	user1, _ := s.UsersService.CreateUser("Alpha", "alpha@gmail.com")
	user2, _ := s.UsersService.CreateUser("Bravo", "bravo@gmail.com")

	if err != nil {
		t.Fatalf("Unable to prepare users needed for the test %v", err)
	}

	//when
	requestBody := server.LinkRequest{SenderId: user1.ID, ReceiverId: user2.ID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/links", bytes.NewReader(jsonBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got server.LinkResponse
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}
	test.AssertStatus(t, response.Code, http.StatusOK)
	assert.Equal(t, got.SenderId, user1.ID)
	assert.Equal(t, got.ReceiverId, user2.ID)
	savedLink, err := s.LinkRepository.Get(user1.ID, user2.ID)
	assert.Equal(t, savedLink.SenderId, user1.ID)
	assert.Equal(t, savedLink.ReceiverId, user2.ID)
}

// Tests to create link request
func TestGetLinks(t *testing.T) {
	//given
	s, err := SetupTest()
	if err != nil {
		return
	}

	user1, _ := s.UsersService.CreateUser("Alpha", "alpha@gmail.com")
	user2, _ := s.UsersService.CreateUser("Bravo", "bravo@gmail.com")
	user3, _ := s.UsersService.CreateUser("Charlie", "bravo@gmail.com")

	_, err = s.UsersService.LinkUsers(user1.ID, user2.ID)
	_, err = s.UsersService.LinkUsers(user2.ID, user3.ID)

	if err != nil {
		t.Fatalf("Unable to prepare users needed for the test %v", err)
	}

	//when
	url := fmt.Sprintf("/api/links/%d", user1.ID)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got []server.LinkResponse
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}
	assert.Equal(t, len(got), 1)
	for _, link := range got {
		if link.SenderId == user3.ID || link.ReceiverId == user3.ID {
			t.Errorf("User 1 is not supposed to have any link to user 3")
		}
	}
	test.AssertStatus(t, response.Code, http.StatusOK)
}

// Tests to accept an existing link request
func TestAcceptLink(t *testing.T) {
	//given
	s, err := SetupTest()
	if err != nil {
		return
	}

	user1, _ := s.UsersService.CreateUser("Alpha", "alpha@gmail.com")
	user2, _ := s.UsersService.CreateUser("Bravo", "bravo@gmail.com")
	_, err = s.UsersService.LinkUsers(user1.ID, user2.ID)

	if err != nil {
		t.Fatalf("Unable to prepare users needed for the test %v", err)
	}

	//when
	requestBody := server.LinkRequest{SenderId: user2.ID, ReceiverId: user1.ID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/links", bytes.NewReader(jsonBody))
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got server.LinkResponse
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}
	test.AssertStatus(t, response.Code, http.StatusOK)

	fetchLink, _ := s.LinkRepository.Get(user1.ID, user2.ID)

	//persisted
	assert.Equal(t, fetchLink.AcceptedAt.After(fetchLink.CreatedAt), true)

	//response
	assert.Equal(t, got.IsLinked, true)
}

// Tests the login of a known user
func TestLogin(t *testing.T) {
	//given
	s, err := SetupTest()
	if err != nil {
		return
	}

	testUser := &internal.User{Name: "Alfred", Email: "alfred@gmail.com"}
	err = s.UserRepository.Save(testUser)
	extraUser := &internal.User{Name: "Bert", Email: "bert@gmail.com"}
	err = s.UserRepository.Save(extraUser)

	//when

	request, _ := http.NewRequest(http.MethodPost, "/api/auth/login", nil)
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got server.LoginResponse
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}
	test.AssertStatus(t, response.Code, http.StatusOK)

	//persisted
	persistedUser, err := s.UserRepository.GetByEmail(testUser.Email)
	assert.Equal(t, persistedUser.Name, "Alfred")
	assert.Equal(t, persistedUser.Email, "alfred@gmail.com")

	//response
	assert.Equal(t, got.Name, "Alfred")
	assert.Equal(t, got.Email, "alfred@gmail.com")
	assert.Equal(t, got.Jwt, "mock token")
}

// Tests the login of a UNKNOWN user
func TestLoginOfNewUser(t *testing.T) {
	//given
	s, err := SetupTest()
	if err != nil {
		return
	}

	//when
	request, _ := http.NewRequest(http.MethodPost, "/api/auth/login", nil)
	response := httptest.NewRecorder()
	s.ServeHTTP(response, request)

	//then
	var got server.LoginResponse
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
	}
	test.AssertStatus(t, response.Code, http.StatusOK)

	//persisted
	persistedUser, err := s.UserRepository.GetByEmail("alfred@gmail.com")
	assert.Equal(t, persistedUser.Name, "Alfred")
	assert.Equal(t, persistedUser.Email, "alfred@gmail.com")

	//response
	assert.Equal(t, got.Name, "Alfred")
	assert.Equal(t, got.Email, "alfred@gmail.com")
	assert.Equal(t, got.Jwt, "mock token")
}
