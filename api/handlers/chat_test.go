package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
    "time"

	"discipleship_journal_api/services"
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
}
