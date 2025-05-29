package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EliasLd/gotalk-backend/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthHandler)

	port := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
