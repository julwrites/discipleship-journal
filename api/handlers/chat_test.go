package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/mock"
)

func TestChatWithAI(t *testing.T) {
	// Decide whether to use real or mock client based on env var
	useReal := os.Getenv("TEST_REAL_BIBLE_API") == "true"
	var client services.BibleAIClient
	if useReal {
		client = services.NewRealBibleAIClient(os.Getenv("BIBLE_API_URL"), os.Getenv("BIBLE_API_KEY"), os.Getenv("LLM_SYSTEM_PROMPTS"))
	} else {
		client = services.NewMockBibleAIClient()
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mockService := new(MockNoteService)
	mockNotificationService := services.NewMockNotificationService()

	handler := NewChatHandler(client, mockService, mockNotificationService, mockDB)

	t.Run("AskAI Success", func(t *testing.T) {
		testUserID := "00000000-0000-0000-0000-000000000001"

		mockService.On("CreateNote", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything, []string{"pending"}).
			Return(&services.Note{ID: "note-123"}, nil).
			Once()

		mockService.On("UpdateNote", mock.Anything, testUserID, "note-123", mock.Anything, mock.Anything, mock.Anything, []string{"active"}).
			Return(nil).
			Maybe()

		reqBody := map[string]string{
			"context": "I am feeling happy today.",
			"prompt":  "What is happiness?",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/ai/ask", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.AskAI(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusOK, rr.Body.String())
		}

		contentType := rr.Header().Get("Content-Type")
		if contentType != "text/event-stream" {
			t.Errorf("Expected Content-Type text/event-stream, got %s", contentType)
		}

		body := rr.Body.String()
		if !strings.Contains(body, "event: start") {
			t.Error("Expected event: start")
		}
		if !strings.Contains(body, "event: chunk") {
			t.Error("Expected event: chunk")
		}
		if !strings.Contains(body, "event: done") {
			t.Error("Expected event: done")
		}
	})

	t.Run("ChatWithAI Success", func(t *testing.T) {
		if useReal {
			t.Skip("Skipping ChatWithAI test in real integration mode due to auth/DB complexity")
		}

		uid := "firebase_uid_123"
		userUUID := uuid.New()

		mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=\\$1").
			WithArgs(uid).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		mockService.On("CreateNote", mock.Anything, userUUID.String(), mock.Anything, mock.Anything, mock.Anything, []string{"pending"}).
			Return(&services.Note{ID: "note-456"}, nil).
			Once()

		mockService.On("UpdateNote", mock.Anything, userUUID.String(), "note-456", mock.Anything, mock.Anything, mock.Anything, []string{"active"}).
			Return(nil).
			Maybe()

		reqBody := map[string]interface{}{
			"passage": "John 3:16",
			"themes":  []string{"Love"},
			"prompt":  "What does this mean?",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/chat", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: uid})
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusOK, rr.Body.String())
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		body := rr.Body.String()
		if !strings.Contains(body, "event: start") {
			t.Error("Expected event: start")
		}

		time.Sleep(10 * time.Millisecond)
	})
}
