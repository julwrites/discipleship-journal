package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"discipleship_journal_api/services"
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

	mockService := new(MockNoteService)
	handler := NewChatHandler(client, mockService)

	// Note: Testing ChatWithAI fully requires mocking Auth Middleware and DB context,
	// or refactoring the handler to be more testable.
	// The current handler depends heavily on:
	// 1. middleware.DecodeAndValidate
	// 2. middleware.UserContextKey (for uid)
	// 3. GetUserUUID -> DB lookup
	// 4. services.NoteService -> DB insert
	//
	// This makes unit testing the "Chat" logic hard without a full integration setup or deep mocking.
	// However, we can test `AskAI` more easily as it doesn't save to DB (based on the code I read).
	// Let's check AskAI again.
	// AskAI uses DecodeAndValidate but doesn't seem to access DB or Auth context for UserUUID.

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
			// The mock returns "This is a mocked AI response to: " + constructed_prompt
			// The handler constructs prompt as: fmt.Sprintf("Context: %s. Question: %s", req.Context, req.Prompt)
			expectedPrompt := "Context: I am feeling happy today.. Question: What is happiness?"
			expected := "This is a mocked AI response to: " + expectedPrompt
			if resp.Response != expected {
				t.Errorf("Got unexpected mock response: %v\nExpected: %v", resp.Response, expected)
			}
		}
	})
}
