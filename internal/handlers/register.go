package handlers

import (
	"encoding/json"
	"net/http"
	"errors"

	appErr "github.com/EliasLd/gotalk-backend/internal/service/errors"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID		string	`json:"id"`
	Username	string	`json:"username"`
	CreatedAt	string	`json:"createdAt"`
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	
	user, err := h.userService.RegisterUser(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrUserAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, appErr.ErrPasswordTooShort),
	    		errors.Is(err, appErr.ErrPasswordMissingDigit),
	     		errors.Is(err, appErr.ErrPasswordMissingUpper),
	     		errors.Is(err, appErr.ErrPasswordMissingLower),
	     		errors.Is(err, appErr.ErrPasswordMissingSymbol):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := registerResponse{
		ID:		user.ID.String(),
		Username:	user.Username,
		CreatedAt:	user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
