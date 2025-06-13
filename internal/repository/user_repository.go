package repository

import (
	"fmt"
	"context"
	
	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

// Contract for any kind of user data access implementation.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error 
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

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	result , err := r.db.Exec(ctx, query, id)
	
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id: %s", id)
	}

	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
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

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, password_hash = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query, user.Username, user.Password, user.UpdatedAt, user.ID)
	return err
}
