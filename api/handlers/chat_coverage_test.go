package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChatHandler_AuthFailures(t *testing.T) {
	mockClient := new(MockBibleAIClient)
	mockNoteService := new(MockNoteService)
	mockNotificationService := new(MockNotificationService)

	// DB is needed for real token auth lookup
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	handler := NewChatHandler(mockClient, mockNoteService, mockNotificationService, mockDB)

	payload := map[string]interface{}{
		"prompt":  "Hello",
		"passage": "John 3:16",
		"themes":  []string{"Love"},
	}
	body, _ := json.Marshal(payload)

	t.Run("NoAuth", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Unauthorized")
	})

	t.Run("InvalidTestUserUUID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), TestUserKey, "invalid-uuid")
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid test user ID")
	})

	t.Run("RealAuth_DBError", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		token := &auth.Token{UID: "firebase-uid"}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
			WithArgs("firebase-uid").
			WillReturnError(errors.New("db error"))

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code) // Handler returns StatusNotFound on DB error (User not found)
		assert.Contains(t, rr.Body.String(), "User not found")
	})
}

func TestChatHandler_ValidationFailures(t *testing.T) {
	handler := NewChatHandler(nil, nil, nil, nil)

	t.Run("ChatWithAI_InvalidBody", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("AskAI_InvalidBody", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/ai/ask", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()

		handler.AskAI(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestChatHandler_OptionsAndErrors(t *testing.T) {
	testUUID := uuid.New()

	t.Run("ChatWithAI_WithOptions", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)
		mockNotificationService := new(MockNotificationService)
		handler := NewChatHandler(mockClient, mockNoteService, mockNotificationService, nil)

		payload := map[string]interface{}{
			"prompt":  "Hello",
			"passage": "Gen 1:1",
			"themes":  []string{"Creation"},
			"version": "KJV",
			"options": map[string]interface{}{
				"stream":      false,
				"ai_provider": "openai",
			},
		}
		body, _ := json.Marshal(payload)

		mockNoteService.On("CreateNote", mock.Anything, testUUID.String(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(&services.Note{ID: "note-id"}, nil)

		// Expect options in context
		mockClient.On("Query", mock.MatchedBy(func(ctx context.Context) bool {
			v, _ := ctx.Value(services.BibleVersionKey).(string)
			p, _ := ctx.Value(services.AIProviderKey).(string)
			return v == "KJV" && p == "openai"
		}), mock.Anything, mock.Anything).Return("Answer", "", nil)

		mockNoteService.On("UpdateNote", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		mockNotificationService.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID.String())
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Notification_Error", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)
		mockNotificationService := new(MockNotificationService)
		handler := NewChatHandler(mockClient, mockNoteService, mockNotificationService, nil)

		payload := map[string]interface{}{
			"prompt": "Hello",
			"passage": "Gen 1:1",
			"themes": []string{},
			"options": map[string]bool{"stream": false},
		}
		body, _ := json.Marshal(payload)

		mockNoteService.On("CreateNote", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(&services.Note{ID: "note-id"}, nil)

		mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return("Answer", "", nil)

		mockNoteService.On("UpdateNote", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		// Notification returns error, handler should log it and continue (return 200)
		mockNotificationService.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("notification error"))

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID.String())
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Double_Failure", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		mockNoteService := new(MockNoteService)
		mockNotificationService := new(MockNotificationService)
		handler := NewChatHandler(mockClient, mockNoteService, mockNotificationService, nil)

		// Query fails, then UpdateNote(failed) also fails
		payload := map[string]interface{}{
			"prompt": "Hello",
			"passage": "Gen 1:1",
			"options": map[string]bool{"stream": false},
		}
		body, _ := json.Marshal(payload)

		mockNoteService.On("CreateNote", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(&services.Note{ID: "note-id"}, nil)

		mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).
			Return("", "", errors.New("ai error"))

		// First UpdateNote (mark failed) fails
		mockNoteService.On("UpdateNote", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, []string{"failed"}).
			Return(errors.New("update error"))

		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID.String())
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
