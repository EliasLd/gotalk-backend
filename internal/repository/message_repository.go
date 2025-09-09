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
	GetMessagesByConversation(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error)
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

// Retrieves all messages from a conversation
func (r *messageRepository) GetMessagesByConversation(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, content, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&msg.Content,
			&msg.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

func (r *messageRepository) DeleteMessage(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
