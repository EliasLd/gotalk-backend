package auth

import (
	"testing"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/repository"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
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

func TestValidateToken_ExpiredToken(t *testing.T) {
	userID := uuid.New()

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}).SignedString(jwtSecret)

	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	_, err = ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateToken_WrongAlgorithm(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodNone, &Claims{
		UserID: uuid.New(),
	})

	signedToken, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("Error signing token with 'none': %v", err)
	}

	_, err = ValidateToken(signedToken)
	if err == nil {
		t.Fatal("Expected error for token with wrong algorithm, got nil")
	}
}
