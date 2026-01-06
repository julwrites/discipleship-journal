package services

import (
	"testing"
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
			name:     "Empty list items with newlines (which become spaces)",
			input:    "<ul><li>Item 1</li>\n<li>   </li>\n<li>Item 2</li></ul>",
			// The newlines become spaces. So "<ul><li>Item 1</li> <li>   </li> <li>Item 2</li></ul>"
			// Then empty li regex removes " <li>   </li> ".
			// Wait, the regex matches `<li`... it doesn't match spaces *between* tags.
			// Input: "<ul><li>Item 1</li>\n<li>   </li>\n<li>Item 2</li></ul>"
			// 1. Newlines -> Space: "<ul><li>Item 1</li> <li>   </li> <li>Item 2</li></ul>"
			// 2. Empty Para -> No change.
			// 3. Empty Li -> "<li>   </li>" matches.
			// Result: "<ul><li>Item 1</li>  <li>Item 2</li></ul>" (Two spaces between items)
			expected: "<ul><li>Item 1</li>  <li>Item 2</li></ul>",
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
			// Since newline regex replaces newlines with space, we should be careful with exact matching
			// but for our test cases it should be deterministic.
			if got != tt.expected {
				t.Errorf("cleanHTML() = %q, want %q", got, tt.expected)
			}
		})
	}
}
