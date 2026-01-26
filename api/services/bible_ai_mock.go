package services

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MockBibleAIClient is a mock implementation for testing.
type MockBibleAIClient struct {
	ShouldError bool
}

// NewMockBibleAIClient creates a new mock client.
func NewMockBibleAIClient() *MockBibleAIClient {
	return &MockBibleAIClient{}
}

func (m *MockBibleAIClient) GetPassage(ctx context.Context, reference string, version string) (map[string]interface{}, error) {
	if m.ShouldError {
		return nil, fmt.Errorf("mock error")
	}
	verseText := "For God so loved the world... (Mocked)"
	return map[string]interface{}{
		"reference": reference,
		"text":      verseText,
		"verse":     verseText,
		"version":   version,
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

func (m *MockBibleAIClient) StreamChatCompletion(ctx context.Context, payload map[string]interface{}) (<-chan string, <-chan error, error) {
	if m.ShouldError {
		return nil, nil, fmt.Errorf("mock error")
	}

	outChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(outChan)
		defer close(errChan)

		// Mock response content
		var prompt string
		if p, ok := payload["prompt"].(string); ok {
			prompt = p
		}
		fullResponse := "This is a mocked AI response to: " + prompt

		// Simulate chunks
		words := strings.Split(fullResponse, " ")
		for _, word := range words {
			select {
			case <-ctx.Done():
				return
			case outChan <- word + " ":
				time.Sleep(10 * time.Millisecond) // Simulate delay
			}
		}
	}()

	return outChan, errChan, nil
}

func (m *MockBibleAIClient) GetSystemPrompt(key string) string {
	if key == "system" {
		return "Mock System Prompt"
	}
	return ""
}

func (m *MockBibleAIClient) GetVersions(ctx context.Context, params map[string]string) (map[string]interface{}, error) {
	if m.ShouldError {
		return nil, fmt.Errorf("mock error")
	}
	return map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"id":           "ESV",
				"name":         "English Standard Version",
				"language":     "English",
				"abbreviation": "ESV",
			},
			map[string]interface{}{
				"id":           "NIV",
				"name":         "New International Version",
				"language":     "English",
				"abbreviation": "NIV",
			},
		},
		"meta": map[string]interface{}{
			"total": 2,
			"page":  1,
			"limit": 20,
		},
	}, nil
}
