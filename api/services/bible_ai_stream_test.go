package services

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBibleAIClient_Stream_Formats(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		// 1. Delta format (Standard)
		fmt.Fprintln(w, `data: {"delta": "Delta "}`)
		// 2. Text format (Old)
		fmt.Fprintln(w, `data: {"text": "Text "}`)
		// 3. Choices format (OpenAI)
		fmt.Fprintln(w, `data: {"choices": [{"delta": {"content": "Choice"}}]}`)
		// 4. Invalid JSON (Should be treated as raw string)
		fmt.Fprintln(w, `data: RawString`)

		fmt.Fprintln(w, `data: [DONE]`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	ch, errChan, err := client.StreamChatCompletion(ctx, map[string]interface{}{"prompt": "test"})
	assert.NoError(t, err)

	var chunks []string
	done := false
	for !done {
		select {
		case chunk, ok := <-ch:
			if !ok {
				done = true
			} else {
				chunks = append(chunks, chunk)
			}
		case err, ok := <-errChan:
			if !ok {
				continue
			}
			t.Errorf("Unexpected error: %v", err)
			done = true
		}
	}

	expected := []string{"Delta ", "Text ", "Choice", "RawString"}
	assert.Equal(t, expected, chunks)
}
