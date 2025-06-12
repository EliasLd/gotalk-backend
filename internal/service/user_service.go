package service

import (
	"context"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/EliasLd/gotalk-backend/internal/repository"
	"github.com/EliasLd/gotalk-backend/internal/service/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Defines business logic operations related to users.
type UserService interface {
	RegisterUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*models.User, error)
}

// Concrete implementation of UserService.
type userService struct {
	repo repository.UserRepository
}

type UpdateUserInput struct {
	Username *string
	Password *string
}

// Creates a new UserService instance.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s *userService) RegisterUser(ctx context.Context, username, password string) (*models.User, error) {
	existingUser, err := s.repo.GetUserByUsername(context.Background(), username)
	if err == nil && existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, errors.ErrPasswordHashingFailed
	}

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

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Username != nil {
		user.Username = *input.Username
	}

	if input.Password != nil {
		if err := ValidatePassword(*input.Password); err != nil {
			return nil, err
		}

		hashedPassword, err := hashPassword(*input.Password)
		if err != nil {
			return nil, errors.ErrPasswordHashingFailed
		}

		user.Password = hashedPassword
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
