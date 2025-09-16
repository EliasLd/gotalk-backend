package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	//	apphttp "github.com/EliasLd/gotalk-backend/internal/http"
	"github.com/EliasLd/gotalk-backend/internal/auth"
	"github.com/EliasLd/gotalk-backend/internal/handlers"
	"github.com/EliasLd/gotalk-backend/internal/repository"
	"github.com/EliasLd/gotalk-backend/internal/service"
	"github.com/google/uuid"
)

func TestGetMeRoute(t *testing.T) {
	db := repository.SetupTest(t)
	repo := repository.NewUserRepository(db)
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
	db := repository.SetupTest(t)
	repo := repository.NewUserRepository(db)
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
	db := repository.SetupTest(t)
	repo := repository.NewUserRepository(db)
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
	db := repository.SetupTest(t)
	repo := repository.NewUserRepository(db)
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

func TestRegisterRoute_InvalidPassword(t *testing.T) {
	db := repository.SetupTest(t)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	handler := handlers.NewHandler(userService)
	router := NewRouter(handler)

	username := "testuser_invalid_password"
	invalidPassword := "abc"

	reqBody := `{"username":"` + username + `","password":"` + invalidPassword + `"}`

	req := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "password") {
		t.Errorf("Expected error message to mention 'password', got: %s", rr.Body.String())
	}
}

func TestLoginRoute(t *testing.T) {
	db := repository.SetupTest(t)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	handler := handlers.NewHandler(userService)
	router := NewRouter(handler)

	username := "testuser_login"
	password := "ValidPasswd123!"

	user, err := userService.RegisterUser(context.Background(), username, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	defer repository.CleanUpUser(t, user.ID, repo)

	reqBody := `{"username":"` + username + `","password":"` + password + `"}`

	req := httptest.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", rr.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check if token exists
	token, ok := response["token"].(string)
	if !ok || token == "" {
		t.Errorf("Expected non-empty token in response")
	}
}

func TestLoginRouteFailures(t *testing.T) {
	db := repository.SetupTest(t)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	handler := handlers.NewHandler(userService)
	router := NewRouter(handler)

	username := "failing_user"
	password := "ValidPasswd123!"
	user, err := userService.RegisterUser(context.Background(), username, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	defer repository.CleanUpUser(t, user.ID, repo)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "Unknown user",
			body:           `{"username":"unknown","password":"doesntmatter"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong password",
			body:           `{"username":"` + username + `","password":"WrongPassword1!"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Malformed JSON",
			body:           `{"username": "badjson", "password": }`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing fields",
			body:           `{"username": "no_password"}`,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("[%s] Expected status %d, got %d", tt.name, tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestUpdateMeRoute_Username(t *testing.T) {
	db := repository.SetupTest(t)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	handler := handlers.NewHandler(userService)
	router := NewRouter(handler)

	username := "testuser_update"
	password := "ValidPasswd123!"
	user, err := userService.RegisterUser(context.Background(), username, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	defer repository.CleanUpUser(t, user.ID, repo)

	// Authenticate user
	token, err := auth.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	newUsername := "updateduser"
	reqBody := `{"username":"` + newUsername + `"}`

	req := httptest.NewRequest("PUT", "/me/update", strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", rr.Code)
	}

	// Check response validity
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["username"] != newUsername {
		t.Errorf("Expected updated username %s, got %s", newUsername, response["username"])
	}
}

func TestUpdateMeRoute_Password(t *testing.T) {
	db := repository.SetupTest(t)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	handler := handlers.NewHandler(userService)
	router := NewRouter(handler)

	username := "testuser_update_pwd"
	oldPassword := "ValidPasswd123!"
	user, err := userService.RegisterUser(context.Background(), username, oldPassword)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	defer repository.CleanUpUser(t, user.ID, repo)

	// Authenticate user
	token, err := auth.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	newPassword := "NewValidPasswd456!"
	reqBody := `{"password":"` + newPassword + `"}`

	req := httptest.NewRequest("PUT", "/me/update", strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", rr.Code)
	}

	// Check response validity
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["username"] != username {
		t.Errorf("Username should not have changed, expected %s, got %s", username, response["username"])
	}

	// Now try logging in with the new password
	loggedInUser, err := userService.AuthenticateUser(context.Background(), username, newPassword)
	if err != nil {
		t.Fatalf("Failed to authenticate with new password: %v", err)
	}

	if loggedInUser.ID != user.ID {
		t.Errorf("Authenticated user ID mismatch, expected %s, got %s", user.ID, loggedInUser.ID)
	}
}
