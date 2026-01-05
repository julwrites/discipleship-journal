package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"discipleship_journal_api/services"
)

func TestGetBiblePassage(t *testing.T) {
	// Decide whether to use real or mock client based on env var
	useReal := os.Getenv("TEST_REAL_BIBLE_API") == "true"
	var client services.BibleAIClient
	if useReal {
		client = services.NewRealBibleAIClient(os.Getenv("BIBLE_API_URL"), os.Getenv("BIBLE_API_KEY"), os.Getenv("LLM_SYSTEM_PROMPTS"))
	} else {
		client = services.NewMockBibleAIClient()
	}

	handler := NewBibleHandler(client)

	t.Run("Missing Reference", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/bible/passage", nil)
		rr := httptest.NewRecorder()

		handler.GetBiblePassage(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Valid Reference", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/bible/passage?ref=John+3:16", nil)
		rr := httptest.NewRecorder()

		handler.GetBiblePassage(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if useReal {
			// Real API checks (might be flaky if API changes)
			// API returns {"verse": "John 3:16 (ESV) For God so loved the world..."}
			if _, ok := resp["verse"]; !ok {
				t.Error("Real API response missing 'verse'")
			}
		} else {
			// Mock checks
			if resp["reference"] != "John 3:16" {
				t.Errorf("Expected reference John 3:16, got %v", resp["reference"])
			}
			if resp["text"] != "For God so loved the world... (Mocked)" {
				t.Errorf("Got unexpected mock text: %v", resp["text"])
			}
		}
	})
}
