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
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestHandlerSafety checks for common panic scenarios in handlers
func TestHandlerSafety(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"

	// Scenario 1: NoteService is nil
	t.Run("nil_service", func(t *testing.T) {
		dbMock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer dbMock.Close()

		handler := NewNoteHandler(dbMock, nil) // Explicitly nil service

		req := httptest.NewRequest("GET", "/api/notes", nil)
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		assert.NotPanics(t, func() {
			handler.GetNotes(w, req)
		})

		// It should probably return 500 or handle it gracefully
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// Scenario 2: DB is nil (GetUserUUID will fail)
	t.Run("nil_db", func(t *testing.T) {
		noteServiceMock := new(MockNoteService)
		handler := NewNoteHandler(nil, noteServiceMock) // Explicitly nil DB

		req := httptest.NewRequest("GET", "/api/notes", nil)
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		assert.NotPanics(t, func() {
			handler.GetNotes(w, req)
		})

		// If DB is nil, GetUserUUID should return error, handler should return 404/500
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	// Scenario 3: NoteService returns nil slice (should be fine, but verification)
	t.Run("nil_notes_slice_from_service", func(t *testing.T) {
		dbMock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer dbMock.Close()

		serviceMock := new(MockNoteService)
		handler := NewNoteHandler(dbMock, serviceMock)

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		// Return nil slice
		serviceMock.On("GetNotes", mock.Anything, userUUID, mock.Anything, mock.Anything, mock.Anything).
			Return(([]services.Note)(nil), 0, nil)

		req := httptest.NewRequest("GET", "/api/notes", nil)
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		assert.NotPanics(t, func() {
			handler.GetNotes(w, req)
		})

		assert.Equal(t, http.StatusOK, w.Code)
		var resp NotesResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Data) // Should be empty slice, not nil
		assert.Len(t, resp.Data, 0)
	})

    // Scenario 4: User not in context (should return 401, not panic)
    t.Run("missing_user_context", func(t *testing.T) {
        dbMock, err := pgxmock.NewPool()
        require.NoError(t, err)
        defer dbMock.Close()

        serviceMock := new(MockNoteService)
        handler := NewNoteHandler(dbMock, serviceMock)

        req := httptest.NewRequest("GET", "/api/notes", nil)
        // No user context

        w := httptest.NewRecorder()

        assert.NotPanics(t, func() {
            handler.GetNotes(w, req)
        })

        assert.Equal(t, http.StatusUnauthorized, w.Code)
    })

    // Scenario 5: Note content is nil (Service returns note with nil Content)
    t.Run("nil_content_in_note", func(t *testing.T) {
        dbMock, err := pgxmock.NewPool()
        require.NoError(t, err)
        defer dbMock.Close()

        serviceMock := new(MockNoteService)
        handler := NewNoteHandler(dbMock, serviceMock)

        dbMock.ExpectQuery("SELECT id FROM users").
            WithArgs(firebaseUID).
            WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

        serviceNotes := []services.Note{
            {
                ID:        "note-1",
                UserID:    userUUID,
                Title:     "Title 1",
                Content:   nil, // Nil content
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            },
        }

        serviceMock.On("GetNotes", mock.Anything, userUUID, mock.Anything, mock.Anything, mock.Anything).
            Return(serviceNotes, 1, nil)

        req := httptest.NewRequest("GET", "/api/notes", nil)
        token := &auth.Token{UID: firebaseUID}
        ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
        req = req.WithContext(ctx)

        w := httptest.NewRecorder()

        assert.NotPanics(t, func() {
            handler.GetNotes(w, req)
        })

        assert.Equal(t, http.StatusOK, w.Code)
    })

	// Scenario 6: Nil *auth.Token in context
	t.Run("nil_token_in_context", func(t *testing.T) {
		dbMock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer dbMock.Close()

		handler := NewNoteHandler(dbMock, new(MockNoteService))

		req := httptest.NewRequest("GET", "/api/notes", nil)
		// Inject nil *auth.Token
		var nilToken *auth.Token = nil
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, nilToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		assert.NotPanics(t, func() {
			handler.GetNotes(w, req)
		})

		// If it doesn't panic, it should probably be Unauthorized
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Scenario 7: CreateNote with nil token
	t.Run("create_note_nil_token", func(t *testing.T) {
		dbMock, err := pgxmock.NewPool()
		require.NoError(t, err)
		defer dbMock.Close()

		handler := NewNoteHandler(dbMock, new(MockNoteService))

		reqBody := `{"title": "Test", "content": {"text": "foo"}}`
		req := httptest.NewRequest("POST", "/api/notes", bytes.NewBufferString(reqBody))
		var nilToken *auth.Token = nil
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, nilToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		assert.NotPanics(t, func() {
			handler.CreateNote(w, req)
		})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
