package services

import (
	"context"
	"fmt"
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
			name:     "Basic text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Empty paragraphs",
			input:    "<p>Line 1</p><p>&nbsp;</p><p>Line 2</p>",
			expected: "<p>Line 1</p><p>Line 2</p>",
		},
		{
			name:     "Empty list items",
			input:    "<ul><li>Item 1</li><li></li><li>Item 2</li><li>&nbsp;</li></ul>",
			expected: "<ul><li>Item 1</li><li>Item 2</li></ul>",
		},
		{
			name:     "Empty list items with newlines",
			input:    "<ul><li>Item 1</li>\n<li>   </li>\n<li>Item 2</li></ul>",
			expected: "<ul><li>Item 1</li><li>Item 2</li></ul>",
		},
		{
			name:     "Mixed empty content",
			input:    "<p>Start</p><p><br></p><ul><li></li><li>Valid</li></ul>",
			expected: "<p>Start</p><ul><li>Valid</li></ul>",
		},
		{
			name:     "Empty list items with br",
			input:    "<ul><li><br></li><li>Item</li></ul>",
			expected: "<ul><li>Item</li></ul>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanHTML(tt.input)
			if got != tt.expected {
				t.Errorf("cleanHTML() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestBibleAIClient_GetVersions(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bible-versions" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"versions": ["ESV", "NIV"]}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	res, err := client.GetVersions(ctx, nil)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	versions, ok := res["versions"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, versions, 2)
}

func TestBibleAIClient_GetVersions_Error(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	res, err := client.GetVersions(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestBibleAIClient_GetSystemPrompt(t *testing.T) {
	prompts := `{"ask": "You are a helpful assistant."}`
	client := NewRealBibleAIClient("url", "key", prompts)

	p := client.GetSystemPrompt("ask")
	assert.Equal(t, "You are a helpful assistant.", p)

	p2 := client.GetSystemPrompt("unknown")
	assert.Equal(t, "", p2)
}
