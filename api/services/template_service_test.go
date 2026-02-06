package services

import (
	"context"
	"regexp"
	"testing"
	"time"

	"discipleship_journal_api/models"
	"github.com/google/uuid"
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
