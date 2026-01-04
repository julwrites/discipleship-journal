package services

import (
	"context"
	"fmt"
)

// MockBibleAIClient is a mock implementation for testing.
type MockBibleAIClient struct {
	ShouldError bool
}

// NewMockBibleAIClient creates a new mock client.
func NewMockBibleAIClient() *MockBibleAIClient {
	return &MockBibleAIClient{}
}

func (m *MockBibleAIClient) GetPassage(ctx context.Context, reference string) (map[string]interface{}, error) {
	if m.ShouldError {
		return nil, fmt.Errorf("mock error")
	}
	verseText := "For God so loved the world... (Mocked)"
	return map[string]interface{}{
		"reference": reference,
		"text":      verseText,
		"verse":     verseText,
	}, nil
}

func (m *MockBibleAIClient) ChatCompletion(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	if m.ShouldError {
		return nil, fmt.Errorf("mock error")
	}

	// Extract prompt from payload for dynamic mock response if needed
	// New payload structure: "prompt", "verses", "themes", "context"
	var prompt string
	if p, ok := payload["prompt"].(string); ok {
		prompt = p
	}

	return map[string]interface{}{
		"choices": []interface{}{
			map[string]interface{}{
				"message": map[string]interface{}{
					"content": "This is a mocked AI response to: " + prompt,
				},
			},
		},
	}, nil
}
