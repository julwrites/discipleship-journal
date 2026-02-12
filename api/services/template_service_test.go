package services

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"discipleship_journal_api/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateService_CreateTemplate(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	userID := uuid.New()
	tmpl := &models.StudyTemplate{
		CreatorID:       userID,
		Title:           "Test Template",
		Description:     "Desc",
		IsPublic:        false,
		BibleReferences: []string{"John 3:16"},
	}

	// Expect Exec. Using regexp to match the query partially
	mockDB.ExpectExec("INSERT INTO study_templates").
		WithArgs(
			pgxmock.AnyArg(), // ID
			userID,
			"Test Template",
			"Desc",
			pgxmock.AnyArg(), // structure
			pgxmock.AnyArg(), // prompts
			pgxmock.AnyArg(), // fields
			false,            // is_public
			pgxmock.AnyArg(), // bible_references
			false,            // allow_user_passages
			"",               // template_body
			"",               // required_version
			pgxmock.AnyArg(), // created_at
			pgxmock.AnyArg(), // updated_at
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	created, err := service.CreateTemplate(context.Background(), tmpl)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.NotEqual(t, uuid.Nil, created.ID)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

type GranularMockAIClient struct {
	*MockBibleAIClient
	PassageError error
	ChatError    error
}

func (m *GranularMockAIClient) GetPassage(ctx context.Context, reference string, version string) (map[string]interface{}, error) {
	if m.PassageError != nil {
		return nil, m.PassageError
	}
	return m.MockBibleAIClient.GetPassage(ctx, reference, version)
}

func (m *GranularMockAIClient) ChatCompletion(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	if m.ChatError != nil {
		return nil, m.ChatError
	}
	return m.MockBibleAIClient.ChatCompletion(ctx, payload)
}

func TestTemplateService_GenerateContent_PassageError(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := &GranularMockAIClient{
		MockBibleAIClient: NewMockBibleAIClient(),
		PassageError:      errors.New("passage fetch error"),
	}
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()
	userID := uuid.New()

	rows := pgxmock.NewRows([]string{
		"id", "creator_id", "title", "description", "structure", "prompts", "fields", "is_public",
		"bible_references", "allow_user_passages", "template_body", "required_version", "created_at", "updated_at",
	}).AddRow(
		tmplID, userID, "My Template", "Desc", []byte("{}"), []byte(`{"system":"System Prompt"}`), []byte("[]"), false,
		[]byte(`["John 3:16"]`), false, "", "ESV", time.Now(), time.Now(),
	)

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, creator_id, title, description")).
		WithArgs(tmplID).
		WillReturnRows(rows)

	req := models.GenerateRequest{}

	content, err := service.GenerateContent(context.Background(), tmplID, req)
	assert.NoError(t, err)
	// Should contain error message in the blockquote placeholder or skipped?
	// Implementation: passageTexts = append(passageTexts, fmt.Sprintf("%s: (Error fetching text)", ref))
	assert.Contains(t, content, "(Error fetching text)")
}

func TestTemplateService_GenerateContent_AIError(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := &GranularMockAIClient{
		MockBibleAIClient: NewMockBibleAIClient(),
		ChatError:         errors.New("ai error"),
	}
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()
	userID := uuid.New()

	rows := pgxmock.NewRows([]string{
		"id", "creator_id", "title", "description", "structure", "prompts", "fields", "is_public",
		"bible_references", "allow_user_passages", "template_body", "required_version", "created_at", "updated_at",
	}).AddRow(
		tmplID, userID, "My Template", "Desc", []byte("{}"), []byte("{}"), []byte("[]"), false,
		[]byte("[]"), false, "", "", time.Now(), time.Now(),
	)

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, creator_id, title, description")).
		WithArgs(tmplID).
		WillReturnRows(rows)

	req := models.GenerateRequest{}

	_, err = service.GenerateContent(context.Background(), tmplID, req)
	assert.Error(t, err)
	assert.EqualError(t, err, "ai error")
}

func TestTemplateService_GenerateContent_UserPassages(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()
	userID := uuid.New()

	// allow_user_passages = true
	rows := pgxmock.NewRows([]string{
		"id", "creator_id", "title", "description", "structure", "prompts", "fields", "is_public",
		"bible_references", "allow_user_passages", "template_body", "required_version", "created_at", "updated_at",
	}).AddRow(
		tmplID, userID, "User Passages Template", "Desc", []byte("{}"), []byte("{}"), []byte("[]"), false,
		[]byte(`["Gen 1:1"]`), true, "", "KJV", time.Now(), time.Now(),
	)

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, creator_id, title, description")).
		WithArgs(tmplID).
		WillReturnRows(rows)

	req := models.GenerateRequest{
		UserPassages: []string{"Psalm 23"},
		UserVersion:  "NIV", // Should override template default if logic allows?
		// Logic: if tmpl.RequiredVersion != "" { version = tmpl.RequiredVersion }
		// Here tmpl.RequiredVersion is "KJV", so it should use KJV.
	}

	content, err := service.GenerateContent(context.Background(), tmplID, req)
	assert.NoError(t, err)

	// MockAIClient returns version in the response
	// We can check if KJV was used in the formatted passage block
	// passageTexts = append(passageTexts, fmt.Sprintf("<blockquote><p><strong>%s (%s)</strong></p>%s</blockquote>", ref, version, text))
	assert.Contains(t, content, "Psalm 23 (KJV)")
}

func TestTemplateService_GetTemplate(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()
	userID := uuid.New()

	rows := pgxmock.NewRows([]string{
		"id", "creator_id", "title", "description", "structure", "prompts", "fields", "is_public",
		"bible_references", "allow_user_passages", "template_body", "required_version", "created_at", "updated_at",
	}).AddRow(
		tmplID, userID, "My Template", "Desc", []byte("{}"), []byte("{}"), []byte("[]"), false,
		[]byte("[]"), false, "", "", time.Now(), time.Now(),
	)

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, creator_id, title, description")).
		WithArgs(tmplID).
		WillReturnRows(rows)

	result, err := service.GetTemplate(context.Background(), tmplID)
	assert.NoError(t, err)
	assert.Equal(t, "My Template", result.Title)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTemplateService_ListTemplates(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	userID := uuid.New()

	rows := pgxmock.NewRows([]string{
		"id", "creator_id", "title", "description", "structure", "prompts", "fields", "is_public",
		"bible_references", "allow_user_passages", "template_body", "required_version", "created_at", "updated_at",
	}).AddRow(
		uuid.New(), userID, "T1", "D1", []byte("{}"), []byte("{}"), []byte("[]"), false,
		[]byte("[]"), false, "", "", time.Now(), time.Now(),
	)

	mockDB.ExpectQuery("SELECT .* FROM study_templates WHERE creator_id").
		WithArgs(userID).
		WillReturnRows(rows)

	list, err := service.ListTemplates(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTemplateService_UpdateTemplate(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()
	userID := uuid.New()
	tmpl := &models.StudyTemplate{
		ID:              tmplID,
		CreatorID:       userID,
		Title:           "Updated Title",
		Description:     "Updated Desc",
		BibleReferences: []string{"Gen 1:1"},
	}

	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE study_templates SET")).
		WithArgs(
			"Updated Title",
			"Updated Desc",
			pgxmock.AnyArg(), // structure
			pgxmock.AnyArg(), // prompts
			pgxmock.AnyArg(), // fields
			false,            // is_public
			pgxmock.AnyArg(), // bible_references
			false,            // allow_user_passages
			"",               // template_body
			"",               // required_version
			tmplID,
			userID,
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = service.UpdateTemplate(context.Background(), tmpl)
	assert.NoError(t, err)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTemplateService_DeleteTemplate(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()
	userID := uuid.New()

	mockDB.ExpectExec("DELETE FROM study_templates").
		WithArgs(tmplID, userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = service.DeleteTemplate(context.Background(), tmplID, userID)
	assert.NoError(t, err)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTemplateService_CloneTemplate(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	originalID := uuid.New()
	userID := uuid.New()

	// 1. Expect GetTemplate
	rows := pgxmock.NewRows([]string{
		"id", "creator_id", "title", "description", "structure", "prompts", "fields", "is_public",
		"bible_references", "allow_user_passages", "template_body", "required_version", "created_at", "updated_at",
	}).AddRow(
		originalID, userID, "Original", "Desc", []byte("{}"), []byte("{}"), []byte("[]"), false,
		[]byte("[]"), false, "", "", time.Now(), time.Now(),
	)

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, creator_id, title, description")).
		WithArgs(originalID).
		WillReturnRows(rows)

	// 2. Expect CreateTemplate (Insert)
	mockDB.ExpectExec("INSERT INTO study_templates").
		WithArgs(
			pgxmock.AnyArg(), // New ID
			userID,
			"Original (Copy)", // Title
			"Desc",
			pgxmock.AnyArg(), // structure
			pgxmock.AnyArg(), // prompts
			pgxmock.AnyArg(), // fields
			false,            // is_public
			pgxmock.AnyArg(), // bible_references
			false,            // allow_user_passages
			"",               // template_body
			"",               // required_version
			pgxmock.AnyArg(), // created_at
			pgxmock.AnyArg(), // updated_at
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	cloned, err := service.CloneTemplate(context.Background(), originalID, userID)
	assert.NoError(t, err)
	assert.Equal(t, "Original (Copy)", cloned.Title)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTemplateService_GenerateContent(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()
	userID := uuid.New()

	// 1. Expect GetTemplate
	rows := pgxmock.NewRows([]string{
		"id", "creator_id", "title", "description", "structure", "prompts", "fields", "is_public",
		"bible_references", "allow_user_passages", "template_body", "required_version", "created_at", "updated_at",
	}).AddRow(
		tmplID, userID, "My Template", "Desc", []byte("{}"), []byte(`{"system":"System Prompt"}`), []byte("[]"), false,
		[]byte(`["John 3:16"]`), false, "{{passages}}\n{{generated}}", "ESV", time.Now(), time.Now(),
	)

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, creator_id, title, description")).
		WithArgs(tmplID).
		WillReturnRows(rows)

	req := models.GenerateRequest{
		Inputs: map[string]string{"key": "value"},
	}

	content, err := service.GenerateContent(context.Background(), tmplID, req)
	assert.NoError(t, err)
	assert.Contains(t, content, "For God so loved the world") // From MockBibleAIClient
	assert.Contains(t, content, "This is a mocked AI response") // From MockBibleAIClient

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTemplateService_ListPublicTemplates(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	rows := pgxmock.NewRows([]string{
		"id", "creator_id", "title", "description", "structure", "prompts", "fields", "is_public",
		"bible_references", "allow_user_passages", "template_body", "required_version", "created_at", "updated_at",
	}).AddRow(
		uuid.New(), uuid.New(), "Public Template", "Desc", []byte("{}"), []byte("{}"), []byte("[]"), true,
		[]byte("[]"), false, "", "", time.Now(), time.Now(),
	)

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, creator_id, title, description")).
		WithArgs().
		WillReturnRows(rows)

	list, err := service.ListPublicTemplates(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "Public Template", list[0].Title)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTemplateService_CreateTemplate_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	userID := uuid.New()
	tmpl := &models.StudyTemplate{
		CreatorID:       userID,
		Title:           "Test Template",
		Description:     "Desc",
		IsPublic:        false,
		BibleReferences: []string{"John 3:16"},
	}

	mockDB.ExpectExec("INSERT INTO study_templates").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnError(errors.New("db error"))

	created, err := service.CreateTemplate(context.Background(), tmpl)
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.EqualError(t, err, "db error")

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTemplateService_UpdateTemplate_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()
	userID := uuid.New()
	tmpl := &models.StudyTemplate{
		ID:              tmplID,
		CreatorID:       userID,
		Title:           "Updated Title",
		Description:     "Updated Desc",
		BibleReferences: []string{"Gen 1:1"},
	}

	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE study_templates SET")).
		WithArgs(
			"Updated Title", "Updated Desc", pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), false,
			pgxmock.AnyArg(), false, "", "", tmplID, userID,
		).
		WillReturnError(errors.New("db error"))

	err = service.UpdateTemplate(context.Background(), tmpl)
	assert.Error(t, err)
	assert.EqualError(t, err, "db error")

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTemplateService_DeleteTemplate_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()
	userID := uuid.New()

	mockDB.ExpectExec("DELETE FROM study_templates").
		WithArgs(tmplID, userID).
		WillReturnError(errors.New("db error"))

	err = service.DeleteTemplate(context.Background(), tmplID, userID)
	assert.Error(t, err)
	assert.EqualError(t, err, "db error")

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestTemplateService_GetTemplate_NotFound(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	mockAI := NewMockBibleAIClient()
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, creator_id, title, description")).
		WithArgs(tmplID).
		WillReturnError(pgx.ErrNoRows)

	result, err := service.GetTemplate(context.Background(), tmplID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, models.ErrNotFound, err)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}
