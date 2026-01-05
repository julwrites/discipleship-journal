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
