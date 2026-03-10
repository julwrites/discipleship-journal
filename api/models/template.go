package models

import (
	"time"

	"github.com/google/uuid"
)

type StudyTemplate struct {
	ID                uuid.UUID              `json:"id"`
	CreatorID         uuid.UUID              `json:"creator_id"`
	Title             string                 `json:"title"`
	Description       string                 `json:"description"`
	Structure         map[string]interface{} `json:"structure"` // JSONB
	Prompts           map[string]interface{} `json:"prompts"`   // JSONB
	Fields            []TemplateField        `json:"fields"`    // JSONB
	IsPublic          bool                   `json:"is_public"`
	BibleReferences   []string               `json:"bible_references"`
	AllowUserPassages bool                   `json:"allow_user_passages"`
	TemplateBody      string                 `json:"template_body"`
	RequiredVersion   string                 `json:"required_version"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type TemplateField struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Type        string `json:"type"` // text, textarea, select
	Placeholder string `json:"placeholder,omitempty"`
}

type GenerateRequest struct {
	Inputs       map[string]string `json:"inputs"`
	UserPassages []string          `json:"user_passages,omitempty"`
	UserVersion  string            `json:"user_version,omitempty"`
}
