package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

//	apphttp "github.com/EliasLd/gotalk-backend/internal/http"
	"github.com/EliasLd/gotalk-backend/internal/auth"
	"github.com/EliasLd/gotalk-backend/internal/handlers"
	"github.com/EliasLd/gotalk-backend/internal/service"
	"github.com/EliasLd/gotalk-backend/internal/repository"

)

func TestGetMeRoute(t *testing.T) {
	repo := repository.SetupTest(t)
	userService := service.NewUserService(repo)
	handler := handlers.NewHandler(userService)
	router := NewRouter(handler)
	
	username := "testuser_GetMeRoute"
	password := "ValidPasswd123!"
	user, err := userService.RegisterUser(nil, username, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	// Decode answer
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["username"] != username {
		t.Errorf("Expected username %s, got %s", username, response["username"])
	}
}
