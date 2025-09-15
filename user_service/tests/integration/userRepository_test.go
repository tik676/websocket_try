//go:build integration
// +build integration

package integration

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"user_service/internal/domain"
	"user_service/internal/infrastructure"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupDB(t *testing.T) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	return db
}

func TestUserRegister(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := infrastructure.NewDB(db)

	input := domain.AuthorizationInput{
		Name:     "testuser",
		Password: "password123",
	}

	user, err := repo.Register(input)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	if user.Name != input.Name {
		t.Errorf("expected %s, got: %s", input.Name, user.Name)
	}

	loggedIn, err := repo.Login(input)
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}

	if loggedIn.ID != user.ID {
		t.Errorf("expected ID %d, got %d", user.ID, loggedIn.ID)
	}
}
