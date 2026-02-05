package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"firebase.google.com/go/v4/auth"
)

func TestChatHandler_ChatWithAI(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)
		mockNotificationService := new(MockNotificationService)

		handler := NewChatHandler(mockClient, mockNoteService, mockNotificationService, nil)

		payload := map[string]interface{}{
			"prompt": "Hello",
		}
		body, _ := json.Marshal(payload)

		mockClient.On("ChatCompletion", mock.Anything, mock.Anything).Return(map[string]interface{}{
			"response": "Hello there",
		}, nil)

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))

		// Add auth context
		token := &auth.Token{UID: "user-123"}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestChatHandler_AskAI(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		handler := NewChatHandler(mockClient, nil, nil, nil)

		payload := map[string]interface{}{
			"prompt": "What is faith?",
		}
		body, _ := json.Marshal(payload)

		mockClient.On("ChatCompletion", mock.Anything, mock.Anything).Return(map[string]interface{}{
			"response": "Faith is...",
		}, nil)

		req := httptest.NewRequest("POST", "/api/ai/ask", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.AskAI(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
