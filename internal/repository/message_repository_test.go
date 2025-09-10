package repository

import (
	//"context"
	//"testing"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/google/uuid"
)

// Helper function used to create a test message
func newTestMessage(conversationID, senderID uuid.UUID, content string) *models.Message {
	return &models.Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      time.Now(),
	}
}
