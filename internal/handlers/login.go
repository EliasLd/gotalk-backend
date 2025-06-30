package handlers

import (
	"encoding/json"
	"net/http"
	"errors"

	appErr "github.com/EliasLd/gotalk-backend/internal/service/errors"
	"github.com/EliasLd/gotalk-backend/internal/auth"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
	}

	user, err := h.userService.AuthenticateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrInvalidCredentials):
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := loginResponse {
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
