package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPlan_Handler_ExtendedErrors(t *testing.T) {
	t.Run("Invalid UUID", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)

		r := chi.NewRouter()
		r.Get("/api/reading-plans/{id}", handler.GetPlan)

		req := httptest.NewRequest("GET", "/api/reading-plans/invalid-uuid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetPlanDays Error", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		id := uuid.New()

		mockService.On("GetPlan", mock.Anything, id).Return(&models.ReadingPlan{
			ID: id, Title: "Test Plan",
		}, nil)

		mockService.On("GetPlanDays", mock.Anything, id).Return(([]*models.ReadingPlanDay)(nil), errors.New("db error"))

		r := chi.NewRouter()
		r.Get("/api/reading-plans/{id}", handler.GetPlan)
		req := httptest.NewRequest("GET", "/api/reading-plans/"+id.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSubscribe_Handler_ExtendedErrors(t *testing.T) {
	t.Run("Invalid UUID", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)

		r := chi.NewRouter()
		r.Post("/api/reading-plans/{id}/subscribe", handler.Subscribe)

		req := httptest.NewRequest("POST", "/api/reading-plans/invalid-uuid/subscribe", nil)
		userID := uuid.New()
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUnsubscribe_Handler_ExtendedErrors(t *testing.T) {
	t.Run("Invalid UUID", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)

		r := chi.NewRouter()
		r.Delete("/api/reading-plans/{id}/subscribe", handler.Unsubscribe)

		req := httptest.NewRequest("DELETE", "/api/reading-plans/invalid-uuid/subscribe", nil)
		userID := uuid.New()
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUnmarkDayComplete_Handler_ExtendedErrors(t *testing.T) {
	t.Run("Invalid Plan UUID", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)

		r := chi.NewRouter()
		r.Delete("/api/my-reading-plans/{id}/progress/{day_number}", handler.UnmarkDayComplete)

		req := httptest.NewRequest("DELETE", "/api/my-reading-plans/invalid/progress/1", nil)
		userID := uuid.New()
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Day Number", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		id := uuid.New()

		r := chi.NewRouter()
		r.Delete("/api/my-reading-plans/{id}/progress/{day_number}", handler.UnmarkDayComplete)

		req := httptest.NewRequest("DELETE", "/api/my-reading-plans/"+id.String()+"/progress/invalid", nil)
		userID := uuid.New()
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
