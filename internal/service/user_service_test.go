package service

import (
	"context"
	"testing"

	"github.com/EliasLd/gotalk-backend/internal/repository"
	"github.com/EliasLd/gotalk-backend/internal/service/errors"
)

// Test object used to safely access user repository
type testUserService struct {
	UserService
	repo repository.UserRepository
}

func setupService(t *testing.T) testUserService {
	t.Helper()

	repo := repository.SetupTest(t)

	return testUserService {
		UserService:	NewUserService(repo),
		repo:		repo,
	}
}

func TestRegisterUser_Success(t *testing.T) {
	s := setupService(t)

	username := "testuser_valid"
	password := "ValidPass123!"

	user, err := s.RegisterUser(context.Background(), username, password)

	if err != nil {
		t.Fatalf("Expected no error during user registration, got %v", err)
	}

	if user.Username != username {
		t.Errorf("expected username %s, got %s", username, user.Username)
	}

	defer repository.CleanUpUser(t, user.ID, s.repo)

}

func TestRegisterUser_UsernameAlreadyExists(t *testing.T) {
	s := setupService(t)

	username := "testuser_taken"
	password := "AnotherValid123!"

	existingUser, err := s.RegisterUser(context.Background(), username, password)
	if err != nil {
		t.Fatalf("Expected no error during user registration, got %v", err)
	}

	defer repository.CleanUpUser(t, existingUser.ID, s.repo)

	newUser, err := s.RegisterUser(context.Background(), username, password)
	if err != errors.ErrUserAlreadyExists {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)	
	}
	if newUser != nil {
		t.Errorf("Expected no user, got %v", newUser)
	}
}
