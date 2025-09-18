package service

import (
	"context"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/EliasLd/gotalk-backend/internal/repository"
	"github.com/google/uuid"
)

// Contract for conversation business logic
type ConversationService interface {
	CreateConversation(ctx context.Context, conv *models.Conversation) (*models.Conversation, error)
	GetConversationByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error)
}

// Actual implementation
type conversationService struct {
	repo repository.ConversationRepository
}

// Constructor
func NewConversationService(repo repository.ConversationRepository) ConversationService {
	return &conversationService{repo: repo}
}

func (s *conversationService) CreateConversation(ctx context.Context, conv *models.Conversation) (*models.Conversation, error) {
	if err := s.repo.CreateConversation(ctx, conv); err != nil {
		return nil, err
	}
	return conv, nil
}

func (s *conversationService) GetConversationByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	return s.repo.GetConversationByID(ctx, id)
}
