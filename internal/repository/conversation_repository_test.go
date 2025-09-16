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
