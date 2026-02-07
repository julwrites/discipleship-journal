package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTemplateHandler_Coverage(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)
	userID := uuid.New()

	t.Run("CreateTemplate_Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/templates", nil)
		w := httptest.NewRecorder()
		// No TestUserKey
		handler.CreateTemplate(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("CreateTemplate_ServiceError", func(t *testing.T) {
		mockService.On("CreateTemplate", mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

		body := `{"title": "Test"}`
		req := httptest.NewRequest("POST", "/api/templates", bytes.NewBufferString(body))
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateTemplate(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("ListPublicTemplates_ServiceError", func(t *testing.T) {
		mockService.On("ListPublicTemplates", mock.Anything).Return(nil, errors.New("service error"))

		req := httptest.NewRequest("GET", "/api/templates/public", nil)
		w := httptest.NewRecorder()

		handler.ListPublicTemplates(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Generate_ServiceError", func(t *testing.T) {
		id := uuid.New()
		mockService.On("GenerateContent", mock.Anything, id, mock.Anything).Return("", errors.New("service error"))

		body := `{"passages": []}`
		req := httptest.NewRequest("POST", "/api/templates/"+id.String()+"/generate", bytes.NewBufferString(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", id.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.Generate(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetTemplate_NotFound", func(t *testing.T) {
		id := uuid.New()
		mockService.On("GetTemplate", mock.Anything, id).Return(nil, errors.New("not found"))

		req := httptest.NewRequest("GET", "/api/templates/"+id.String(), nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", id.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.GetTemplate(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
