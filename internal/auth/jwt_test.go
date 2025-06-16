package auth

import (
	"testing"

	"github.com/EliasLd/gotalk-backend/internal/repository"
)

func TestGenerateTokenAndValidateToken_Succes(t *testing.T) {
	user := repository.NewTestUser(t, "testuser")

	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("Unexpected error generating token: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Unexpected error validating token: %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected userID %v, got %v", user.ID, claims.UserID)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token"

	_, err := ValidateToken(invalidToken)
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
}
