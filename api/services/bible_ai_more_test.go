package services

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBibleAI_MockClient_StreamChatCompletion(t *testing.T) {
	mockClient := NewMockBibleAIClient()

	// Test the other branches of GetSystemPrompt
	assert.Equal(t, "Mock System Prompt", mockClient.GetSystemPrompt("system"))
	assert.Equal(t, "", mockClient.GetSystemPrompt("unknown"))

	outChan, errChan, err := mockClient.StreamChatCompletion(context.Background(), map[string]interface{}{
		"messages": []interface{}{
            map[string]interface{}{"role": "user", "content": "test"},
        },
	})
	assert.NoError(t, err)

	// Read from channel
	for msg := range outChan {
		assert.NotNil(t, msg)
	}
    for err := range errChan {
		assert.Nil(t, err)
	}
}

func TestBibleAI_MockClient_GetVersions(t *testing.T) {
    mockClient := NewMockBibleAIClient()
    versions, err := mockClient.GetVersions(context.Background(), nil)
    assert.NoError(t, err)
    assert.NotNil(t, versions)
}
