package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"discipleship_journal_api/models"
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

func (m *MockReadingPlanService) GetPlanProgress(ctx context.Context, userID, planID uuid.UUID) ([]int, error) {
	args := m.Called(ctx, userID, planID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
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

func TestSubscribe_Handler(t *testing.T) {
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
	// Inject user ID for testing
	ctx := context.WithValue(req.Context(), TestUserKey, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}
