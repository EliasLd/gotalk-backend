package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/auth"
	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/google/uuid"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	userID := uuid.New()

	user := &models.User {
		ID:		userID,
		Username:	"testuser_authmiddleware",
		Password:	"ValidPswrd123!",
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
	}

	// Generate a valid jwt token
	token, err := auth.GenerateToken(user)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUserID, ok := UserIDFromContext(r.Context())
		if !ok || ctxUserID != userID.String() {
			t.Errorf("Expected userID %s, got %s", userID.String(), ctxUserID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Apply auth middleware
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	rr := httptest.NewRecorder()
	handler := AuthMiddleware(protectedHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r* http.Request) {
		t.Errorf("Handler should not be called with invalide token")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	rr := httptest.NewRecorder()
	handler := AuthMiddleware(protectedHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401 Unauthorized, got %d", rr.Code)
	}
}
