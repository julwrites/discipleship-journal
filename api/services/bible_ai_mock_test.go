package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockBibleAIClient(t *testing.T) {
	client := NewMockBibleAIClient()

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "mock", client.Name())
	})

	t.Run("Query", func(t *testing.T) {
		text, provider, err := client.Query(context.Background(), "prompt", "schema")
		assert.NoError(t, err)
		assert.Equal(t, "This is a mocked AI response to: prompt", text)
		assert.Equal(t, "mock", provider)
	})

	t.Run("Stream", func(t *testing.T) {
		ch, provider, err := client.Stream(context.Background(), "prompt")
		assert.NoError(t, err)
		assert.Equal(t, "mock", provider)

		var content string
		for chunk := range ch {
			content += chunk
		}
		assert.Equal(t, "This is a mocked AI response to: prompt ", content)
	})

	t.Run("GetVersions", func(t *testing.T) {
		res, err := client.GetVersions(context.Background(), nil)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("StreamChatCompletion", func(t *testing.T) {
		ch, errChan, err := client.StreamChatCompletion(context.Background(), map[string]interface{}{"prompt": "test"})
		assert.NoError(t, err)

		var content string
		done := false
		for !done {
			select {
			case chunk, ok := <-ch:
				if !ok {
					done = true
				} else {
					content += chunk
				}
			case err, ok := <-errChan:
				if ok {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		}
		assert.Equal(t, "This is a mocked AI response to: test ", content)
	})
}
