package handlers

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestGetNotes(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	mockService := new(MockNoteService)
	handler := NewNoteHandler(mockDB, mockService)

	uid := "firebase-uid-123"
	userUUID := "user-uuid-123"

	// Mock getUserUUID query
	mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
		WithArgs(uid).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// Mock count query
	mockDB.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes WHERE user_id=").
		WithArgs(userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

	// Mock notes query
	now := time.Now()
	// Note: using map[string]interface{} for content as the handler scans into a map
	mockDB.ExpectQuery("SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id=").
		WithArgs(userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "updated_at"}).
			AddRow("note-1", userUUID, "Title 1", map[string]interface{}{}, now, now))

	req := httptest.NewRequest("GET", "/api/notes", nil)
	token := &auth.Token{UID: uid}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.GetNotes(w, req)

	require.Equal(t, http.StatusOK, w.Code, "Response body: %s", w.Body.String())

	var resp NotesResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Data, 1)
	if len(resp.Data) > 0 {
		assert.Equal(t, "note-1", resp.Data[0].ID)
	}
	assert.Equal(t, 1, resp.Meta.Total)
	assert.Equal(t, 1, resp.Meta.TotalPages)

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetNote(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	mockService := new(MockNoteService)
	handler := NewNoteHandler(mockDB, mockService)

	uid := "firebase-uid-123"
	userUUID := "user-uuid-123"
	noteID := "note-1"
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		// Mock getUserUUID query
		mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
			WithArgs(uid).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		// Mock GetNote query
		mockDB.ExpectQuery("SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id=").
			WithArgs(noteID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "updated_at"}).
				AddRow(noteID, userUUID, "Title 1", map[string]interface{}{"text": "content"}, now, now))

		req := httptest.NewRequest("GET", "/api/notes/"+noteID, nil)
		token := &auth.Token{UID: uid}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.GetNote(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp Note
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, noteID, resp.ID)
		assert.Equal(t, "Title 1", resp.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Mock getUserUUID query
		mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
			WithArgs(uid).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		// Mock GetNote query
		mockDB.ExpectQuery("SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id=").
			WithArgs(noteID).
			WillReturnError(pgx.ErrNoRows)

		req := httptest.NewRequest("GET", "/api/notes/"+noteID, nil)
		token := &auth.Token{UID: uid}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.GetNote(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		otherUserUUID := "other-user-uuid"

		// Mock getUserUUID query
		mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
			WithArgs(uid).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		// Mock GetNote query - returns note belonging to other user
		mockDB.ExpectQuery("SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id=").
			WithArgs(noteID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "updated_at"}).
				AddRow(noteID, otherUserUUID, "Title 1", map[string]interface{}{"text": "content"}, now, now))

		req := httptest.NewRequest("GET", "/api/notes/"+noteID, nil)
		token := &auth.Token{UID: uid}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.GetNote(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestCreateNote(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	mockService := new(MockNoteService)
	handler := NewNoteHandler(mockDB, mockService)

	uid := "firebase-uid-123"
	userUUID := "user-uuid-123"

	// Mock getUserUUID query
	mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
		WithArgs(uid).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// Mock insert note using Service
	reqBody := CreateNoteRequest{
		Title:   "New Note",
		Content: map[string]interface{}{"text": "content"},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// NoteService.CreateNote expectation
	mockService.On("CreateNote", mock.Anything, userUUID, reqBody.Title, mock.Anything).
		Return(&services.Note{ID: "new-note-id"}, nil)

	req := httptest.NewRequest("POST", "/api/notes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	token := &auth.Token{UID: uid}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.CreateNote(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "new-note-id", resp["id"])

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	mockService.AssertExpectations(t)
}

func TestUpdateNote(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	mockService := new(MockNoteService)
	handler := NewNoteHandler(mockDB, mockService)

	uid := "firebase-uid-123"
	userUUID := "user-uuid-123"
	noteID := "note-1"

	t.Run("Success", func(t *testing.T) {
		// Mock getUserUUID query
		mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
			WithArgs(uid).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		// Mock ownership check
		mockDB.ExpectQuery("SELECT user_id FROM notes WHERE id=").
			WithArgs(noteID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id"}).AddRow(userUUID))

		// Mock update
		mockDB.ExpectExec("UPDATE notes SET title=\\$1, content=\\$2, updated_at=NOW\\(\\) WHERE id=\\$3").
			WithArgs("Updated Title", map[string]interface{}{"text": "updated content"}, noteID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		reqBody := CreateNoteRequest{
			Title:   "Updated Title",
			Content: map[string]interface{}{"text": "updated content"},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/api/notes/"+noteID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		token := &auth.Token{UID: uid}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.UpdateNote(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		otherUserUUID := "other-user-uuid"

		// Mock getUserUUID query
		mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
			WithArgs(uid).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		// Mock ownership check - returns other user
		mockDB.ExpectQuery("SELECT user_id FROM notes WHERE id=").
			WithArgs(noteID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id"}).AddRow(otherUserUUID))

		reqBody := CreateNoteRequest{
			Title:   "Updated Title",
			Content: map[string]interface{}{"text": "updated content"},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/api/notes/"+noteID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		token := &auth.Token{UID: uid}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.UpdateNote(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestDeleteNote(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	mockService := new(MockNoteService)
	handler := NewNoteHandler(mockDB, mockService)

	uid := "firebase-uid-123"
	userUUID := "user-uuid-123"
	noteID := "note-1"

	// Mock getUserUUID query
	mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
		WithArgs(uid).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// Mock delete note using Service
	mockService.On("DeleteNote", mock.Anything, userUUID, noteID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/notes/"+noteID, nil)
	token := &auth.Token{UID: uid}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	// Mock URL Param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", noteID)
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.DeleteNote(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	mockService.AssertExpectations(t)
}
