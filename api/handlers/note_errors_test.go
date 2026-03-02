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

func TestNoteHandler_CreateNote_Errors(t *testing.T) {
	mockService := new(MockNoteService)
	handler := NewNoteHandler(nil, mockService)

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/notes", strings.NewReader(`{"title":"Test"}`))
		w := httptest.NewRecorder()
		handler.CreateNote(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/notes", strings.NewReader(`{invalid}`))
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.CreateNote(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/notes", strings.NewReader(`{"title":"Test"}`))
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockService.On("CreateNote", mock.Anything, userID.String(), "Test", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("service error"))

		handler.CreateNote(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNoteHandler_GetTags_Errors(t *testing.T) {
	mockService := new(MockNoteService)
	handler := NewNoteHandler(nil, mockService)

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/notes/tags", nil)
		w := httptest.NewRecorder()
		handler.GetTags(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/notes/tags", nil)
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockService.On("GetUserTags", mock.Anything, userID.String()).
			Return(nil, errors.New("service error"))

		handler.GetTags(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNoteHandler_UpdateNote_Errors(t *testing.T) {
	mockService := new(MockNoteService)
	handler := NewNoteHandler(nil, mockService)

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/notes/id", strings.NewReader(`{}`))
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Put("/api/notes/{id}", handler.UpdateNote)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Note: UpdateNote doesn't check if ID is valid UUID if it's treated as string by service
		// But let's check if the handler parses it or not.
		// Looking at code: noteID := chi.URLParam(r, "id")
		// It uses string ID.
		// So invalid ID format might not trigger bad request unless service complains.
		// But empty ID might be an issue? chi handles routing so empty ID usually 404.
	})

	t.Run("Invalid Body", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/notes/id", strings.NewReader(`{invalid}`))
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Put("/api/notes/{id}", handler.UpdateNote)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		noteID := "note-id"
		req := httptest.NewRequest("PUT", "/api/notes/"+noteID, strings.NewReader(`{"title":"New"}`))
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockService.On("UpdateNote", mock.Anything, userID.String(), noteID, "New", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("service error"))

		r := chi.NewRouter()
		r.Put("/api/notes/{id}", handler.UpdateNote)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Note Not Found", func(t *testing.T) {
		noteID := "note-id"
		req := httptest.NewRequest("PUT", "/api/notes/"+noteID, strings.NewReader(`{"title":"New"}`))
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockService.On("UpdateNote", mock.Anything, userID.String(), noteID, "New", mock.Anything, mock.Anything, mock.Anything).
			Return(models.ErrNotFound)

		r := chi.NewRouter()
		r.Put("/api/notes/{id}", handler.UpdateNote)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestNoteHandler_DeleteNote_Errors(t *testing.T) {
	mockService := new(MockNoteService)
	handler := NewNoteHandler(nil, mockService)

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/notes/id", nil)
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Delete("/api/notes/{id}", handler.DeleteNote)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		noteID := "note-id"
		req := httptest.NewRequest("DELETE", "/api/notes/"+noteID, nil)
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockService.On("DeleteNote", mock.Anything, userID.String(), noteID).
			Return(errors.New("service error"))

		r := chi.NewRouter()
		r.Delete("/api/notes/{id}", handler.DeleteNote)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNoteHandler_GetNote_Errors(t *testing.T) {
	mockService := new(MockNoteService)
	handler := NewNoteHandler(nil, mockService)

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/notes/id", nil)
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Get("/api/notes/{id}", handler.GetNote)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		noteID := "note-id"
		req := httptest.NewRequest("GET", "/api/notes/"+noteID, nil)
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockService.On("GetNote", mock.Anything, userID.String(), noteID).
			Return(nil, errors.New("service error"))

		r := chi.NewRouter()
		r.Get("/api/notes/{id}", handler.GetNote)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		noteID := "note-id"
		req := httptest.NewRequest("GET", "/api/notes/"+noteID, nil)
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, userID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockService.On("GetNote", mock.Anything, userID.String(), noteID).
			Return(nil, models.ErrNotFound)

		r := chi.NewRouter()
		r.Get("/api/notes/{id}", handler.GetNote)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
