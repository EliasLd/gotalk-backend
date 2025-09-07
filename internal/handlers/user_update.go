package handlers

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/EliasLd/gotalk-backend/internal/service"
	appErr "github.com/EliasLd/gotalk-backend/internal/service/errors"
	"github.com/EliasLd/gotalk-backend/internal/http/middleware"
	"github.com/google/uuid"
)

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

type UpdateUserResponse struct {
	ID		string `json:"id"`
	Username	string `json:"username"`
	UpdatedAt	string `json:"updatedAt"`
}

// User to gracefully handle user data updates
func (h *Handler) HandleUpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == nil && req.Password == nil {
		http.Error(w, "At least one field (username of password) must be provided", http.StatusBadRequest)
		return
	}
	
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
	}
	updatedUser, err := h.userService.UpdateUser(r.Context(), userUUID, service.UpdateUserInput { 
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound)
		case errors.Is(err, appErr.ErrPasswordTooShort),
	    		errors.Is(err, appErr.ErrPasswordMissingDigit),
	     		errors.Is(err, appErr.ErrPasswordMissingUpper),
	     		errors.Is(err, appErr.ErrPasswordMissingLower),
	     		errors.Is(err, appErr.ErrPasswordMissingSymbol),
			errors.Is(err, appErr.ErrPasswordHashingFailed):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := UpdateUserResponse {
		ID:		updatedUser.ID.String(),
		Username:	updatedUser.Username,
		UpdatedAt:	updatedUser.UpdatedAt.Format("2006-01-02T15:04:05z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
