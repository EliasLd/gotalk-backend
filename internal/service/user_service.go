package service

import (
	"context"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/EliasLd/gotalk-backend/internal/repository"
	"github.com/EliasLd/gotalk-backend/internal/service/errors"
	"github.com/google/uuid"
)

// Defines business logic operations related to users.
type UserService interface {
	RegisterUser(ctx context.Context, username, password string) (*models.User, error)
}

// Concrete implementation of UserService.
type userService struct {
	repo repository.UserRepository
}

// Creates a new UserService instance.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, username, password string) (*models.User, error) {
	existingUser, err := s.repo.GetUserByUsername(context.Background(), username)
	if err == nil && existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// TODO: Add password hash here
	hashedPassword := password

	user := &models.User {
		ID:		uuid.New(),
		Username:	username,
		Password:	hashedPassword,
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
	}

	if err := s.repo.CreateUser(context.Background(), user); err != nil {
		return nil, err
	}

	return user, nil
}


