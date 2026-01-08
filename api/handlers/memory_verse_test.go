package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"discipleship_journal_api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMemoryVerseService
type MockMemoryVerseService struct {
	mock.Mock
}

func (m *MockMemoryVerseService) SearchVerses(ctx context.Context, userID uuid.UUID, query string, tags []string) ([]*models.MemoryVerse, error) {
	args := m.Called(ctx, userID, query, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MemoryVerse), args.Error(1)
}

func (m *MockMemoryVerseService) CreateVerse(ctx context.Context, verse *models.MemoryVerse) (*models.MemoryVerse, error) {
	args := m.Called(ctx, verse)
	return args.Get(0).(*models.MemoryVerse), args.Error(1)
}

func (m *MockMemoryVerseService) GetSystemPacks(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func TestMemoryVerseHandler_SearchVerses(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	// Mock Auth Context
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	// Setup Request
	req := httptest.NewRequest("GET", "/api/memory-verses?q=love", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Expectations
	verses := []*models.MemoryVerse{{Reference: "John 3:16"}}
	mockService.On("SearchVerses", mock.Anything, userID, "love", []string(nil)).Return(verses, nil)

	handler.SearchVerses(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]*models.MemoryVerse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Len(t, response["data"], 1)
}

func TestMemoryVerseHandler_CreateVerse(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	payload := `{"reference": "Rom 8:28", "text": "And we know...", "tags": ["comfort"]}`
	req := httptest.NewRequest("POST", "/api/memory-verses", strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	expectedVerse := &models.MemoryVerse{ID: uuid.New(), Reference: "Rom 8:28"}

	// We use mock.MatchedBy to validate the argument passed to service
	mockService.On("CreateVerse", mock.Anything, mock.MatchedBy(func(v *models.MemoryVerse) bool {
		return v.Reference == "Rom 8:28" && *v.UserID == userID
	})).Return(expectedVerse, nil)

	handler.CreateVerse(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
