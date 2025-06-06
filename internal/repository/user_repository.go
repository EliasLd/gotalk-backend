package repository

import (
	"context"
	
	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Contract for any kind of user data access implementation.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

// Concrete implementation of UserRepository
type userRepository struct {
	db *pgxpool.Pool
}

// Constructor, returns a new instance of the repository
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// Insert a new user into the database.
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// Retrieves a user by its username
func (r *userRepository) GetUserByUsername(ctx context.Context, username string ) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
