package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Contract for conversation data access
type ConversationRepository interface {
	CreateConversation(ctx context.Context, conv *models.Conversation) error
	GetConversationByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error)
	DeleteConversation(ctx context.Context, id uuid.UUID) error
	ListPublicConversations(ctx context.Context) ([]*models.Conversation, error)
	AddMember(ctx context.Context, conversationID, userID uuid.UUID) error
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

func (r *conversationRepository) GetConversationByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	query := `
		SELECT id, is_public, name, created_at
		FROM conversations
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var conv models.Conversation
	err := row.Scan(&conv.ID, &conv.IsPublic, &conv.Name, &conv.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("Could not get conversation: %w", err)
	}

	return &conv, nil
}

func (r *conversationRepository) DeleteConversation(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM conversations WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *conversationRepository) ListPublicConversations(ctx context.Context) ([]*models.Conversation, error) {
	query := `
		SELECT id, is_public, name, created_at FROM conversations
		WHERE is_public = true
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Could not query public conversations: %w", err)
	}

	var conversations []*models.Conversation
	for rows.Next() {
		var conv models.Conversation
		if err := rows.Scan(&conv.ID, &conv.IsPublic, &conv.Name, &conv.CreatedAt); err != nil {
			return nil, fmt.Errorf("Could not scan conversation: %w", err)
		}
		conversations = append(conversations, &conv)
	}
	return conversations, nil
}

func (r *conversationRepository) AddMember(ctx context.Context, conversationID, userID uuid.UUID) error {
	query := `
		INSERT INTO conversation_members (user_id, conversation_id, joined_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, userID, conversationID)
	if err != nil {
		return fmt.Errorf("Could not add member: %w", err)
	}

	return nil
}
