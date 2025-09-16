package service

import (
	"context"
	"testing"

	"github.com/EliasLd/gotalk-backend/internal/repository"
	"github.com/EliasLd/gotalk-backend/internal/service/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Test object used to safely access user repository
type testUserService struct {
	UserService
	repo repository.UserRepository
}

func setupService(t *testing.T) testUserService {
	t.Helper()

	db := repository.SetupTest(t)
	repo := repository.NewUserRepository(db)

	return testUserService{
		UserService: NewUserService(repo),
		repo:        repo,
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

// Table-driven tests to check all password validation errors
func TestRegisterUser_InvalidPasswords(t *testing.T) {
	s := setupService(t)
	username := "weakuser"

	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "Too short",
			password: "Aa1!",
			wantErr:  errors.ErrPasswordTooShort,
		},
		{
			name:     "Missing digit",
			password: "Aa!abcdefgh",
			wantErr:  errors.ErrPasswordMissingDigit,
		},
		{
			name:     "Missing uppercase",
			password: "abcde1234!",
			wantErr:  errors.ErrPasswordMissingUpper,
		},
		{
			name:     "Missing lowercase",
			password: "ABCDE1234!",
			wantErr:  errors.ErrPasswordMissingLower,
		},
		{
			name:     "Missing symbol",
			password: "Abcdef1234",
			wantErr:  errors.ErrPasswordMissingSymbol,
		},
	}

	for _, test_case := range tests {
		t.Run(test_case.name, func(t *testing.T) {
			_, err := s.RegisterUser(context.Background(), username+"_"+test_case.name, test_case.password)

			if err == nil {
				t.Fatalf("Expected error %v, got nil", test_case.wantErr)
			}
			if err != test_case.wantErr {
				t.Errorf("Expected error %v, got %v", test_case.wantErr, err)
			}
		})
	}
}

func TestGetUserByID_Success(t *testing.T) {
	s := setupService(t)

	username := "testuser_by_id"
	password := "ValidPass123!"

	user, err := s.RegisterUser(context.Background(), username, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	defer repository.CleanUpUser(t, user.ID, s.repo)

	found, err := s.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if found == nil || found.ID != user.ID {
		t.Errorf("Expected user ID %vn got %v", user.ID, found.ID)
	}
}

func TestGetUserByID_Bad_ID(t *testing.T) {
	s := setupService(t)

	randomID := uuid.New()

	user, err := s.GetUserByID(context.Background(), randomID)

	if err == nil {
		t.Error("Expected an error when retrieving non-existent user, got nil")
	}
	if user != nil {
		t.Errorf("Expeceted no user, got %v", user)
	}
}

func TestGetUserByUsername_Success(t *testing.T) {
	s := setupService(t)

	username := "testuser_by_username"
	password := "ValidPass123!"

	user, err := s.RegisterUser(context.Background(), username, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	defer repository.CleanUpUser(t, user.ID, s.repo)

	found, err := s.GetUserByUsername(context.Background(), username)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if found == nil || found.Username != username {
		t.Errorf("Expected username %v, got %v", username, found.Username)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	s := setupService(t)

	username := "testuser_delete"
	password := "ValidPass123!"

	user, err := s.RegisterUser(context.Background(), username, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	if err := s.DeleteUser(context.Background(), user.ID); err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Confirm user no longer exists
	found, err := s.GetUserByID(context.Background(), user.ID)
	if err == nil && found != nil {
		t.Errorf("Expected error or nil user after deletion, got user: %v", found)
	}
}

func TestDeleteUser_Not_Found(t *testing.T) {
	s := setupService(t)

	randomID := uuid.New()

	err := s.DeleteUser(context.Background(), randomID)

	if err == nil {
		t.Error("Expected an error when deleting non-existent user, got nil")
	}
}

func TestUpdateUser(t *testing.T) {
	s := setupService(t)

	username := "testuser_update"
	password := "ValidPass123!"

	user, err := s.RegisterUser(context.Background(), username, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	defer repository.CleanUpUser(t, user.ID, s.repo)

	username = "testuser_new_user"
	password = "newValidPass123!"

	input := UpdateUserInput{
		Username: &username,
		Password: &password,
	}

	user, err = s.UpdateUser(context.Background(), user.ID, input)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	check_user, err := s.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if check_user.Username != username {
		t.Errorf("Expected username to be updated to %s, got %s", username, check_user.Username)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(check_user.Password), []byte(password)); err != nil {
		t.Errorf("Password hash does not match new password: %v", err)
	}
}
