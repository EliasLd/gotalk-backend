package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

//	apphttp "github.com/EliasLd/gotalk-backend/internal/http"
	"github.com/EliasLd/gotalk-backend/internal/auth"
	"github.com/EliasLd/gotalk-backend/internal/handlers"
	"github.com/EliasLd/gotalk-backend/internal/service"
	"github.com/EliasLd/gotalk-backend/internal/repository"
	"github.com/google/uuid"

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

	defer repository.CleanUpUser(t, user.ID, repo)

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

func TestGetMe_Unauthorized(t *testing.T) {
	repo := repository.SetupTest(t)
	userService := service.NewUserService(repo)
	handler := handlers.NewHandler(userService)
	router := NewRouter(handler)

	req := httptest.NewRequest("GET", "/me", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Missing or invalid Authorization header") {
		t.Errorf("Unexpected response body: %s", rr.Body.String())
	}
}

func TestRegisterRoute(t *testing.T) {
	repo := repository.SetupTest(t)
	userService := service.NewUserService(repo)
	handler := handlers.NewHandler(userService)
	router := NewRouter(handler)

	username := "testuser_register"
	password := "ValidPasswd123!"

	reqBody := `{"username":"` + username + `","password":"` + password + `"}`

	req := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d", rr.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["username"] != username {
		t.Errorf("Expected username %s, got %v", username, response["username"])
	}

	if response["id"] == "" || response["createdAt"] == "" {
		t.Error("Expected non-empty ID and CreatedAt fields")
	}

	userID, err := uuid.Parse(response["id"].(string))
	if err != nil {
		t.Fatalf("Failed to parse user ID: %v", err)
	}

	defer repository.CleanUpUser(t, userID, repo)

}

func TestRegisterRoute_UserAlreadyExists(t *testing.T) {
	repo := repository.SetupTest(t)
	userService := service.NewUserService(repo)
	handler := handlers.NewHandler(userService)
	router := NewRouter(handler)

	username := "testuser_register_duplicate"
	password := "ValidPasswd123!"

	firstBody := `{"username":"` + username + `","password":"` + password + `"}`

	firstReq := httptest.NewRequest("POST", "/register", strings.NewReader(firstBody))
	firstReq.Header.Set("Content-Type", "application/json")
	firstRec := httptest.NewRecorder()
	router.ServeHTTP(firstRec, firstReq)

	if firstRec.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created on first register, got %d", firstRec.Code)
	}

	secondReq := httptest.NewRequest("POST", "/register", strings.NewReader(firstBody))
	secondReq.Header.Set("Content-Type", "application/json")
	secondRec := httptest.NewRecorder()
	router.ServeHTTP(secondRec, secondReq)

	if secondRec.Code != http.StatusConflict {
		t.Fatalf("Expected status 409 Conflict on duplicate register, got %d", secondRec.Code)
	}

	userID := ParseUserIDFromResponse(t, firstRec.Body)

	defer repository.CleanUpUser(t, userID, repo)
}
