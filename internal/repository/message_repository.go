package repository

import (
	"context"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Contract for messag data access
type MessageRepository interface {
	CreateMessage(ctx context.Context, message *models.Message) error
}

// Concrete implementation
type messageRepository struct {
	db *pgxpool.Pool
}

// Constructor
func NewMessageRepository(db *pgxpool.Pool) MessageRepository {
	return &messageRepository{db: db}
}

// Inserts a new message
func (r *messageRepository) CreateMessage(ctx context.Context, message *models.Message) error {
	if message.ID == uuid.Nil {
		message.ID = uuid.New()
	}
	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO messages (id, conversation_id, sender_id, content, created_at)
		VALUES {$1, $2, $3, $4, $5}
	`

	_, err := r.db.Exec(ctx, query,
		message.ID,
		message.ConversationID,
		message.SenderID,
		message.Content,
		message.CreatedAt,
	)

	return err
}
