package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMemoryVerseHandler_GetPackDetails(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	packID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Get("/api/verse-packs/{id}", handler.GetPackDetails)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/verse-packs/"+packID.String(), nil)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		pack := &models.VersePack{ID: packID, Title: "My Pack", UserID: &userID}
		verses := []*models.MemoryVerse{{ID: uuid.New(), Reference: "John 3:16"}}

		mockService.On("GetPack", mock.Anything, packID, userID).Return(pack, nil).Once()
		mockService.On("GetVerses", mock.Anything, packID, userID).Return(verses, nil).Once()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Check pack details
		packMap, ok := resp["pack"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "My Pack", packMap["title"])

		// Check verses
		versesList, ok := resp["verses"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, versesList, 1)
	})

	t.Run("pack_not_found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/verse-packs/"+packID.String(), nil)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockService.On("GetPack", mock.Anything, packID, userID).Return(nil, models.ErrNotFound).Once()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("verses_error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/verse-packs/"+packID.String(), nil)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		pack := &models.VersePack{ID: packID, Title: "My Pack"}

		mockService.On("GetPack", mock.Anything, packID, userID).Return(pack, nil).Once()
		mockService.On("GetVerses", mock.Anything, packID, userID).Return(nil, assert.AnError).Once()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestMemoryVerseHandler_DeletePack(t *testing.T) {
	mockService := new(MockMemoryVerseService)
	handler := NewMemoryVerseHandler(mockService)

	userID := uuid.New()
	packID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	r := chi.NewRouter()
	r.Delete("/api/verse-packs/{id}", handler.DeletePack)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/verse-packs/"+packID.String(), nil)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockService.On("DeletePack", mock.Anything, packID, userID).Return(nil).Once()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"success": true}`, w.Body.String())
	})

	t.Run("not_found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/verse-packs/"+packID.String(), nil)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockService.On("DeletePack", mock.Anything, packID, userID).Return(models.ErrNotFound).Once()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
