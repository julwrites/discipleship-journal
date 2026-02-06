package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
    "time"

	"errors"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockNote(id, userID string) *services.Note {
    return &services.Note{
        ID: id,
        UserID: userID,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func TestChatHandler_ChatWithAI(t *testing.T) {
    testUserID := "00000000-0000-0000-0000-000000000001"

	t.Run("Success_Blocking", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)
		mockNotificationService := new(MockNotificationService)

		handler := NewChatHandler(mockClient, mockNoteService, mockNotificationService, nil)

		payload := map[string]interface{}{
			"prompt": "Hello",
            "passage": "John 3:16",
            "themes": []string{"Love"},
            "options": map[string]bool{"stream": false},
		}
		body, _ := json.Marshal(payload)

		mockClient.On("ChatCompletion", mock.Anything, mock.Anything).Return(map[string]interface{}{
			"response": "Hello there",
		}, nil)

        // Mock CreateNote
        mockNoteService.On("CreateNote", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
            Return(createMockNote("note-1", testUserID), nil)

        // Mock UpdateNote
        mockNoteService.On("UpdateNote", mock.Anything, testUserID, "note-1", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
            Return(nil)

        // Mock Query (since default is blocking)
        mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return("Response from AI", "", nil)

        // Mock Notification
        mockNotificationService.On("SendNotification", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

    t.Run("Success_Streaming", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)
		mockNotificationService := new(MockNotificationService)

		handler := NewChatHandler(mockClient, mockNoteService, mockNotificationService, nil)

		payload := map[string]interface{}{
			"prompt": "Hello",
            "passage": "John 3:16",
            "themes": []string{"Love"},
            "options": map[string]bool{"stream": true},
		}
		body, _ := json.Marshal(payload)

        mockNoteService.On("CreateNote", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
            Return(createMockNote("note-stream", testUserID), nil)

        mockNoteService.On("UpdateNote", mock.Anything, testUserID, "note-stream", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
            Return(nil)

        // Mock Stream
        outChan := make(chan string, 2)
        outChan <- "Hello "
        outChan <- "World"
        close(outChan)

        mockClient.On("Stream", mock.Anything, mock.Anything).Return((<-chan string)(outChan), "", nil)

        // Mock Notification (called after stream)
        mockNotificationService.On("SendNotification", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
        assert.Contains(t, rr.Body.String(), "event: start")
        assert.Contains(t, rr.Body.String(), "Hello")
        assert.Contains(t, rr.Body.String(), "World")
        assert.Contains(t, rr.Body.String(), "event: done")
	})

	t.Run("RealUserAuth_Success", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)
		mockNotificationService := new(MockNotificationService)
		mockDB, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mockDB.Close()

		handler := NewChatHandler(mockClient, mockNoteService, mockNotificationService, mockDB)

		firebaseUID := "firebase-uid-123"
		userUUID := uuid.New()
		userUUIDStr := userUUID.String()

		// Mock DB User Lookup
		mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		payload := map[string]interface{}{
			"prompt": "Hello",
            "passage": "John 3:16",
            "themes": []string{"Love"},
            "options": map[string]bool{"stream": false},
		}
		body, _ := json.Marshal(payload)

		mockNoteService.On("CreateNote", mock.Anything, userUUIDStr, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(createMockNote("note-real", userUUIDStr), nil)
		mockNoteService.On("UpdateNote", mock.Anything, userUUIDStr, "note-real", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return("Response", "", nil)
		mockNotificationService.On("SendNotification", mock.Anything, userUUIDStr, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("CreateNote_Error", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)

		handler := NewChatHandler(mockClient, mockNoteService, nil, nil)

		payload := map[string]interface{}{
			"prompt": "Hello",
            "passage": "John 3:16",
		}
		body, _ := json.Marshal(payload)

		mockNoteService.On("CreateNote", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("db error"))

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Query_Error", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)

		handler := NewChatHandler(mockClient, mockNoteService, nil, nil)

		payload := map[string]interface{}{
			"prompt": "Hello",
            "passage": "John 3:16",
			"options": map[string]bool{"stream": false},
		}
		body, _ := json.Marshal(payload)

		mockNoteService.On("CreateNote", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(createMockNote("note-err", testUserID), nil)

		mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return("", "", errors.New("ai error"))

		// Expect UpdateNote to be called with failed status
		mockNoteService.On("UpdateNote", mock.Anything, testUserID, "note-err", mock.Anything, mock.Anything, mock.Anything, []string{"failed"}).
			Return(nil)

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Stream_Init_Error", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)

		handler := NewChatHandler(mockClient, mockNoteService, nil, nil)

		payload := map[string]interface{}{
			"prompt": "Hello",
            "passage": "John 3:16",
			"options": map[string]bool{"stream": true},
		}
		body, _ := json.Marshal(payload)

		mockNoteService.On("CreateNote", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(createMockNote("note-stream-err", testUserID), nil)

		mockClient.On("Stream", mock.Anything, mock.Anything).Return((<-chan string)(nil), "", errors.New("stream error"))

		// Expect UpdateNote to be called with failed status
		mockNoteService.On("UpdateNote", mock.Anything, testUserID, "note-stream-err", mock.Anything, mock.Anything, mock.Anything, []string{"failed"}).
			Return(nil)

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		// Handler writes event: error
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "event: error")
	})
}

func TestChatHandler_AskAI(t *testing.T) {
    testUserID := "00000000-0000-0000-0000-000000000001"

	t.Run("Success", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
        mockNoteService := new(MockNoteService)
		// No notification service for AskAI in some cases? Or maybe yes?
		// AskAI calls handleRequest which calls sendNotification if NoteService updates to active.
		// So yes, it needs notification mock too?
		// Check NewChatHandler signature in test: previously passed nil.
		// If nil, sendNotification check `if h.NotificationService != nil`.
		// So passing nil is safe.
		handler := NewChatHandler(mockClient, mockNoteService, nil, nil)

		payload := map[string]interface{}{
			"prompt": "What is faith?",
            "context": "Context here",
            "options": map[string]bool{"stream": false},
		}
		body, _ := json.Marshal(payload)

        mockNoteService.On("CreateNote", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
            Return(createMockNote("note-2", testUserID), nil)

        mockNoteService.On("UpdateNote", mock.Anything, testUserID, "note-2", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
            Return(nil)

		mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return("Faith is...", "", nil)

		req := httptest.NewRequest("POST", "/api/ai/ask", bytes.NewBuffer(body))
        ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
        req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.AskAI(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid_Body", func(t *testing.T) {
		handler := NewChatHandler(nil, nil, nil, nil)

		req := httptest.NewRequest("POST", "/api/ai/ask", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()

		handler.AskAI(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
