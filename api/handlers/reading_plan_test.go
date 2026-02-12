package handlers

import (
	"context"
	"errors"
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

// MockReadingPlanService is a mock implementation of ReadingPlanService
type MockReadingPlanService struct {
	mock.Mock
}

func (m *MockReadingPlanService) GetAllPlans(ctx context.Context) ([]*models.ReadingPlan, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ReadingPlan), args.Error(1)
}

func (m *MockReadingPlanService) GetPlan(ctx context.Context, id uuid.UUID) (*models.ReadingPlan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ReadingPlan), args.Error(1)
}

func (m *MockReadingPlanService) GetPlanDays(ctx context.Context, planID uuid.UUID) ([]*models.ReadingPlanDay, error) {
	args := m.Called(ctx, planID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ReadingPlanDay), args.Error(1)
}

func (m *MockReadingPlanService) Subscribe(ctx context.Context, userID, planID uuid.UUID) (*models.UserReadingPlan, error) {
	args := m.Called(ctx, userID, planID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserReadingPlan), args.Error(1)
}

func (m *MockReadingPlanService) GetUserPlans(ctx context.Context, userID uuid.UUID) ([]*models.UserReadingPlan, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UserReadingPlan), args.Error(1)
}

func (m *MockReadingPlanService) MarkDayComplete(ctx context.Context, userID, planID uuid.UUID, dayNumber int) error {
	args := m.Called(ctx, userID, planID, dayNumber)
	return args.Error(0)
}

func (m *MockReadingPlanService) UnmarkDayComplete(ctx context.Context, userID, planID uuid.UUID, dayNumber int) error {
	args := m.Called(ctx, userID, planID, dayNumber)
	return args.Error(0)
}

func (m *MockReadingPlanService) GetPlanProgress(ctx context.Context, userID, planID uuid.UUID) ([]int, error) {
	args := m.Called(ctx, userID, planID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockReadingPlanService) Unsubscribe(ctx context.Context, userID, planID uuid.UUID) error {
	args := m.Called(ctx, userID, planID)
	return args.Error(0)
}

func TestGetAllPlans_Handler(t *testing.T) {
	mockService := new(MockReadingPlanService)
	handler := NewReadingPlanHandler(mockService)

	mockService.On("GetAllPlans", mock.Anything).Return([]*models.ReadingPlan{
		{Title: "Test Plan"},
	}, nil)

	req := httptest.NewRequest("GET", "/api/reading-plans", nil)
	w := httptest.NewRecorder()

	handler.GetAllPlans(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllPlans_Handler_Error(t *testing.T) {
	mockService := new(MockReadingPlanService)
	handler := NewReadingPlanHandler(mockService)

	mockService.On("GetAllPlans", mock.Anything).Return(([]*models.ReadingPlan)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/reading-plans", nil)
	w := httptest.NewRecorder()

	handler.GetAllPlans(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetPlan_Handler(t *testing.T) {
	mockService := new(MockReadingPlanService)
	handler := NewReadingPlanHandler(mockService)

	id := uuid.New()

	mockService.On("GetPlan", mock.Anything, id).Return(&models.ReadingPlan{
		ID: id, Title: "Test Plan",
	}, nil)

	mockService.On("GetPlanDays", mock.Anything, id).Return([]*models.ReadingPlanDay{}, nil)

	r := chi.NewRouter()
	r.Get("/api/reading-plans/{id}", handler.GetPlan)

	req := httptest.NewRequest("GET", "/api/reading-plans/"+id.String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetPlan_Handler_Errors(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		id := uuid.New()

		mockService.On("GetPlan", mock.Anything, id).Return((*models.ReadingPlan)(nil), models.ErrNotFound)

		r := chi.NewRouter()
		r.Get("/api/reading-plans/{id}", handler.GetPlan)
		req := httptest.NewRequest("GET", "/api/reading-plans/"+id.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		id := uuid.New()

		mockService.On("GetPlan", mock.Anything, id).Return((*models.ReadingPlan)(nil), errors.New("db error"))

		r := chi.NewRouter()
		r.Get("/api/reading-plans/{id}", handler.GetPlan)
		req := httptest.NewRequest("GET", "/api/reading-plans/"+id.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestReadingPlanHandler_Unauthorized(t *testing.T) {
	mockService := new(MockReadingPlanService)
	handler := NewReadingPlanHandler(mockService)

	tests := []struct {
		name        string
		method      string
		path        string
		handlerFunc http.HandlerFunc
	}{
		{"Subscribe", "POST", "/api/reading-plans/uuid/subscribe", handler.Subscribe},
		{"Unsubscribe", "DELETE", "/api/reading-plans/uuid/subscribe", handler.Unsubscribe},
		{"GetUserPlans", "GET", "/api/my-reading-plans", handler.GetUserPlans},
		{"MarkDayComplete", "POST", "/api/my-reading-plans/uuid/progress", handler.MarkDayComplete},
		{"UnmarkDayComplete", "DELETE", "/api/my-reading-plans/uuid/progress/1", handler.UnmarkDayComplete},
		{"GetPlanProgress", "GET", "/api/my-reading-plans/uuid/progress", handler.GetPlanProgress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			// No context with UserID

			w := httptest.NewRecorder()
			tt.handlerFunc(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestSubscribe_Handler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)

		planID := uuid.New()
		userID := uuid.New()

		mockService.On("Subscribe", mock.Anything, userID, planID).Return(&models.UserReadingPlan{
			UserID: userID, ReadingPlanID: planID, Status: "active",
		}, nil)

		r := chi.NewRouter()
		r.Post("/api/reading-plans/{id}/subscribe", handler.Subscribe)

		req := httptest.NewRequest("POST", "/api/reading-plans/"+planID.String()+"/subscribe", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Already Exists", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)

		planID := uuid.New()
		userID := uuid.New()

		mockService.On("Subscribe", mock.Anything, userID, planID).Return(nil, models.ErrAlreadyExists)

		r := chi.NewRouter()
		r.Post("/api/reading-plans/{id}/subscribe", handler.Subscribe)

		req := httptest.NewRequest("POST", "/api/reading-plans/"+planID.String()+"/subscribe", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)

		planID := uuid.New()
		userID := uuid.New()

		mockService.On("Subscribe", mock.Anything, userID, planID).Return(nil, errors.New("db error"))

		r := chi.NewRouter()
		r.Post("/api/reading-plans/{id}/subscribe", handler.Subscribe)

		req := httptest.NewRequest("POST", "/api/reading-plans/"+planID.String()+"/subscribe", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUnsubscribe_Handler(t *testing.T) {
	mockService := new(MockReadingPlanService)
	handler := NewReadingPlanHandler(mockService)

	planID := uuid.New()
	userID := uuid.New()

	mockService.On("Unsubscribe", mock.Anything, userID, planID).Return(nil)

	r := chi.NewRouter()
	r.Delete("/api/reading-plans/{id}/subscribe", handler.Unsubscribe)

	req := httptest.NewRequest("DELETE", "/api/reading-plans/"+planID.String()+"/subscribe", nil)
	ctx := context.WithValue(req.Context(), TestUserKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestUnsubscribe_Handler_Errors(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		planID := uuid.New()
		userID := uuid.New()

		mockService.On("Unsubscribe", mock.Anything, userID, planID).Return(models.ErrNotFound)

		r := chi.NewRouter()
		r.Delete("/api/reading-plans/{id}/subscribe", handler.Unsubscribe)
		req := httptest.NewRequest("DELETE", "/api/reading-plans/"+planID.String()+"/subscribe", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req.WithContext(ctx))

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		planID := uuid.New()
		userID := uuid.New()

		mockService.On("Unsubscribe", mock.Anything, userID, planID).Return(errors.New("db error"))

		r := chi.NewRouter()
		r.Delete("/api/reading-plans/{id}/subscribe", handler.Unsubscribe)
		req := httptest.NewRequest("DELETE", "/api/reading-plans/"+planID.String()+"/subscribe", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req.WithContext(ctx))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetUserPlans_Handler(t *testing.T) {
	mockService := new(MockReadingPlanService)
	handler := NewReadingPlanHandler(mockService)

	userID := uuid.New()

	mockService.On("GetUserPlans", mock.Anything, userID).Return([]*models.UserReadingPlan{
		{UserID: userID, Status: "active"},
	}, nil)

	r := chi.NewRouter()
	r.Get("/api/my-reading-plans", handler.GetUserPlans)

	req := httptest.NewRequest("GET", "/api/my-reading-plans", nil)
	ctx := context.WithValue(req.Context(), TestUserKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetUserPlans_Handler_Error(t *testing.T) {
	mockService := new(MockReadingPlanService)
	handler := NewReadingPlanHandler(mockService)
	userID := uuid.New()

	mockService.On("GetUserPlans", mock.Anything, userID).Return(([]*models.UserReadingPlan)(nil), errors.New("db error"))

	r := chi.NewRouter()
	r.Get("/api/my-reading-plans", handler.GetUserPlans)
	req := httptest.NewRequest("GET", "/api/my-reading-plans", nil)
	ctx := context.WithValue(req.Context(), TestUserKey, userID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req.WithContext(ctx))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMarkDayComplete_Handler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)

		planID := uuid.New()
		userID := uuid.New()
		dayNumber := 1

		mockService.On("MarkDayComplete", mock.Anything, userID, planID, dayNumber).Return(nil)

		r := chi.NewRouter()
		r.Post("/api/my-reading-plans/{id}/progress", handler.MarkDayComplete)

		reqBody := `{"day_number": 1}`
		req := httptest.NewRequest("POST", "/api/my-reading-plans/"+planID.String()+"/progress", strings.NewReader(reqBody))
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Body", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)

		planID := uuid.New()
		userID := uuid.New()

		r := chi.NewRouter()
		r.Post("/api/my-reading-plans/{id}/progress", handler.MarkDayComplete)

		reqBody := `{"day_number": "invalid"}`
		req := httptest.NewRequest("POST", "/api/my-reading-plans/"+planID.String()+"/progress", strings.NewReader(reqBody))
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		planID := uuid.New()
		userID := uuid.New()

		mockService.On("MarkDayComplete", mock.Anything, userID, planID, 1).Return(models.ErrNotFound)

		r := chi.NewRouter()
		r.Post("/api/my-reading-plans/{id}/progress", handler.MarkDayComplete)
		reqBody := `{"day_number": 1}`
		req := httptest.NewRequest("POST", "/api/my-reading-plans/"+planID.String()+"/progress", strings.NewReader(reqBody))
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req.WithContext(ctx))

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		planID := uuid.New()
		userID := uuid.New()

		mockService.On("MarkDayComplete", mock.Anything, userID, planID, 1).Return(errors.New("db error"))

		r := chi.NewRouter()
		r.Post("/api/my-reading-plans/{id}/progress", handler.MarkDayComplete)
		reqBody := `{"day_number": 1}`
		req := httptest.NewRequest("POST", "/api/my-reading-plans/"+planID.String()+"/progress", strings.NewReader(reqBody))
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req.WithContext(ctx))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUnmarkDayComplete_Handler(t *testing.T) {
	mockService := new(MockReadingPlanService)
	handler := NewReadingPlanHandler(mockService)

	planID := uuid.New()
	userID := uuid.New()
	dayNumber := 1

	mockService.On("UnmarkDayComplete", mock.Anything, userID, planID, dayNumber).Return(nil)

	r := chi.NewRouter()
	r.Delete("/api/my-reading-plans/{id}/progress/{day_number}", handler.UnmarkDayComplete)

	req := httptest.NewRequest("DELETE", "/api/my-reading-plans/"+planID.String()+"/progress/1", nil)
	ctx := context.WithValue(req.Context(), TestUserKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUnmarkDayComplete_Handler_Errors(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		planID := uuid.New()
		userID := uuid.New()

		mockService.On("UnmarkDayComplete", mock.Anything, userID, planID, 1).Return(models.ErrNotFound)

		r := chi.NewRouter()
		r.Delete("/api/my-reading-plans/{id}/progress/{day_number}", handler.UnmarkDayComplete)
		req := httptest.NewRequest("DELETE", "/api/my-reading-plans/"+planID.String()+"/progress/1", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req.WithContext(ctx))

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		planID := uuid.New()
		userID := uuid.New()

		mockService.On("UnmarkDayComplete", mock.Anything, userID, planID, 1).Return(errors.New("db error"))

		r := chi.NewRouter()
		r.Delete("/api/my-reading-plans/{id}/progress/{day_number}", handler.UnmarkDayComplete)
		req := httptest.NewRequest("DELETE", "/api/my-reading-plans/"+planID.String()+"/progress/1", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req.WithContext(ctx))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetPlanProgress_Handler(t *testing.T) {
	mockService := new(MockReadingPlanService)
	handler := NewReadingPlanHandler(mockService)

	planID := uuid.New()
	userID := uuid.New()

	mockService.On("GetPlanProgress", mock.Anything, userID, planID).Return([]int{1, 2, 3}, nil)

	r := chi.NewRouter()
	r.Get("/api/my-reading-plans/{id}/progress", handler.GetPlanProgress)

	req := httptest.NewRequest("GET", "/api/my-reading-plans/"+planID.String()+"/progress", nil)
	ctx := context.WithValue(req.Context(), TestUserKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetPlanProgress_Handler_Errors(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		planID := uuid.New()
		userID := uuid.New()

		mockService.On("GetPlanProgress", mock.Anything, userID, planID).Return(([]int)(nil), models.ErrNotFound)

		r := chi.NewRouter()
		r.Get("/api/my-reading-plans/{id}/progress", handler.GetPlanProgress)
		req := httptest.NewRequest("GET", "/api/my-reading-plans/"+planID.String()+"/progress", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req.WithContext(ctx))

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockReadingPlanService)
		handler := NewReadingPlanHandler(mockService)
		planID := uuid.New()
		userID := uuid.New()

		mockService.On("GetPlanProgress", mock.Anything, userID, planID).Return(([]int)(nil), errors.New("db error"))

		r := chi.NewRouter()
		r.Get("/api/my-reading-plans/{id}/progress", handler.GetPlanProgress)
		req := httptest.NewRequest("GET", "/api/my-reading-plans/"+planID.String()+"/progress", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req.WithContext(ctx))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
