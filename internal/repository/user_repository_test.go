package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
	db := SetupTest(t)
	repo := NewUserRepository(db)

	user := NewTestUser(t, "testuser_create")
	// Clean up at test end
	defer CleanUpUser(t, user.ID, repo)

	if err := repo.CreateUser(context.Background(), user); err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	t.Logf("User added successfully")
}

func TestDeleteUser(t *testing.T) {
	db := SetupTest(t)
	repo := NewUserRepository(db)

	user := NewTestUser(t, "testuser_delete")

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
	db := SetupTest(t)
	repo := NewUserRepository(db)

	user := NewTestUser(t, "testuser_lookup_by_username")

	// Clean up at test end
	defer CleanUpUser(t, user.ID, repo)

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
	db := SetupTest(t)
	repo := NewUserRepository(db)

	user := NewTestUser(t, "testuser_lookup_by_id")

	defer CleanUpUser(t, user.ID, repo)

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
	db := SetupTest(t)
	repo := NewUserRepository(db)

	user1 := NewTestUser(t, "same_name")
	user2 := NewTestUser(t, "same_name")

	defer CleanUpUser(t, user1.ID, repo)
	// Should log a warning because user2 is never created
	defer CleanUpUser(t, user2.ID, repo)

	if err := repo.CreateUser(context.Background(), user1); err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	if err := repo.CreateUser(context.Background(), user2); err == nil {
		t.Errorf("Expected error when adding second user with the same name")
	}
}
