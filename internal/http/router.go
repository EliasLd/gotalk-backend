package http

import (
	"net/http"

	"github.com/EliasLd/gotalk-backend/internal/handlers"
	"github.com/EliasLd/gotalk-backend/internal/http/middleware"
)

func NewRouter(handler * handlers.Handler) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/register", handler.HandleRegister)
	mux.HandleFunc("/login", handler.HandleLogin)

	// Private routes
	mux.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(handler.HandleGetMe)))

	return mux
}
