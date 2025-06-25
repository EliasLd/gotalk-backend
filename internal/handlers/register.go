package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/EliasLd/gotalk-backend/internal/service"
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
