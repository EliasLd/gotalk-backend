package service

import (
	"context"
	"testing"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/EliasLd/gotalk-backend/internal/repository"
	"github.com/google/uuid"
)

func TestConversationService_CreateAndGetConversation(t *testing.T) {
	db := repository.SetupTest(t)
	repo := repository.NewConversationRepository(db)
	service := NewConversationService(repo)
	ctx := context.Background()

	conv := &models.Conversation{
		IsPublic:  true,
		Name:      "Test Conversation",
		CreatedAt: time.Now(),
	}

	created, err := service.CreateConversation(ctx, conv)
	if err != nil {
		t.Fatalf("unexpected error creating conversation: %v", err)
	}

	if created.ID == uuid.Nil {
		t.Errorf("expected a generated ID, got nil")
	}
	if created.Name != "Test Conversation" {
		t.Errorf("expected name %q, got %q", "Test Conversation", created.Name)
	}

	// Fetch again by ID
	fetched, err := service.GetConversationByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("unexpected error fetching conversation: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("expected ID %v, got %v", created.ID, fetched.ID)
	}
}

func TestConversationService_GetConversationByID_NotFound(t *testing.T) {
	db := repository.SetupTest(t)
	repo := repository.NewConversationRepository(db)
	service := NewConversationService(repo)
	ctx := context.Background()

	_, err := service.GetConversationByID(ctx, uuid.New())
	if err == nil {
		t.Fatalf("expected error when conversation not found, got nil")
	}
}
