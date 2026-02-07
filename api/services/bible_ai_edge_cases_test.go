package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRealBibleAIClient_InvalidJSON(t *testing.T) {
	// Should warn but not panic
	client := NewRealBibleAIClient("url", "key", "{invalid-json")
	assert.NotNil(t, client)
	assert.Empty(t, client.SystemPrompts)
}

func TestRealBibleAIClient_GetPassage_EdgeCases(t *testing.T) {
	t.Run("Invalid_JSON_Response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain") // Force Resty to skip auto-unmarshal
			// Return 200 OK but with invalid JSON body
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{invalid-json`))
		}))
		defer server.Close()

		client := NewRealBibleAIClient(server.URL, "key", "")
		// Should return success with empty verse/text because parsing failed and fallback failed
		// Or maybe it returns empty map?
		// The logic:
		// if result.Verse == "" { log... json.Unmarshal(body)... if v... else ... log }
		// verseText := html.UnescapeString(result.Verse) -> ""
		// returns map with empty strings.

		res, err := client.GetPassage(context.Background(), "John 3:16", "ESV")
		assert.NoError(t, err)
		assert.Equal(t, "", res["verse"])
		assert.Equal(t, "", res["text"])
	})

	t.Run("Missing_Verse_And_Text", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// JSON is valid but has no relevant fields
			_, _ = w.Write([]byte(`{"other": "field"}`))
		}))
		defer server.Close()

		client := NewRealBibleAIClient(server.URL, "key", "")

		res, err := client.GetPassage(context.Background(), "John 3:16", "ESV")
		assert.NoError(t, err)
		assert.Equal(t, "", res["verse"])
	})
}

func TestRealBibleAIClient_Stream_ErrorPropagation(t *testing.T) {
	// This test ensures that when streaming fails and fallback fails, the error is propagated/logged
	// which hits the slog.Error line in the goroutine.

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always return 500
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")

	// Stream starts a goroutine.
	outChan, _, err := client.Stream(context.Background(), "prompt")
	assert.NoError(t, err)

	// Consume outChan until closed
	// Because of our fix, if errChan has error, it should be processed.
	// We can't strictly assert logs without hooking logger, but we can ensure execution completes.

	timeout := time.After(2 * time.Second)
	done := make(chan bool)

	go func() {
		for range outChan {
			// consume
		}
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-timeout:
		t.Fatal("Stream timed out")
	}
}
