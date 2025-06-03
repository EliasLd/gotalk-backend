package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/EliasLd/gotalk-backend/internal/database"
	"github.com/EliasLd/gotalk-backend/internal/handlers"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, pursuing with system environment variables")
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)

	port := os.Getenv("PORT")
	addr := ":" + port
	log.Printf("Server running on http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
