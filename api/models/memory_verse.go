package models

import (
	"time"

	"github.com/google/uuid"
)

type MemoryVerse struct {
	ID        uuid.UUID   `json:"id"`
	UserID    *uuid.UUID  `json:"user_id,omitempty"` // Null for system packs
	PackName  string      `json:"pack_name"`
	Reference string      `json:"reference"`
	Text      string      `json:"text"`
	Version   string      `json:"version"`
	Tags      []string    `json:"tags"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
