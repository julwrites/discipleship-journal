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

// MockMemoryVerseService
type MockMemoryVerseService struct {
	mock.Mock
}

func (m *MockMemoryVerseService) GetPacks(ctx context.Context, userID uuid.UUID, typeFilter string) ([]*models.VersePack, error) {
	args := m.Called(ctx, userID, typeFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.VersePack), args.Error(1)
}

func (m *MockMemoryVerseService) GetPack(ctx context.Context, packID uuid.UUID, userID uuid.UUID) (*models.VersePack, error) {
	args := m.Called(ctx, packID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VersePack), args.Error(1)
}

func (m *MockMemoryVerseService) CreatePack(ctx context.Context, pack *models.VersePack) (*models.VersePack, error) {
	args := m.Called(ctx, pack)
	return args.Get(0).(*models.VersePack), args.Error(1)
}

func (m *MockMemoryVerseService) GetVerses(ctx context.Context, packID uuid.UUID, userID uuid.UUID) ([]*models.MemoryVerse, error) {
	args := m.Called(ctx, packID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MemoryVerse), args.Error(1)
}

func (m *MockMemoryVerseService) CreateVerse(ctx context.Context, verse *models.MemoryVerse) (*models.MemoryVerse, error) {
	args := m.Called(ctx, verse)
	return args.Get(0).(*models.MemoryVerse), args.Error(1)
}

func (m *MockMemoryVerseService) ClonePack(ctx context.Context, packID uuid.UUID, userID uuid.UUID, newTitle string) (*models.VersePack, error) {
	args := m.Called(ctx, packID, userID, newTitle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VersePack), args.Error(1)
}

func (m *MockMemoryVerseService) DeletePack(ctx context.Context, packID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, packID, userID)
	return args.Error(0)
}

func (m *MockMemoryVerseService) SearchVerses(ctx context.Context, userID uuid.UUID, query string) ([]*models.MemoryVerse, error) {
	args := m.Called(ctx, userID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.MemoryVerse), args.Error(1)
}

func (m *MockMemoryVerseService) UpdateVerse(ctx context.Context, verse *models.MemoryVerse, userID uuid.UUID) error {
	args := m.Called(ctx, verse, userID)
	return args.Error(0)
}

func (m *MockMemoryVerseService) DeleteVerse(ctx context.Context, verseID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, verseID, userID)
	return args.Error(0)
}

func (m *MockMemoryVerseService) SetVersePreference(ctx context.Context, userID, verseID uuid.UUID, version string) error {
	args := m.Called(ctx, userID, verseID, version)
	return args.Error(0)
}

func (m *MockMemoryVerseService) RemoveVersePreference(ctx context.Context, userID, verseID uuid.UUID) error {
	args := m.Called(ctx, userID, verseID)
	return args.Error(0)
}

func TestMemoryVerseHandler_GetPacks(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("GET", "/api/verse-packs?type=user", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	packs := []*models.VersePack{{Title: "My Pack"}}
	mockService.On("GetPacks", mock.Anything, userID, "user").Return(packs, nil)

	handler.GetPacks(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string][]*models.VersePack
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response["data"], 1)
}

func TestMemoryVerseHandler_CreatePack(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	payload := `{"title": "New Pack"}`
	req := httptest.NewRequest("POST", "/api/verse-packs", strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	expectedPack := &models.VersePack{ID: uuid.New(), Title: "New Pack"}
	mockService.On("CreatePack", mock.Anything, mock.MatchedBy(func(p *models.VersePack) bool {
		return p.Title == "New Pack" && *p.UserID == userID
	})).Return(expectedPack, nil)

	handler.CreatePack(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestMemoryVerseHandler_CreateVerseInPack(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	packID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	// Need to setup URL params using Chi Context or simulate it?
	// With httptest, we usually mock the router logic or manually inject params if the handler reads them directly.
	// The handler uses chi.URLParam. We can wrap the handler in a chi router to test properly or inject context.
	r := chi.NewRouter()
	r.Post("/api/verse-packs/{id}/verses", handler.CreateVerseInPack)

	payload := `{"reference": "John 3:16"}`
	req := httptest.NewRequest("POST", "/api/verse-packs/"+packID.String()+"/verses", strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// 1. GetPack (check ownership)
	mockService.On("GetPack", mock.Anything, packID, userID).Return(&models.VersePack{ID: packID, UserID: &userID}, nil)

	// 2. CreateVerse
	mockService.On("CreateVerse", mock.Anything, mock.MatchedBy(func(v *models.MemoryVerse) bool {
		return v.VersePackID == packID && v.Reference == "John 3:16"
	})).Return(&models.MemoryVerse{ID: uuid.New()}, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestMemoryVerseHandler_SearchVerses(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("GET", "/api/memory-verses?q=love", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	verses := []*models.MemoryVerse{{Reference: "John 3:16"}}
	mockService.On("SearchVerses", mock.Anything, userID, "love").Return(verses, nil)

	handler.SearchVerses(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string][]*models.MemoryVerse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response["data"], 1)
}

func TestMemoryVerseHandler_UpdateVerse(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Put("/api/memory-verses/{verseId}", handler.UpdateVerse)

	payload := `{"reference": "John 3:17", "version": "NIV"}`
	req := httptest.NewRequest("PUT", "/api/memory-verses/"+verseID.String(), strings.NewReader(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("UpdateVerse", mock.Anything, mock.MatchedBy(func(v *models.MemoryVerse) bool {
		return v.ID == verseID && v.Reference == "John 3:17" && v.Version == "NIV"
	}), userID).Return(nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMemoryVerseHandler_DeleteVerse(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	verseID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Delete("/api/memory-verses/{verseId}", handler.DeleteVerse)

	req := httptest.NewRequest("DELETE", "/api/memory-verses/"+verseID.String(), nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("DeleteVerse", mock.Anything, verseID, userID).Return(nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
