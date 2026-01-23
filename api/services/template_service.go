package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
	GenerateContent(ctx context.Context, templateID uuid.UUID, req models.GenerateRequest) (string, error)
}

type templateService struct {
	db       DBInterfaceWithQuery
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
	bibleRefsJSON, _ := json.Marshal(tmpl.BibleReferences)

	query := `
		INSERT INTO study_templates (id, creator_id, title, description, structure, prompts, fields, is_public, bible_references, allow_user_passages, template_body, required_version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`
	_, err := s.db.Exec(ctx, query,
		tmpl.ID, tmpl.CreatorID, tmpl.Title, tmpl.Description, structureJSON, promptsJSON, fieldsJSON, tmpl.IsPublic,
		bibleRefsJSON, tmpl.AllowUserPassages, tmpl.TemplateBody, tmpl.RequiredVersion,
		tmpl.CreatedAt, tmpl.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (s *templateService) GetTemplate(ctx context.Context, id uuid.UUID) (*models.StudyTemplate, error) {
	query := `
		SELECT id, creator_id, title, description, structure, prompts, fields, is_public, bible_references, allow_user_passages, template_body, required_version, created_at, updated_at
		FROM study_templates
		WHERE id = $1
	`
	var t models.StudyTemplate
	var structureBytes, promptsBytes, fieldsBytes, bibleRefsBytes []byte

	err := s.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.CreatorID, &t.Title, &t.Description, &structureBytes, &promptsBytes, &fieldsBytes, &t.IsPublic,
		&bibleRefsBytes, &t.AllowUserPassages, &t.TemplateBody, &t.RequiredVersion,
		&t.CreatedAt, &t.UpdatedAt,
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
	_ = json.Unmarshal(bibleRefsBytes, &t.BibleReferences)

	return &t, nil
}

func (s *templateService) ListTemplates(ctx context.Context, userID uuid.UUID) ([]*models.StudyTemplate, error) {
	query := `
		SELECT id, creator_id, title, description, structure, prompts, fields, is_public, bible_references, allow_user_passages, template_body, required_version, created_at, updated_at
		FROM study_templates
		WHERE creator_id = $1
		ORDER BY updated_at DESC
	`
	return s.scanTemplates(ctx, query, userID)
}

func (s *templateService) ListPublicTemplates(ctx context.Context) ([]*models.StudyTemplate, error) {
	query := `
		SELECT id, creator_id, title, description, structure, prompts, fields, is_public, bible_references, allow_user_passages, template_body, required_version, created_at, updated_at
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
		var structureBytes, promptsBytes, fieldsBytes, bibleRefsBytes []byte
		if err := rows.Scan(
			&t.ID, &t.CreatorID, &t.Title, &t.Description, &structureBytes, &promptsBytes, &fieldsBytes, &t.IsPublic,
			&bibleRefsBytes, &t.AllowUserPassages, &t.TemplateBody, &t.RequiredVersion,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(structureBytes, &t.Structure)
		_ = json.Unmarshal(promptsBytes, &t.Prompts)
		_ = json.Unmarshal(fieldsBytes, &t.Fields)
		_ = json.Unmarshal(bibleRefsBytes, &t.BibleReferences)
		templates = append(templates, &t)
	}
	return templates, nil
}

func (s *templateService) UpdateTemplate(ctx context.Context, tmpl *models.StudyTemplate) error {
	structureJSON, _ := json.Marshal(tmpl.Structure)
	promptsJSON, _ := json.Marshal(tmpl.Prompts)
	fieldsJSON, _ := json.Marshal(tmpl.Fields)
	bibleRefsJSON, _ := json.Marshal(tmpl.BibleReferences)

	query := `
		UPDATE study_templates
		SET title=$1, description=$2, structure=$3, prompts=$4, fields=$5, is_public=$6,
		bible_references=$7, allow_user_passages=$8, template_body=$9, required_version=$10,
		updated_at=NOW()
		WHERE id=$11 AND creator_id=$12
	`
	cmd, err := s.db.Exec(ctx, query,
		tmpl.Title, tmpl.Description, structureJSON, promptsJSON, fieldsJSON, tmpl.IsPublic,
		bibleRefsJSON, tmpl.AllowUserPassages, tmpl.TemplateBody, tmpl.RequiredVersion,
		tmpl.ID, tmpl.CreatorID,
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
		CreatorID:         userID,
		Title:             original.Title + " (Copy)",
		Description:       original.Description,
		Structure:         original.Structure,
		Prompts:           original.Prompts,
		Fields:            original.Fields,
		IsPublic:          false,
		BibleReferences:   original.BibleReferences,
		AllowUserPassages: original.AllowUserPassages,
		TemplateBody:      original.TemplateBody,
		RequiredVersion:   original.RequiredVersion,
	}

	return s.CreateTemplate(ctx, clone)
}

func (s *templateService) GenerateContent(ctx context.Context, templateID uuid.UUID, req models.GenerateRequest) (string, error) {
	tmpl, err := s.GetTemplate(ctx, templateID)
	if err != nil {
		return "", err
	}

	// 1. Resolve Version
	version := "ESV" // Default
	if tmpl.RequiredVersion != "" {
		version = tmpl.RequiredVersion
	} else if req.UserVersion != "" {
		version = req.UserVersion
	}

	// 2. Resolve Passages
	var references []string
	if tmpl.AllowUserPassages {
		references = req.UserPassages
	} else {
		references = tmpl.BibleReferences
	}

	// 3. Fetch Passages Text (for output) and format references (for prompt)
	var passageTexts []string
	for _, ref := range references {
		res, err := s.aiClient.GetPassage(ctx, ref, version)
		if err != nil {
			// Log error but continue
			passageTexts = append(passageTexts, fmt.Sprintf("%s: (Error fetching text)", ref))
			continue
		}

		// Extract text
		text := ""
		if t, ok := res["text"].(string); ok {
			text = t
		} else if t, ok := res["verse"].(string); ok {
			text = t
		}

		passageTexts = append(passageTexts, fmt.Sprintf("<blockquote><p><strong>%s (%s)</strong></p>%s</blockquote>", ref, version, text))
	}
	passagesBlock := strings.Join(passageTexts, "\n\n")
	referencesBlock := strings.Join(references, ", ")

	// 4. Construct AI Prompt
	globalPrompt := s.aiClient.GetSystemPrompt("system")
	if globalPrompt == "" {
		globalPrompt = "You are a helpful Bible study assistant."
	}

	templatePromptRaw, _ := tmpl.Prompts["system"].(string)
	if templatePromptRaw == "" {
		templatePromptRaw = "Analyze the provided passages and inputs."
	}

	// Replace params in template prompt
	templatePrompt := templatePromptRaw
	for k, v := range req.Inputs {
		templatePrompt = strings.ReplaceAll(templatePrompt, "{{"+k+"}}", v)
	}

	// Construct full prompt
	// Format:
	// [Global Prompt]
	//
	// Context:
	// Bible References: [Refs]
	//
	// Instructions:
	// [Template Prompt]
	//
	// User Inputs:
	// [Inputs]

	fullPrompt := fmt.Sprintf("%s\n\nBible References: %s\n\n%s\n\n", globalPrompt, referencesBlock, templatePrompt)

	// Add User Inputs for context if not fully covered by params
	fullPrompt += "User Inputs:\n"
	for k, v := range req.Inputs {
		fullPrompt += fmt.Sprintf("%s: %s\n", k, v)
	}

	// 5. Call AI
	payload := map[string]interface{}{
		"prompt": fullPrompt,
		"type":   "raw_template_generation",
	}

	resp, err := s.aiClient.ChatCompletion(ctx, payload)
	if err != nil {
		return "", err
	}

	aiOutput := ""
	if text, ok := resp["text"].(string); ok {
		aiOutput = text
	} else if choices, ok := resp["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if msg, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := msg["content"].(string); ok {
					aiOutput = content
				}
			}
		}
	}

	// 6. Format Final Output
	finalOutput := ""

	if tmpl.TemplateBody != "" {
		body := tmpl.TemplateBody
		body = strings.ReplaceAll(body, "{{passages}}", passagesBlock)
		body = strings.ReplaceAll(body, "{{generated}}", aiOutput)
		for k, v := range req.Inputs {
			body = strings.ReplaceAll(body, "{{"+k+"}}", v)
		}
		finalOutput = body
	} else {
		// Default Format
		if passagesBlock != "" {
			finalOutput = passagesBlock + "\n\n"
		}
		finalOutput += aiOutput
	}

	return finalOutput, nil
}
