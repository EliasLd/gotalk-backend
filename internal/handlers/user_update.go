package handlers

import (
	"github.com/EliasLd/gotalk-backend/internal/middleware"	
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
