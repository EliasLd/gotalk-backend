package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/EliasLd/gotalk-backend/internal/database"
	"github.com/EliasLd/gotalk-backend/internal/handlers"
	httpHandler "github.com/EliasLd/gotalk-backend/internal/http"
	"github.com/EliasLd/gotalk-backend/internal/service"
	"github.com/EliasLd/gotalk-backend/internal/repository"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, pursuing with system environment variables")
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	userRepo 	:= repository.NewUserRepository(database.DB)
	userService 	:= service.NewUserService(userRepo) 

	handler := handlers.NewHandler(userService)
	router 	:= httpHandler.NewRouter(handler)

	port := os.Getenv("PORT")
	addr := ":" + port
	log.Printf("Server running on http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
