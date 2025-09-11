package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID `json:"id"`
	IsPublic  bool      `json:"isPublic"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}
