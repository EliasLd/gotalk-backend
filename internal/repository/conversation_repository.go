package repository

import (
	"context"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Contract for conversation data access
type ConversationRepository interface {
	CreateConversation(ctx context.Context, conv *models.Conversation) error
}

// Concrete implementation
type conversationRepository struct {
	db *pgxpool.Pool
}

// Constructor
func NewConversationRepository(db *pgxpool.Pool) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) CreateConversation(ctx context.Context, conv *models.Conversation) error {
	if conv.ID == uuid.Nil {
		conv.ID = uuid.New()
	}

	if conv.CreatedAt.IsZero() {
		conv.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO conversations (id, is_public, name, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query,
		conv.ID,
		conv.IsPublic,
		conv.Name,
		conv.CreatedAt,
	)

	return err
}

