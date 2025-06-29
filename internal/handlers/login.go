package handlers

import (
	"encoding/json"
	"net/http"

	appErr "github.com/EliasLd/gotalk-backend/internal/service/errors"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}
