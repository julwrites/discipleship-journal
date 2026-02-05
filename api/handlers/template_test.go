package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"discipleship_journal_api/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTemplateService
type MockTemplateService struct {
	mock.Mock
}

func (m *MockTemplateService) CreateTemplate(ctx context.Context, tmpl *models.StudyTemplate) (*models.StudyTemplate, error) {
	args := m.Called(ctx, tmpl)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyTemplate), args.Error(1)
}

func (m *MockTemplateService) GetTemplate(ctx context.Context, id uuid.UUID) (*models.StudyTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyTemplate), args.Error(1)
}

func (m *MockTemplateService) ListTemplates(ctx context.Context, userID uuid.UUID) ([]*models.StudyTemplate, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.StudyTemplate), args.Error(1)
}

func (m *MockTemplateService) ListPublicTemplates(ctx context.Context) ([]*models.StudyTemplate, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.StudyTemplate), args.Error(1)
}

func (m *MockTemplateService) UpdateTemplate(ctx context.Context, tmpl *models.StudyTemplate) error {
	args := m.Called(ctx, tmpl)
	return args.Error(0)
}

func (m *MockTemplateService) DeleteTemplate(ctx context.Context, id, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockTemplateService) CloneTemplate(ctx context.Context, id, userID uuid.UUID) (*models.StudyTemplate, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyTemplate), args.Error(1)
}

func (m *MockTemplateService) GenerateContent(ctx context.Context, templateID uuid.UUID, req models.GenerateRequest) (string, error) {
	args := m.Called(ctx, templateID, req)
	return args.String(0), args.Error(1)
}

func TestTemplateHandler_CreateTemplate(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	payload := `{"title": "New Template"}`
	req := httptest.NewRequest("POST", "/api/templates", strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	expectedTmpl := &models.StudyTemplate{ID: uuid.New(), Title: "New Template", CreatorID: userID}
	mockService.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(tmpl *models.StudyTemplate) bool {
		return tmpl.Title == "New Template" && tmpl.CreatorID == userID
	})).Return(expectedTmpl, nil)

	handler.CreateTemplate(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.StudyTemplate
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "New Template", response.Title)
}

func TestTemplateHandler_GetTemplate(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	tmplID := uuid.New()
	req := httptest.NewRequest("GET", "/api/templates/"+tmplID.String(), nil)

	r := chi.NewRouter()
	r.Get("/api/templates/{id}", handler.GetTemplate)

	w := httptest.NewRecorder()

	expectedTmpl := &models.StudyTemplate{ID: tmplID, Title: "My Template"}
	mockService.On("GetTemplate", mock.Anything, tmplID).Return(expectedTmpl, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.StudyTemplate
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "My Template", response.Title)
}

func TestTemplateHandler_ListMyTemplates(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("GET", "/api/templates", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	list := []*models.StudyTemplate{{Title: "Template 1"}}
	mockService.On("ListTemplates", mock.Anything, userID).Return(list, nil)

	handler.ListMyTemplates(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string][]*models.StudyTemplate
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response["data"], 1)
}

func TestTemplateHandler_UpdateTemplate(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	userID := uuid.New()
	tmplID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Put("/api/templates/{id}", handler.UpdateTemplate)

	payload := `{"title": "Updated Template"}`
	req := httptest.NewRequest("PUT", "/api/templates/"+tmplID.String(), strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("UpdateTemplate", mock.Anything, mock.MatchedBy(func(tmpl *models.StudyTemplate) bool {
		return tmpl.ID == tmplID && tmpl.Title == "Updated Template" && tmpl.CreatorID == userID
	})).Return(nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTemplateHandler_DeleteTemplate(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	userID := uuid.New()
	tmplID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Delete("/api/templates/{id}", handler.DeleteTemplate)

	req := httptest.NewRequest("DELETE", "/api/templates/"+tmplID.String(), nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("DeleteTemplate", mock.Anything, tmplID, userID).Return(nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
