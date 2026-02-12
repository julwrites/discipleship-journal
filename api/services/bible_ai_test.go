package services

import (
	"testing"
)

func TestNewRealBibleAIClient_SystemPrompts(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		wantKey   string
		wantVal   string
	}{
		{
			name:      "Valid JSON",
			jsonInput: `{"ask": "test"}`,
			wantKey:   "ask",
			wantVal:   "test",
		},
		{
			name: "Invalid JSON with newline in string (Fixable)",
			jsonInput: `{"ask": "line 1
line 2"}`,
			wantKey: "ask",
			wantVal: "line 1\nline 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewRealBibleAIClient("http://localhost", "key", tt.jsonInput)
			if client == nil {
				t.Fatal("NewRealBibleAIClient returned nil")
			}

			val, ok := client.SystemPrompts[tt.wantKey]
			if !ok {
				t.Errorf("SystemPrompts[%q] not found", tt.wantKey)
			} else if val != tt.wantVal {
				t.Errorf("SystemPrompts[%q] = %q, want %q", tt.wantKey, val, tt.wantVal)
			}
		})
	}
}

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
			name:  "Empty list items with newlines",
			input: "<ul><li>Item 1</li>\n<li>   </li>\n<li>Item 2</li></ul>",
			// Explanation of change:
			// 1. Newlines -> Space: "<ul><li>Item 1</li> <li>   </li> <li>Item 2</li></ul>"
			// 2. Empty Li -> Removes `<li>   </li>`. Result: "<ul><li>Item 1</li>  <li>Item 2</li></ul>"
			// 3. List Whitespace -> Collapses "</li>  <li>". Result: "<ul><li>Item 1</li><li>Item 2</li></ul>"
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
