package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
		client = services.NewRealBibleAIClient(os.Getenv("BIBLE_API_URL"), os.Getenv("BIBLE_API_KEY"))
	} else {
		client = services.NewMockBibleAIClient()
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mockService := new(MockNoteService)
	handler := NewChatHandler(client, mockService, mockDB)

	t.Run("AskAI Success", func(t *testing.T) {
		reqBody := map[string]string{
			"context": "I am feeling happy today.",
			"prompt":  "What is happiness?",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/ai/ask", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.AskAI(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusOK, rr.Body.String())
		}

		var resp ChatResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if useReal {
			if len(resp.Response) == 0 {
				t.Error("Real API returned empty response")
			}
		} else {
			// The mock returns "This is a mocked AI response to: " + prompt
			// With the new implementation, AskAI sends req.Prompt as "prompt"
			expected := "This is a mocked AI response to: " + reqBody["prompt"]
			if resp.Response != expected {
				t.Errorf("Got unexpected mock response: %v\nExpected: %v", resp.Response, expected)
			}
		}
	})

	t.Run("ChatWithAI Success", func(t *testing.T) {
		if useReal {
			t.Skip("Skipping ChatWithAI test in real integration mode due to auth/DB complexity")
		}

		// Setup Mock DB expectation
		uid := "firebase_uid_123"
		userUUID := uuid.New()

		// pgxmock requires escaping $ for args if regex is used, but by default ExpectQuery takes a regex string
		// So we should escape $.
		mockDB.ExpectQuery("SELECT id FROM users WHERE firebase_uid=\\$1").
			WithArgs(uid).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

		// Setup Mock NoteService expectation
		// We expect CreateNote to be called with the userUUID we just returned
		mockService.On("CreateNote", mock.Anything, userUUID.String(), mock.Anything, mock.Anything).
			Return(&services.Note{ID: uuid.New().String()}, nil).
			Once()

		reqBody := map[string]interface{}{
			"passage": "John 3:16",
			"themes":  []string{"Love"},
			"prompt":  "What does this mean?",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/chat", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Inject Auth Context
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: uid})
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.ChatWithAI(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusOK, rr.Body.String())
		}

		// Ensure all expectations were met
		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		mockService.AssertExpectations(t)
	})
}
