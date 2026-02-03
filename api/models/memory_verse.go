package models

import (
	"time"

	"github.com/google/uuid"
)

type VersePack struct {
	ID          uuid.UUID  `json:"id"`
	UserID      *uuid.UUID `json:"user_id,omitempty"` // Null for system packs
	Title       string     `json:"title"`
	Identifier  string     `json:"identifier,omitempty"`
	Description string     `json:"description,omitempty"`
	IsPublic    bool       `json:"is_public"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	VerseCount  int        `json:"verse_count"` // Computed field
}

type MemoryVerse struct {
	ID          uuid.UUID `json:"id"`
	VersePackID uuid.UUID `json:"verse_pack_id"`
	Reference   string    `json:"reference"`
	Title       string    `json:"title"`
	Version     string    `json:"version"`
	VersionSource string  `json:"version_source,omitempty"` // "override", "user_default", "original"
	Tags        []string  `json:"tags"`
	PackTitle   string    `json:"pack_title,omitempty"` // Populated in searches
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
