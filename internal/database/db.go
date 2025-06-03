package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Global PostgreSQL connection pool.
var DB *pgxpool.Pool

// Initialize a PostgreSQL connection pool
func Connect() error {

	dsn := os.Getenv("DATABASE_URL")

	// Set a timeout to avoid hanging
	// on slow or unrecheable database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("config parsing error: %w", err)
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	DB = dbpool
	fmt.Println("Successfully connected to PostgreSQL")
	return nil
}

// Shuts down the connection pool cleanly
func Close() {
	if DB != nil {
		DB.Close()
		fmt.Println("Closed connection to PostgreSQL")
	}
}
