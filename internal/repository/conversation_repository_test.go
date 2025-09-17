package repository

import (
	"context"
	"testing"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/google/uuid"
)

func TestConversationRepository(t *testing.T) {
	db := SetupTest(t)
	repo := NewConversationRepository(db)

	tests := []struct {
		name     string
		action   func(t *testing.T, repo ConversationRepository) error
		validate func(t *testing.T, repo ConversationRepository, conv *models.Conversation)
	}{
		{
			name: "CreateConversation succeeds",
			action: func(t *testing.T, repo ConversationRepository) error {
				conv := &models.Conversation{
					ID:        uuid.New(),
					IsPublic:  true,
					Name:      "Public Conv",
					CreatedAt: time.Now(),
				}
				return repo.CreateConversation(context.Background(), conv)
			},
			validate: func(t *testing.T, repo ConversationRepository, _ *models.Conversation) {
			},
		},
		{
			name: "GetConversationByID returns correct conversation",
			action: func(t *testing.T, repo ConversationRepository) error {
				conv := &models.Conversation{
					ID:        uuid.New(),
					IsPublic:  false,
					Name:      "Private Conv",
					CreatedAt: time.Now(),
				}
				if err := repo.CreateConversation(context.Background(), conv); err != nil {
					return err
				}
				return nil
			},
			validate: func(t *testing.T, repo ConversationRepository, _ *models.Conversation) {
				conv := &models.Conversation{
					ID:        uuid.New(),
					IsPublic:  false,
					Name:      "To Fetch",
					CreatedAt: time.Now(),
				}
				if err := repo.CreateConversation(context.Background(), conv); err != nil {
					t.Fatalf("Failed to create test conversation: %v", err)
				}
				defer repo.DeleteConversation(context.Background(), conv.ID)

				fetched, err := repo.GetConversationByID(context.Background(), conv.ID)
				if err != nil {
					t.Fatalf("Expected no error fetching conversation, got: %v", err)
				}
				if fetched.Name != conv.Name {
					t.Errorf("Expected Name %v, got %v", conv.Name, fetched.Name)
				}
			},
		},
		{
			name: "DeleteConversation removes conversation",
			action: func(t *testing.T, repo ConversationRepository) error {
				conv := &models.Conversation{
					ID:        uuid.New(),
					IsPublic:  true,
					Name:      "To Delete",
					CreatedAt: time.Now(),
				}
				if err := repo.CreateConversation(context.Background(), conv); err != nil {
					return err
				}
				return repo.DeleteConversation(context.Background(), conv.ID)
			},
			validate: func(t *testing.T, repo ConversationRepository, _ *models.Conversation) {
				conv := &models.Conversation{
					ID:        uuid.New(),
					IsPublic:  true,
					Name:      "To Delete Validate",
					CreatedAt: time.Now(),
				}
				if err := repo.CreateConversation(context.Background(), conv); err != nil {
					t.Fatalf("Failed to create conversation: %v", err)
				}
				if err := repo.DeleteConversation(context.Background(), conv.ID); err != nil {
					t.Fatalf("Failed to delete conversation: %v", err)
				}

				_, err := repo.GetConversationByID(context.Background(), conv.ID)
				if err == nil {
					t.Errorf("Expected error fetching deleted conversation, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.action(t, repo)
			if err != nil {
				t.Fatalf("Action failed: %v", err)
			}
			if tt.validate != nil {
				tt.validate(t, repo, nil)
			}
		})
	}
}
