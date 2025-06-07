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

// Helper function used to prepare the testing environment
func setupTest(t *testing.T) UserRepository {
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

	return NewUserRepository(database.DB)

}

// Helper function used to clean database by deleting a newly added user
func cleanUpUser(t *testing.T, id uuid.UUID, repo UserRepository) {
	t.Logf("Now deleting the newly added user...")
	err := repo.DeleteUser(context.Background(), id)
	if err != nil {
		t.Logf("Warning: failed to clean up user: %v", err)
	}
}

// Helper function used to create a new user for each test
func newTestUser(t *testing.T, username string) *models.User {
	t.Helper()
	return &models.User {
		ID:		uuid.New(),
		Username:	username,
		Password:	"some_password",
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
	}
}

func TestCreateUser(t *testing.T) {
	repo := setupTest(t)

	user := newTestUser(t, "testuser_create")	
	// Clean up at test end
	defer cleanUpUser(t, user.ID, repo)

	if err := repo.CreateUser(context.Background(), user); err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	t.Logf("User added successfully")
}

func TestDeleteUser(t *testing.T) {
	repo := setupTest(t)

	user := newTestUser(t, "testuser_delete")

	if err := repo.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if err := repo.DeleteUser(context.Background(), user.ID); err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	_, err := repo.GetUserByUsername(context.Background(), user.Username)
	if err == nil {
		t.Errorf("Expected error when fetching deleted user, got nil")
	}

	nonExistentID := uuid.New()
	err = repo.DeleteUser(context.Background(), nonExistentID)
	if err == nil {
		t.Errorf("Expected error when deleteing non-existent user, got nil")
	}
}

func TestGetUserByUsername(t *testing.T) {
	repo :=  setupTest(t)

	user := newTestUser(t, "testuser_lookup_by_username")

	// Clean up at test end
	defer cleanUpUser(t, user.ID, repo)

	if err := repo.CreateUser(context.Background(), user); err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	found, err := repo.GetUserByUsername(context.Background(), user.Username)
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected user but got nil")
	}

	if found.ID != user.ID {
		t.Errorf("Expected ID %v, got %v", user.ID, found.ID)
	}

	t.Logf("Usernames are matching")
}

func TestGetUserByID(t *testing.T) {
	repo := setupTest(t)

	user := newTestUser(t, "testuser_lookup_by_id")
	
	defer cleanUpUser(t, user.ID, repo)

	if err := repo.CreateUser(context.Background(), user); err != nil {
		t.Errorf("Failed to create user: %v", user)
	}

	found, err := repo.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected user, got nil")
	}

	if found.ID != user.ID {
		t.Errorf("Expected ID %v, got %v", user.ID, found.ID)
	}

	t.Logf("IDs are matching")
}

func TestCreateUserSameUsername(t *testing.T) {
	repo := setupTest(t)

	user1 := newTestUser(t, "same_name")
	user2 := newTestUser(t, "same_name")

	defer cleanUpUser(t, user1.ID, repo)
	// Should log a warning because user2 is never created
	defer cleanUpUser(t, user2.ID, repo)

	if err := repo.CreateUser(context.Background(), user1); err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	if err := repo.CreateUser(context.Background(), user2); err == nil {
		t.Errorf("Expected error when adding second user with the same name")
	}
}
