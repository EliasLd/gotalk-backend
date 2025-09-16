package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/database"
	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// This file contains a collection of helper functions
// intended to be used in every db/repo/business tests

// .env file should be located in the project's root directory
var env_file_path string = "../../.env"

// Helper function used to prepare the testing environment for user repository
func SetupTest(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if err := godotenv.Load(env_file_path); err != nil {
		t.Fatalf("Failed to load environment variables from %s", env_file_path)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Fatalf("DATABASE_URL env var is empty")
	}
	fmt.Println("Using DATABASE_URL: ", dbURL)

	if err := database.Connect(); err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	return database.DB
}

// Helper function used to clean database by deleting a newly added user
func CleanUpUser(t *testing.T, id uuid.UUID, repo UserRepository) {
	t.Logf("Now deleting the newly added user...")
	err := repo.DeleteUser(context.Background(), id)
	if err != nil {
		t.Logf("Warning: failed to clean up user: %v", err)
	}
}

// Helper function used to create a new user for each test
func NewTestUser(t *testing.T, username string) *models.User {
	t.Helper()
	return &models.User{
		ID:        uuid.New(),
		Username:  username,
		Password:  "some_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
