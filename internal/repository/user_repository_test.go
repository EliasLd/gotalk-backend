package repository

import (
	"fmt"
	"context"
	"testing"
	"time"
	"os"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/EliasLd/gotalk-backend/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// .env file should be located in the project's root directory
var env_file_path string = "../../.env"

func TestCreateUser(t *testing.T) {
	if err := godotenv.Load(env_file_path); err != nil {
		t.Fatalf("Failed to read environment variables")
	}
	
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Fatalf("DATABASE_URL env var is empty")
	}
	fmt.Println("Using DATABASE_URL: ", dbURL)

	if err := database.Connect(); err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	repo := NewUserRepository(database.DB)

	// Test user
	user := &models.User{
		ID:		uuid.New(),
		Username:	"testuser",
		Password:	"testpasswordhash",
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
	}

	if err := repo.CreateUser(context.Background(), user); err != nil {
		t.Errorf("Failed to create user: %v", err)
	}
}
