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
	return map[string]interface{}{
		"reference": reference,
		"text":      "For God so loved the world... (Mocked)",
	}, nil
}

func (m *MockBibleAIClient) ChatCompletion(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	if m.ShouldError {
		return nil, fmt.Errorf("mock error")
	}

	// Extract prompt from payload for dynamic mock response if needed
	var prompt string
	if msgs, ok := payload["messages"].([]map[string]string); ok && len(msgs) > 0 {
		prompt = msgs[len(msgs)-1]["content"]
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
