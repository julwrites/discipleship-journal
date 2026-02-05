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
	t.Run("Success", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)
		mockNotificationService := new(MockNotificationService)

		handler := NewChatHandler(mockClient, mockNoteService, mockNotificationService, nil)

        // Use valid UUID for test user
        testUserID := "00000000-0000-0000-0000-000000000001"

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

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))

		// Use TestUserKey to bypass DB lookup
		ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestChatHandler_AskAI(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
        mockNoteService := new(MockNoteService)
		handler := NewChatHandler(mockClient, mockNoteService, nil, nil)

        testUserID := "00000000-0000-0000-0000-000000000001"

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
