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
	"discipleship_journal_api/services"

	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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
	// We can't easily defer dbMock.Close() here because it needs to live as long as the test
	// But since it's a mock, it's probably fine, or we can return a cleanup func.

	noteServiceMock := new(MockNoteService)
	handler := NewNoteHandler(dbMock, noteServiceMock)
	return dbMock, noteServiceMock, handler
}

func TestGetNotes(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		dbMock, _, handler := setupTest(t)
		defer dbMock.Close()

		// 1. Mock getUserUUID
		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		// 2. Mock Count Query
		dbMock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes WHERE user_id=\\$1").
			WithArgs(userUUID).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

		// 3. Mock Select Query
		// Note: content is map[string]interface{}, but stored as jsonb. pgx scans it.
		// We mock it as json.RawMessage or bytes? Or map?
		// In TestCreateNoteHandler I used matchedBy contentJSON.
		// Here we are returning rows. pgxmock accepts any values.
		content := map[string]interface{}{"text": "content"}

		dbMock.ExpectQuery("SELECT id, user_id, title, content, created_at, updated_at FROM notes").
			WithArgs(userUUID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "updated_at"}).
				AddRow("note-1", userUUID, "Title 1", content, now, now))

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
		assert.Equal(t, 1, resp.Meta.TotalPages)

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

		// Mock getUserUUID behavior
		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		// Mock NoteService behavior
		noteServiceMock.On("CreateNote", mock.Anything, userUUID, title, mock.MatchedBy(func(c json.RawMessage) bool {
			return string(c) == string(contentJSON)
		})).Return(createdNote, nil)

		// Create request
		body, _ := json.Marshal(validReq)
		req := httptest.NewRequest("POST", "/api/notes", bytes.NewBuffer(body))

		// Add auth context
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.CreateNote(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]string
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, "note-123", resp["id"])

		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		dbMock, _, handler := setupTest(t)
		defer dbMock.Close()

		// Mock getUserUUID behavior (called before validation)
		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		invalidReq := CreateNoteRequest{
			Title: "", // Invalid
		}
		body, _ := json.Marshal(invalidReq)
		req := httptest.NewRequest("POST", "/api/notes", bytes.NewBuffer(body))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.CreateNote(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
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

		noteServiceMock.On("DeleteNote", mock.Anything, userUUID, noteID).Return(pgx.ErrNoRows)

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
