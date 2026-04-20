package services

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBibleAIClient_VerseReferenceRegex(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedRef   string
		expectedVer   string
		expectedText  string
		expectedMatch bool
	}{
		{
			name:          "Standard verse",
			input:         "John 3:16 (ESV) For God so loved the world...",
			expectedRef:   "John 3:16",
			expectedVer:   "ESV",
			expectedText:  "For God so loved the world...",
			expectedMatch: true,
		},
		{
			name:          "Verse range",
			input:         "1 John 3:16-18 (NIV) This is how we know...",
			expectedRef:   "1 John 3:16-18",
			expectedVer:   "NIV",
			expectedText:  "This is how we know...",
			expectedMatch: true,
		},
		{
			name:          "Whole chapter",
			input:         "Romans 3 (ESV) Then what advantage has the Jew? Or what is the value of circumcision? Much in every way. To begin with...",
			expectedRef:   "Romans 3",
			expectedVer:   "ESV",
			expectedText:  "Then what advantage has the Jew? Or what is the value of circumcision? Much in every way. To begin with...",
			expectedMatch: true,
		},
		{
			name:          "Multi-line whole chapter",
			input:         "Romans 3 (ESV) Then what advantage has the Jew? \n2 Or what is the value of circumcision? \n3 Much in every way...",
			expectedRef:   "Romans 3",
			expectedVer:   "ESV",
			expectedText:  "Then what advantage has the Jew? \n2 Or what is the value of circumcision? \n3 Much in every way...",
			expectedMatch: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matches := verseReferenceRegex.FindStringSubmatch(tc.input)
			if tc.expectedMatch {
				assert.Len(t, matches, 4)
				assert.Equal(t, tc.expectedRef, matches[1])
				assert.Equal(t, tc.expectedVer, matches[2])
				assert.Equal(t, tc.expectedText, matches[3])
			} else {
				assert.Empty(t, matches)
			}
		})
	}
}
