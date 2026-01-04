package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPassage(t *testing.T) {
	// Mock server that mimics the external Bible AI API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/query", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "test-key", r.Header.Get("X-API-KEY"))

		// Check request body
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		query := req["query"].(map[string]interface{})
		assert.Equal(t, []interface{}{"John 3:16"}, query["verses"])

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"verse": "John 3:16 (ESV) For God so loved the world...",
		})
	}))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "test-key")

	result, err := client.GetPassage(context.Background(), "John 3:16")
	assert.NoError(t, err)
	assert.Equal(t, "John 3:16 (ESV) For God so loved the world...", result["verse"])
	assert.Equal(t, "John 3:16 (ESV) For God so loved the world...", result["text"])
	assert.Equal(t, "John 3:16 (ESV) For God so loved the world...", result["content"])
}

func TestChatCompletion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/query", r.URL.Path)

		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		// Return standard OQueryResponse
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"text": "This is the answer.",
			"references": []interface{}{
				map[string]interface{}{"verse": "John 3:16"},
			},
		})
	}))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "test-key")
	payload := map[string]interface{}{
		"prompt": "Explain John 3:16",
		"verses": []string{"John 3:16"},
	}

	result, err := client.ChatCompletion(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, "This is the answer.", result["text"])

	// Check backward compatibility wrapper
	choices := result["choices"].([]interface{})
	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	assert.Equal(t, "This is the answer.", message["content"])
}
