package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Removes newlines between tags",
			input:    "<p>Line 1</p>\n\n<p>Line 2</p>",
			expected: "<p>Line 1</p> <p>Line 2</p>", // Collapses multiple newlines to single space
		},
		{
			name:     "Removes empty paragraphs",
			input:    "<p>Line 1</p><p></p><p>Line 2</p>",
			expected: "<p>Line 1</p><p>Line 2</p>",
		},
		{
			name:     "Removes paragraphs with whitespace",
			input:    "<p>Line 1</p><p>  </p><p>Line 2</p>",
			expected: "<p>Line 1</p><p>Line 2</p>",
		},
		{
			name:     "Removes paragraphs with br",
			input:    "<p>Line 1</p><p><br></p><p>Line 2</p>",
			expected: "<p>Line 1</p><p>Line 2</p>",
		},
		{
			name:     "Removes paragraphs with self-closing br",
			input:    "<p>Line 1</p><p><br/></p><p>Line 2</p>",
			expected: "<p>Line 1</p><p>Line 2</p>",
		},
		{
			name:     "Removes paragraphs with &nbsp;",
			input:    "<p>Line 1</p><p>&nbsp;</p><p>Line 2</p>",
			expected: "<p>Line 1</p><p>Line 2</p>",
		},
		{
			name:     "Replaces newlines with space in text",
			input:    "<p>Line\n1</p>",
			expected: "<p>Line 1</p>",
		},
		{
			name:     "Complex mixed case",
			input:    "<p>Start</p>\n\n<p></p>\n<p>End</p>",
			expected: "<p>Start</p>  <p>End</p>", // Collapsed newlines result in spaces
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanHTML(tt.input)
			if got != tt.expected {
				t.Errorf("cleanHTML(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRealBibleAIClient_Unescape(t *testing.T) {
	// Create a mock server that returns escaped HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/query" {
			var req QueryRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return
			}

			// Check if it's a verse request or chat/prompt
			if len(req.Query.Verses) > 0 {
				// Return escaped HTML for verse
				_ = json.NewEncoder(w).Encode(VerseResponse{
					Verse: "&lt;h3&gt;Title&lt;/h3&gt;&lt;p&gt;Verse Text&lt;/p&gt;",
				})
			} else {
				// Return escaped HTML for prompt
				_ = json.NewEncoder(w).Encode(OQueryResponse{
					Text: "&lt;p&gt;AI Response&lt;/p&gt;",
				})
			}
		}
	}))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "test-key", "")

	t.Run("GetPassage unescapes HTML", func(t *testing.T) {
		res, err := client.GetPassage(context.Background(), "John 3:16")
		assert.NoError(t, err)
		assert.Equal(t, "<h3>Title</h3><p>Verse Text</p>", res["verse"])
	})

	t.Run("ChatCompletion unescapes HTML", func(t *testing.T) {
		res, err := client.ChatCompletion(context.Background(), map[string]interface{}{
			"prompt": "test",
		})
		assert.NoError(t, err)
		assert.Equal(t, "<p>AI Response</p>", res["text"])
	})
}
