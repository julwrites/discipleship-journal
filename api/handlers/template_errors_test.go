package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTemplateHandler_ListMyTemplates_Unauthorized(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	req := httptest.NewRequest("GET", "/api/templates", nil)
	// No user context
	w := httptest.NewRecorder()

	handler.ListMyTemplates(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTemplateHandler_ListMyTemplates_ServiceError(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("GET", "/api/templates", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("ListTemplates", mock.Anything, userID).Return(nil, assert.AnError)

	handler.ListMyTemplates(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTemplateHandler_GetTemplate_InvalidUUID(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	req := httptest.NewRequest("GET", "/api/templates/invalid-uuid", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/api/templates/{id}", handler.GetTemplate)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTemplateHandler_UpdateTemplate_Unauthorized(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	tmplID := uuid.New()
	req := httptest.NewRequest("PUT", "/api/templates/"+tmplID.String(), strings.NewReader(`{}`))
	// No user context
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Put("/api/templates/{id}", handler.UpdateTemplate)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTemplateHandler_UpdateTemplate_InvalidUUID(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("PUT", "/api/templates/invalid-uuid", strings.NewReader(`{}`))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Put("/api/templates/{id}", handler.UpdateTemplate)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTemplateHandler_UpdateTemplate_InvalidBody(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())
	tmplID := uuid.New()

	req := httptest.NewRequest("PUT", "/api/templates/"+tmplID.String(), strings.NewReader(`{invalid}`))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Put("/api/templates/{id}", handler.UpdateTemplate)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTemplateHandler_DeleteTemplate_Unauthorized(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	tmplID := uuid.New()
	req := httptest.NewRequest("DELETE", "/api/templates/"+tmplID.String(), nil)
	// No user context
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/api/templates/{id}", handler.DeleteTemplate)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTemplateHandler_DeleteTemplate_InvalidUUID(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("DELETE", "/api/templates/invalid-uuid", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/api/templates/{id}", handler.DeleteTemplate)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTemplateHandler_CloneTemplate_Unauthorized(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	tmplID := uuid.New()
	req := httptest.NewRequest("POST", "/api/templates/"+tmplID.String()+"/clone", nil)
	// No user context
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/templates/{id}/clone", handler.CloneTemplate)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTemplateHandler_CloneTemplate_InvalidUUID(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("POST", "/api/templates/invalid-uuid/clone", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/templates/{id}/clone", handler.CloneTemplate)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTemplateHandler_Generate_InvalidUUID(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	req := httptest.NewRequest("POST", "/api/templates/invalid-uuid/generate", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/templates/{id}/generate", handler.Generate)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTemplateHandler_Generate_InvalidBody(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	tmplID := uuid.New()
	req := httptest.NewRequest("POST", "/api/templates/"+tmplID.String()+"/generate", strings.NewReader(`{invalid}`))
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/templates/{id}/generate", handler.Generate)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
