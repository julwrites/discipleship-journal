package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"discipleship_journal_api/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TemplateService interface {
	CreateTemplate(ctx context.Context, tmpl *models.StudyTemplate) (*models.StudyTemplate, error)
	GetTemplate(ctx context.Context, id uuid.UUID) (*models.StudyTemplate, error)
	ListTemplates(ctx context.Context, userID uuid.UUID) ([]*models.StudyTemplate, error)
	ListPublicTemplates(ctx context.Context) ([]*models.StudyTemplate, error)
	UpdateTemplate(ctx context.Context, tmpl *models.StudyTemplate) error
	DeleteTemplate(ctx context.Context, id, userID uuid.UUID) error
	CloneTemplate(ctx context.Context, id, userID uuid.UUID) (*models.StudyTemplate, error)
	GenerateContent(ctx context.Context, templateID uuid.UUID, inputs map[string]string) (string, error)
}

type templateService struct {
	db      DBInterfaceWithQuery
	aiClient BibleAIClient
}

func NewTemplateService(db DBInterfaceWithQuery, aiClient BibleAIClient) TemplateService {
	return &templateService{db: db, aiClient: aiClient}
}

func (s *templateService) CreateTemplate(ctx context.Context, tmpl *models.StudyTemplate) (*models.StudyTemplate, error) {
	tmpl.ID = uuid.New()
	tmpl.CreatedAt = time.Now()
	tmpl.UpdatedAt = time.Now()

	structureJSON, _ := json.Marshal(tmpl.Structure)
	promptsJSON, _ := json.Marshal(tmpl.Prompts)
	fieldsJSON, _ := json.Marshal(tmpl.Fields)

	query := `
		INSERT INTO study_templates (id, creator_id, title, description, structure, prompts, fields, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	_, err := s.db.Exec(ctx, query,
		tmpl.ID, tmpl.CreatorID, tmpl.Title, tmpl.Description, structureJSON, promptsJSON, fieldsJSON, tmpl.IsPublic, tmpl.CreatedAt, tmpl.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (s *templateService) GetTemplate(ctx context.Context, id uuid.UUID) (*models.StudyTemplate, error) {
	query := `
		SELECT id, creator_id, title, description, structure, prompts, fields, is_public, created_at, updated_at
		FROM study_templates
		WHERE id = $1
	`
	var t models.StudyTemplate
	var structureBytes, promptsBytes, fieldsBytes []byte

	err := s.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.CreatorID, &t.Title, &t.Description, &structureBytes, &promptsBytes, &fieldsBytes, &t.IsPublic, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	_ = json.Unmarshal(structureBytes, &t.Structure)
	_ = json.Unmarshal(promptsBytes, &t.Prompts)
	_ = json.Unmarshal(fieldsBytes, &t.Fields)

	return &t, nil
}

func (s *templateService) ListTemplates(ctx context.Context, userID uuid.UUID) ([]*models.StudyTemplate, error) {
	query := `
		SELECT id, creator_id, title, description, structure, prompts, fields, is_public, created_at, updated_at
		FROM study_templates
		WHERE creator_id = $1
		ORDER BY updated_at DESC
	`
	return s.scanTemplates(ctx, query, userID)
}

func (s *templateService) ListPublicTemplates(ctx context.Context) ([]*models.StudyTemplate, error) {
	query := `
		SELECT id, creator_id, title, description, structure, prompts, fields, is_public, created_at, updated_at
		FROM study_templates
		WHERE is_public = true
		ORDER BY created_at DESC
	`
	return s.scanTemplates(ctx, query)
}

func (s *templateService) scanTemplates(ctx context.Context, query string, args ...any) ([]*models.StudyTemplate, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*models.StudyTemplate
	for rows.Next() {
		var t models.StudyTemplate
		var structureBytes, promptsBytes, fieldsBytes []byte
		if err := rows.Scan(
			&t.ID, &t.CreatorID, &t.Title, &t.Description, &structureBytes, &promptsBytes, &fieldsBytes, &t.IsPublic, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(structureBytes, &t.Structure)
		_ = json.Unmarshal(promptsBytes, &t.Prompts)
		_ = json.Unmarshal(fieldsBytes, &t.Fields)
		templates = append(templates, &t)
	}
	return templates, nil
}

func (s *templateService) UpdateTemplate(ctx context.Context, tmpl *models.StudyTemplate) error {
	structureJSON, _ := json.Marshal(tmpl.Structure)
	promptsJSON, _ := json.Marshal(tmpl.Prompts)
	fieldsJSON, _ := json.Marshal(tmpl.Fields)

	query := `
		UPDATE study_templates
		SET title=$1, description=$2, structure=$3, prompts=$4, fields=$5, is_public=$6, updated_at=NOW()
		WHERE id=$7 AND creator_id=$8
	`
	cmd, err := s.db.Exec(ctx, query,
		tmpl.Title, tmpl.Description, structureJSON, promptsJSON, fieldsJSON, tmpl.IsPublic, tmpl.ID, tmpl.CreatorID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *templateService) DeleteTemplate(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM study_templates WHERE id=$1 AND creator_id=$2`
	cmd, err := s.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *templateService) CloneTemplate(ctx context.Context, id, userID uuid.UUID) (*models.StudyTemplate, error) {
	original, err := s.GetTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	clone := &models.StudyTemplate{
		CreatorID:   userID,
		Title:       original.Title + " (Copy)",
		Description: original.Description,
		Structure:   original.Structure,
		Prompts:     original.Prompts,
		Fields:      original.Fields,
		IsPublic:    false,
	}

	return s.CreateTemplate(ctx, clone)
}

func (s *templateService) GenerateContent(ctx context.Context, templateID uuid.UUID, inputs map[string]string) (string, error) {
	tmpl, err := s.GetTemplate(ctx, templateID)
	if err != nil {
		return "", err
	}

	// Basic Prompt Construction
	// We assume a 'system' key in prompts for the instruction
	systemPrompt, ok := tmpl.Prompts["system"].(string)
	if !ok {
		systemPrompt = "You are a helpful Bible study assistant. Use the user provided inputs to create a study guide."
	}

	userPrompt := "Please generate a study guide based on the following inputs:\n"
	for key, val := range inputs {
		userPrompt += fmt.Sprintf("%s: %s\n", key, val)
	}

	// Add structure guidance
	structureJSON, _ := json.MarshalIndent(tmpl.Structure, "", "  ")
	userPrompt += fmt.Sprintf("\nPlease format the output according to this JSON structure if possible, or use it as a section guide:\n%s", string(structureJSON))

	// Prepare payload for ChatCompletion
	// ChatCompletion signature is: (ctx, payload map[string]interface{})
	payload := map[string]interface{}{
		"prompt": userPrompt,
		"type":   "custom", // We might need to handle custom prompt injection in the client
	}

	// WARNING: The current RealBibleAIClient uses pre-defined system prompts via `type`.
	// To support dynamic templates, we need to bypass the template lookup if a custom system prompt is needed,
	// OR we rely on the `prompt` field being the full user message.
	// Since `RealBibleAIClient.ChatCompletion` logic is:
	// 1. Select template by `type` (default 'ask').
	// 2. If template found, replace {PROMPT}.
	// 3. If no template found, use prompt as is.

	// So, if we pass a `type` that doesn't exist in `SystemPrompts` map, it falls back to raw prompt.
	// However, we want to inject the *System Prompt* from the template.
	// The `BibleAIClient` as written is somewhat rigid around `SystemPrompts` loaded at config time.
	// We might need to modify `ChatCompletion` or append the system instruction to the user prompt.

	// Strategy: Append the template's system prompt to the beginning of the user prompt.
	fullPrompt := fmt.Sprintf("System Instruction: %s\n\nUser Request: %s", systemPrompt, userPrompt)

	payload["prompt"] = fullPrompt
	payload["type"] = "raw_template_generation" // Likely doesn't exist, triggering fallback

	resp, err := s.aiClient.ChatCompletion(ctx, payload)
	if err != nil {
		return "", err
	}

	// Extract text from response
	if text, ok := resp["text"].(string); ok {
		return text, nil
	}

	// Fallback check for nested choice (OpenAI style)
	if choices, ok := resp["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if msg, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := msg["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no content in response")
}
