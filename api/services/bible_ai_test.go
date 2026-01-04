package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePassageFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple paragraph",
			html:     `<p>For God so loved the world</p>`,
			expected: "For God so loved the world",
			wantErr:  false,
		},
		{
			name:     "bold text",
			html:     `<p>For God so <b>loved</b> the world</p>`,
			expected: "For God so **loved** the world",
			wantErr:  false,
		},
		{
			name:     "italic text",
			html:     `<p>For God so <i>loved</i> the world</p>`,
			expected: "For God so *loved* the world",
			wantErr:  false,
		},
		{
			name:     "strong and em",
			html:     `<p><strong>Important</strong> and <em>emphasis</em></p>`,
			expected: "**Important** and *emphasis*",
			wantErr:  false,
		},
		{
			name:     "superscript",
			html:     `<p>Jesus<sup>1</sup> said</p>`,
			expected: "Jesus^1^ said",
			wantErr:  false,
		},
		{
			name:     "line break",
			html:     `<p>Line one<br>Line two</p>`,
			expected: "Line one\nLine two",
			wantErr:  false,
		},
		{
			name:     "multiple paragraphs",
			html:     `<p>First paragraph</p><p>Second paragraph</p>`,
			expected: "First paragraph\n\nSecond paragraph",
			wantErr:  false,
		},
		{
			name:     "heading",
			html:     `<h1>Chapter 1</h1><p>Text</p>`,
			expected: "**Chapter 1**\n\nText",
			wantErr:  false,
		},
		{
			name:     "nested tags",
			html:     `<p><b><i>bold italic</i></b> text</p>`,
			expected: "***bold italic*** text",
			wantErr:  false,
		},
		{
			name:     "span passes through",
			html:     `<p><span class=\"chapternum\">1</span> Text</p>`,
			expected: "1 Text",
			wantErr:  false,
		},
		{
			name:     "empty html",
			html:     ``,
			expected: "",
			wantErr:  false,
		},
		{
			name:     "html with attributes",
			html:     `<p class=\"verse\">Verse <sup class=\"footnote\">1</sup> text</p>`,
			expected: "Verse ^1^ text",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePassageFromHTML(tt.html)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}