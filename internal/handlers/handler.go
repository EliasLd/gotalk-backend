package handlers

import "github.com/EliasLd/gotalk-backend/internal/service"

type Handler struct {
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler {
		userService:	userService,
	}
}
