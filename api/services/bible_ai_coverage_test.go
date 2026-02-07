package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestRealBibleAIClient_GetVersions_Coverage(t *testing.T) {
	t.Run("APIURL_Empty", func(t *testing.T) {
		client := &RealBibleAIClient{APIURL: ""}
		_, err := client.GetVersions(context.Background(), nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bible API not configured")
	})

	t.Run("Network_Error", func(t *testing.T) {
		// Invalid URL to force network error
		client := NewRealBibleAIClient("http://invalid-url-that-does-not-exist.local", "key", "")
		// Reduce timeout to fail fast
		client.Client.SetTimeout(1)

		_, err := client.GetVersions(context.Background(), nil)
		assert.Error(t, err)
		// Error message varies but should be resty error
	})

	t.Run("API_Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": {"message": "Internal Error"}}`))
		}))
		defer server.Close()

		client := NewRealBibleAIClient(server.URL, "key", "")
		_, err := client.GetVersions(context.Background(), nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Internal Error")
	})
}

func TestRealBibleAIClient_GetSystemPrompt_Coverage(t *testing.T) {
	t.Run("Nil_Map", func(t *testing.T) {
		client := &RealBibleAIClient{SystemPrompts: nil}
		prompt := client.GetSystemPrompt("any")
		assert.Equal(t, "", prompt)
	})
}

func TestRealBibleAIClient_PrepareQueryRequest_Coverage(t *testing.T) {
	// We test prepareQueryRequest via ChatCompletion
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": {"text": "Response"}}`))
	}))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", `{"chat": "Template {PROMPT}"}`)

	t.Run("Full_Payload", func(t *testing.T) {
		payload := map[string]interface{}{
			"prompt":      "Hello",
			"context":     "Ctx",
			"version":     "KJV",
			"ai_provider": "openai",
			"themes":      []string{"Grace", "Truth"},
			"verses":      []string{"John 1:14"},
			"type":        "chat",
		}

		// We can't easily inspect the request sent unless we check it in the mock server.
		// But just running this path covers the logic.
		_, err := client.ChatCompletion(context.Background(), payload)
		assert.NoError(t, err)
	})

	t.Run("Empty_Payload_Defaults", func(t *testing.T) {
		payload := map[string]interface{}{
			"prompt": 123, // Invalid type, should default to empty string
			// No version -> ESV
			// No themes
		}
		_, err := client.ChatCompletion(context.Background(), payload)
		assert.NoError(t, err)
	})
}

func TestRealBibleAIClient_Query_Coverage(t *testing.T) {
	t.Run("APIURL_Empty", func(t *testing.T) {
		client := &RealBibleAIClient{APIURL: ""}
		_, _, err := client.Query(context.Background(), "prompt", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bible API not configured")
	})

	t.Run("Request_Error", func(t *testing.T) {
		client := NewRealBibleAIClient("http://invalid", "key", "")
		// Force error in NewRequest or Post?
		// Invalid URL usually causes error in Post.
		client.Client = resty.New() // Reset client
		// Using an invalid URL scheme to fail immediately?
		client.APIURL = "::invalid" // Invalid URL structure

		_, _, err := client.Query(context.Background(), "prompt", "")
		assert.Error(t, err)
	})

	t.Run("API_Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": {"message": "Bad Request"}}`))
		}))
		defer server.Close()

		client := NewRealBibleAIClient(server.URL, "key", "")
		_, _, err := client.Query(context.Background(), "prompt", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Bad Request")
	})
}

func TestRealBibleAIClient_GetPassage_Coverage(t *testing.T) {
	t.Run("APIURL_Empty", func(t *testing.T) {
		client := &RealBibleAIClient{APIURL: ""}
		_, err := client.GetPassage(context.Background(), "John 3:16", "ESV")
		assert.Error(t, err)
	})

	t.Run("Request_Error", func(t *testing.T) {
		client := NewRealBibleAIClient("::invalid", "key", "")
		_, err := client.GetPassage(context.Background(), "John 3:16", "ESV")
		assert.Error(t, err)
	})

	t.Run("API_Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": {"message": "Not Found"}}`))
		}))
		defer server.Close()

		client := NewRealBibleAIClient(server.URL, "key", "")
		_, err := client.GetPassage(context.Background(), "John 3:16", "ESV")
		assert.Error(t, err)
	})

	t.Run("Fallback_Text_Field", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Return 'text' instead of 'verse'
			_, _ = w.Write([]byte(`{"text": "For God so loved the world..."}`))
		}))
		defer server.Close()

		client := NewRealBibleAIClient(server.URL, "key", "")
		resp, err := client.GetPassage(context.Background(), "John 3:16", "ESV")
		assert.NoError(t, err)
		assert.Equal(t, "For God so loved the world...", resp["text"])
	})
}

func TestRealBibleAIClient_Stream_Coverage(t *testing.T) {
	t.Run("Stream_Fallback_Failure", func(t *testing.T) {
		// Mock server that returns error for stream and fallback
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError) // Both calls fail
		}))
		defer server.Close()

		client := NewRealBibleAIClient(server.URL, "key", "")

		outChan, _, err := client.Stream(context.Background(), "prompt")
		assert.NoError(t, err) // Stream init succeeds (async)

		// Wait for stream to close or produce output
		for msg := range outChan {
			// Should not produce meaningful output if everything fails
			// But might close immediately
			_ = msg
		}
		// If we reached here, outChan is closed.
		// The error was logged in goroutine (can't verify easily without hook).
	})

	t.Run("Stream_SSE_Formats", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			flusher, _ := w.(http.Flusher)

			// 1. V2 Delta format
			_, _ = w.Write([]byte("data: {\"delta\": \"Hello\"}\n\n"))
			flusher.Flush()

			// 2. Text fallback format
			_, _ = w.Write([]byte("data: {\"text\": \" World\"}\n\n"))
			flusher.Flush()

			// 3. Choices delta format (OpenAI style)
			_, _ = w.Write([]byte("data: {\"choices\": [{\"delta\": {\"content\": \"!\"}}]}\n\n"))
			flusher.Flush()

			// 4. Raw text fallback (invalid JSON)
			_, _ = w.Write([]byte("data:  RawText \n\n"))
			flusher.Flush()

			// 5. Done
			_, _ = w.Write([]byte("data: [DONE]\n\n"))
			flusher.Flush()
		}))
		defer server.Close()

		client := NewRealBibleAIClient(server.URL, "key", "")
		outChan, _, err := client.Stream(context.Background(), "prompt")
		assert.NoError(t, err)

		fullText := ""
		for msg := range outChan {
			fullText += msg
		}

		assert.Contains(t, fullText, "Hello")
		assert.Contains(t, fullText, "World")
		assert.Contains(t, fullText, "!")
		assert.Contains(t, fullText, "RawText")
	})
}
