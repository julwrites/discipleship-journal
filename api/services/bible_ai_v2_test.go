package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetPassage_V2Format(t *testing.T) {
	// Simulate V2 API response where 'verse' contains the reference
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"verse": "John 3:16 (ESV) For God so loved the world, that he gave his only Son.",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "test-key", "")
	ctx := context.Background()

	result, err := client.GetPassage(ctx, "John 3:16", "ESV")
	if err != nil {
		t.Fatalf("GetPassage failed: %v", err)
	}

	// Expect parsed fields
	expectedRef := "John 3:16"
	expectedText := "For God so loved the world, that he gave his only Son."

	if ref, ok := result["reference"].(string); !ok || ref != expectedRef {
		t.Errorf("Expected reference %q, got %q", expectedRef, ref)
	}

	if text, ok := result["text"].(string); !ok || text != expectedText {
		t.Errorf("Expected text %q, got %q", expectedText, text)
	}
}

func TestStreamChatCompletion_RawText(t *testing.T) {
	// Simulate SSE with raw text data instead of JSON
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// Send raw text lines
		fmt.Fprintf(w, "data: Hello\n\n")
		fmt.Fprintf(w, "data:  World\n\n") // Space after colon is trimmed? "data: " prefix removal
		fmt.Fprintf(w, "data: [DONE]\n\n")
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "test-key", "")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	payload := map[string]interface{}{"prompt": "hi"}
	outChan, errChan, err := client.StreamChatCompletion(ctx, payload)
	if err != nil {
		t.Fatalf("StreamChatCompletion failed: %v", err)
	}

	var result string
	done := make(chan bool)
	go func() {
		for msg := range outChan {
			result += msg
		}
		close(done)
	}()

	select {
	case <-done:
	case err, ok := <-errChan:
		if ok {
			t.Fatalf("Stream error: %v", err)
		}
		// If closed, wait for done
		<-done
	case <-ctx.Done():
		t.Fatal("Timeout")
	}

	expected := "Hello World"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
