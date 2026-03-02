package handlers

import (
	"bytes"
	"context"
	"database/sql"
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
	"github.com/DATA-DOG/go-sqlmock"
	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *MockNoteService, *NoteHandler) {
	db, dbMock, err := sqlmock.New()
	_ = db

	_ = dbMock
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	noteServiceMock := new(MockNoteService)
	handler := NewNoteHandler(db, noteServiceMock)
	return db, dbMock, noteServiceMock, handler
}

func TestGetNotes(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

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

	t.Run("with_params", func(t *testing.T) {
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		serviceNotes := []services.Note{}

		startDateStr := "2023-01-01T00:00:00Z"
		endDateStr := "2023-12-31T23:59:59Z"
		startDate, _ := time.Parse(time.RFC3339, startDateStr)
		endDate, _ := time.Parse(time.RFC3339, endDateStr)

		expectedFilter := services.NoteFilter{
			SearchQuery: "term",
			Tag:         "tag1",
			SortBy:      "title",
			SortOrder:   "desc",
			StartDate:   &startDate,
			EndDate:     &endDate,
		}

		noteServiceMock.On("GetNotes", mock.Anything, userUUID, 2, 50, expectedFilter).Return(serviceNotes, 0, nil)

		req := httptest.NewRequest("GET", "/api/notes?page=2&limit=50&q=term&tag=tag1&sortBy=title&sortOrder=desc&startDate="+startDateStr+"&endDate="+endDateStr, nil)
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.GetNotes(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		noteServiceMock.AssertExpectations(t)
	})

	t.Run("invalid params", func(t *testing.T) {
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		// Should default to empty filter if parsing fails? Or ignore invalid fields?
		// Logic in GetNotes handler:
		// startDate, _ := time.Parse(...) // Ignores error
		// limit, err := strconv.Atoi(...) // If error, defaults to 20 or ignores?

		// Let's assume it calls service with defaults/partial filter
		noteServiceMock.On("GetNotes", mock.Anything, userUUID, 1, 20, services.NoteFilter{}).Return([]services.Note{}, 0, nil)

		req := httptest.NewRequest("GET", "/api/notes?page=invalid&limit=invalid&startDate=invalid", nil)
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.GetNotes(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		noteServiceMock.AssertExpectations(t)
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
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("CreateNote", mock.Anything, userUUID, title, mock.MatchedBy(func(c json.RawMessage) bool {
			return string(c) == string(contentJSON)
		}), mock.Anything, mock.Anything).Return(createdNote, nil)

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
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

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
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

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
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

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
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("UpdateNote", mock.Anything, userUUID, noteID, title, mock.MatchedBy(func(c json.RawMessage) bool {
			return string(c) == string(contentJSON)
		}), mock.Anything, mock.Anything).Return(nil)

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
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("UpdateNote", mock.Anything, userUUID, noteID, title, mock.Anything, mock.Anything, mock.Anything).Return(models.ErrNotFound)

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
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

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
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

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

func TestCreateNoteHandler_InvalidBody(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"

	t.Run("invalid json", func(t *testing.T) {
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		req := httptest.NewRequest("POST", "/api/notes", bytes.NewBufferString("invalid-json"))
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.CreateNote(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		noteServiceMock.AssertNotCalled(t, "CreateNote")
	})
}

func TestUpdateNoteHandler_InvalidBody(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	noteID := "note-123"

	t.Run("invalid json", func(t *testing.T) {
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		req := httptest.NewRequest("PUT", "/api/notes/"+noteID, bytes.NewBufferString("invalid-json"))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", noteID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.UpdateNote(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		noteServiceMock.AssertNotCalled(t, "UpdateNote")
	})
}

func TestGetNotes_Unauthorized(t *testing.T) {
	db, dbMock, noteServiceMock, handler := setupTest(t)
	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	req := httptest.NewRequest("GET", "/api/notes", nil)
	// No token, no TestUserKey

	w := httptest.NewRecorder()
	handler.GetNotes(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	noteServiceMock.AssertNotCalled(t, "GetNotes")
}

func TestGetNotes_UserNotFound(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	db, dbMock, noteServiceMock, handler := setupTest(t)
	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnError(errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/notes", nil)
	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetNotes(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	noteServiceMock.AssertNotCalled(t, "GetNotes")
}

func TestCreateNote_UserNotFound(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	db, dbMock, noteServiceMock, handler := setupTest(t)
	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnError(errors.New("db error"))

	req := httptest.NewRequest("POST", "/api/notes", bytes.NewBufferString("{}"))
	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateNote(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	noteServiceMock.AssertNotCalled(t, "CreateNote")
}

func TestUpdateNote_UserNotFound(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	db, dbMock, noteServiceMock, handler := setupTest(t)
	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnError(errors.New("db error"))

	req := httptest.NewRequest("PUT", "/api/notes/1", bytes.NewBufferString("{}"))
	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.UpdateNote(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	noteServiceMock.AssertNotCalled(t, "UpdateNote")
}

func TestDeleteNote_UserNotFound(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	db, dbMock, noteServiceMock, handler := setupTest(t)
	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnError(errors.New("db error"))

	req := httptest.NewRequest("DELETE", "/api/notes/1", nil)
	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteNote(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	noteServiceMock.AssertNotCalled(t, "DeleteNote")
}

func TestGetNote_UserNotFound(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	db, dbMock, noteServiceMock, handler := setupTest(t)
	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnError(errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/notes/1", nil)
	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetNote(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	noteServiceMock.AssertNotCalled(t, "GetNote")
}

func TestGetNoteHandler_DBError(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	noteID := "note-123"

	db, dbMock, noteServiceMock, handler := setupTest(t)

	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

	noteServiceMock.On("GetNote", mock.Anything, userUUID, noteID).Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/notes/"+noteID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", noteID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetNote(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateNoteHandler_ServiceError(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	title := "Test Note"
	contentMap := map[string]interface{}{"text": "hello"}
	contentJSON, _ := json.Marshal(contentMap)

	validReq := CreateNoteRequest{
		Title:   title,
		Content: contentMap,
	}

	db, dbMock, noteServiceMock, handler := setupTest(t)

	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

	noteServiceMock.On("CreateNote", mock.Anything, userUUID, title, mock.MatchedBy(func(c json.RawMessage) bool {
		return string(c) == string(contentJSON)
	}), mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

	body, _ := json.Marshal(validReq)
	req := httptest.NewRequest("POST", "/api/notes", bytes.NewBuffer(body))
	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateNote(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateNoteHandler_ServiceError(t *testing.T) {
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

	db, dbMock, noteServiceMock, handler := setupTest(t)

	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

	noteServiceMock.On("UpdateNote", mock.Anything, userUUID, noteID, title, mock.MatchedBy(func(c json.RawMessage) bool {
		return string(c) == string(contentJSON)
	}), mock.Anything, mock.Anything).Return(errors.New("service error"))

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

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateNote_Unauthorized(t *testing.T) {
	db, _, noteServiceMock, handler := setupTest(t)
	_ = db
	req := httptest.NewRequest("POST", "/api/notes", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	handler.CreateNote(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	noteServiceMock.AssertNotCalled(t, "CreateNote")
}

func TestUpdateNote_Unauthorized(t *testing.T) {
	db, _, noteServiceMock, handler := setupTest(t)
	_ = db
	req := httptest.NewRequest("PUT", "/api/notes/1", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	handler.UpdateNote(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	noteServiceMock.AssertNotCalled(t, "UpdateNote")
}

func TestDeleteNote_Unauthorized(t *testing.T) {
	db, _, noteServiceMock, handler := setupTest(t)
	_ = db
	req := httptest.NewRequest("DELETE", "/api/notes/1", nil)
	w := httptest.NewRecorder()
	handler.DeleteNote(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	noteServiceMock.AssertNotCalled(t, "DeleteNote")
}

func TestGetNote_Unauthorized(t *testing.T) {
	db, _, noteServiceMock, handler := setupTest(t)
	_ = db
	req := httptest.NewRequest("GET", "/api/notes/1", nil)
	w := httptest.NewRecorder()
	handler.GetNote(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	noteServiceMock.AssertNotCalled(t, "GetNote")
}

func TestGetTags_Unauthorized(t *testing.T) {
	db, _, noteServiceMock, handler := setupTest(t)
	_ = db
	req := httptest.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()
	handler.GetTags(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	noteServiceMock.AssertNotCalled(t, "GetUserTags")
}

func TestCreateTag_Unauthorized(t *testing.T) {
	db, _, noteServiceMock, handler := setupTest(t)
	_ = db
	req := httptest.NewRequest("POST", "/api/tags", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	handler.CreateTag(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	noteServiceMock.AssertNotCalled(t, "CreateTag")
}

func TestDeleteTag_Unauthorized(t *testing.T) {
	db, _, noteServiceMock, handler := setupTest(t)
	_ = db
	req := httptest.NewRequest("DELETE", "/api/tags/1", nil)
	w := httptest.NewRecorder()
	handler.DeleteTag(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	noteServiceMock.AssertNotCalled(t, "DeleteTag")
}
