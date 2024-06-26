package internal

import (
	"context"
	"github.com/markbates/goth"
	client "github.com/martingallauner/bookclub/internal/client"
	repository "github.com/martingallauner/bookclub/internal/repository"
	server "github.com/martingallauner/bookclub/internal/server"
	internal "github.com/martingallauner/bookclub/internal"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"net/http"
	"testing"
	"time"
)

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

// helper method to run for each test //TODO: please don't start a new container for each test
func setupTest() (*server.BookclubServer, error) {
	container, err := CreatePostgresContainer()
	if err != nil {
		log.Fatal(err)
	}

	db, err := internal.SetupDatabaseWithDSN(container.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewBookclubServer(client.Client{}, &repository.PostgresBookRepository{Database: db}, &repository.PostgresUserRepository{Database: db}, &repository.PostgresLinkRepository{Database: db}, &MockAuthService{}, &MockJwtService{})
	return s, err
}

type MockAuthService struct{}

func (svc *MockAuthService) CompleteUserAuth(w http.ResponseWriter, r *http.Request) (goth.User, error) {

	user := goth.User{Name: "Alfred", Email: "alfred@gmail.com"}
	return user, nil
}

type MockJwtService struct{}

func (svc *MockJwtService) CreateToken(issuer string, id int) (string, error) {
	return "mock token", nil
}

var (
	ctx = context.Background()
)

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
