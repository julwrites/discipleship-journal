package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/models"
	"discipleship_journal_api/services"

	"firebase.google.com/go/v4/auth"
	chi "github.com/go-chi/chi/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (pgxmock.PgxPoolIface, *MockNoteService, *NoteHandler) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	noteServiceMock := new(MockNoteService)
	handler := NewNoteHandler(dbMock, noteServiceMock)
	return dbMock, noteServiceMock, handler
}

func TestGetNotes(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		dbMock, noteServiceMock, handler := setupTest(t)
		defer dbMock.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		content := json.RawMessage(`{"text": "content"}`)
		serviceNotes := []services.Note{
			{
				ID:        "note-1",
				UserID:    userUUID,
				Title:     "Title 1",
				Content:   content,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		// Updated to match filter struct
		noteServiceMock.On("GetNotes", mock.Anything, userUUID, 1, 20, services.NoteFilter{}).Return(serviceNotes, 1, nil)

		req := httptest.NewRequest("GET", "/api/notes", nil)
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.GetNotes(w, req)

		require.Equal(t, http.StatusOK, w.Code, "Response body: %s", w.Body.String())

		var resp NotesResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Data, 1)
		if len(resp.Data) > 0 {
			assert.Equal(t, "note-1", resp.Data[0].ID)
		}
		assert.Equal(t, 1, resp.Meta.Total)

		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestCreateNoteHandler(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	title := "Test Note"
	contentMap := map[string]interface{}{"text": "hello"}
	contentJSON, _ := json.Marshal(contentMap)

	validReq := CreateNoteRequest{
		Title:   title,
		Content: contentMap,
	}

	createdNote := &services.Note{
		ID:        "note-123",
		UserID:    userUUID,
		Title:     title,
		Content:   contentJSON,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		dbMock, noteServiceMock, handler := setupTest(t)
		defer dbMock.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("CreateNote", mock.Anything, userUUID, title, mock.MatchedBy(func(c json.RawMessage) bool {
			return string(c) == string(contentJSON)
		})).Return(createdNote, nil)

		body, _ := json.Marshal(validReq)
		req := httptest.NewRequest("POST", "/api/notes", bytes.NewBuffer(body))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.CreateNote(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})
}

func TestDeleteNoteHandler(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	noteID := "note-123"

	t.Run("success", func(t *testing.T) {
		dbMock, noteServiceMock, handler := setupTest(t)
		defer dbMock.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("DeleteNote", mock.Anything, userUUID, noteID).Return(nil)

		req := httptest.NewRequest("DELETE", "/api/notes/"+noteID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.DeleteNote(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"success": true}`, w.Body.String())
		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		dbMock, noteServiceMock, handler := setupTest(t)
		defer dbMock.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("DeleteNote", mock.Anything, userUUID, noteID).Return(models.ErrNotFound)

		req := httptest.NewRequest("DELETE", "/api/notes/"+noteID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.DeleteNote(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})

	t.Run("db error", func(t *testing.T) {
		dbMock, noteServiceMock, handler := setupTest(t)
		defer dbMock.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("DeleteNote", mock.Anything, userUUID, noteID).Return(errors.New("db error"))

		req := httptest.NewRequest("DELETE", "/api/notes/"+noteID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.DeleteNote(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})
}

func TestUpdateNoteHandler(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	noteID := "note-123"
	title := "Updated Title"
	contentMap := map[string]interface{}{"text": "updated"}
	contentJSON, _ := json.Marshal(contentMap)

	validReq := CreateNoteRequest{
		Title:   title,
		Content: contentMap,
	}

	t.Run("success", func(t *testing.T) {
		dbMock, noteServiceMock, handler := setupTest(t)
		defer dbMock.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("UpdateNote", mock.Anything, userUUID, noteID, title, mock.MatchedBy(func(c json.RawMessage) bool {
			return string(c) == string(contentJSON)
		})).Return(nil)

		body, _ := json.Marshal(validReq)
		req := httptest.NewRequest("PUT", "/api/notes/"+noteID, bytes.NewBuffer(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.UpdateNote(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"success": true}`, w.Body.String())
		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		dbMock, noteServiceMock, handler := setupTest(t)
		defer dbMock.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("UpdateNote", mock.Anything, userUUID, noteID, title, mock.Anything).Return(models.ErrNotFound)

		body, _ := json.Marshal(validReq)
		req := httptest.NewRequest("PUT", "/api/notes/"+noteID, bytes.NewBuffer(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.UpdateNote(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})
}

func TestGetNoteHandler(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	noteID := "note-123"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		dbMock, noteServiceMock, handler := setupTest(t)
		defer dbMock.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		serviceNote := &services.Note{
			ID:        noteID,
			UserID:    userUUID,
			Title:     "Title",
			Content:   json.RawMessage(`{"text": "hello"}`),
			CreatedAt: now,
			UpdatedAt: now,
		}

		noteServiceMock.On("GetNote", mock.Anything, userUUID, noteID).Return(serviceNote, nil)

		req := httptest.NewRequest("GET", "/api/notes/"+noteID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.GetNote(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp Note
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, noteID, resp.ID)

		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		dbMock, noteServiceMock, handler := setupTest(t)
		defer dbMock.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("GetNote", mock.Anything, userUUID, noteID).Return(nil, models.ErrNotFound)

		req := httptest.NewRequest("GET", "/api/notes/"+noteID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.GetNote(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})
}
