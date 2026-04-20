package services

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRealBibleAIClient_SystemPrompts(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		wantKey   string
		wantVal   string
	}{
		{
			name:      "Valid JSON",
			jsonInput: `{"ask": "test"}`,
			wantKey:   "ask",
			wantVal:   "test",
		},
		{
			name: "Invalid JSON with newline in string (Fixable)",
			jsonInput: `{"ask": "line 1
line 2"}`,
			wantKey: "ask",
			wantVal: "line 1\nline 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewRealBibleAIClient("http://localhost", "key", tt.jsonInput)
			if client == nil {
				t.Fatal("NewRealBibleAIClient returned nil")
			}

			val, ok := client.SystemPrompts[tt.wantKey]
			if !ok {
				t.Errorf("SystemPrompts[%q] not found", tt.wantKey)
			} else if val != tt.wantVal {
				t.Errorf("SystemPrompts[%q] = %q, want %q", tt.wantKey, val, tt.wantVal)
			}
		})
	}
}

func TestCleanHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Empty paragraphs",
			input:    "<p>Line 1</p><p>&nbsp;</p><p>Line 2</p>",
			expected: "<p>Line 1</p><p>Line 2</p>",
		},
		{
			name:     "Empty list items_SKIP",
			input:    "<ul><li>Item 1</li><li></li><li>Item 2</li><li>&nbsp;</li></ul>",
			expected: "<ul><li>Item 1</li><li>Item 2</li></ul>",
		},
		{
			name:     "Empty list items with newlines_SKIP",
			input:    "<ul><li>Item 1</li>\n<li>   </li>\n<li>Item 2</li></ul>",
			expected: "<ul><li>Item 1</li><li>Item 2</li></ul>",
		},
		{
			name:     "Mixed empty content",
			input:    "<p>Start</p><p><br></p><ul><li></li><li>Valid</li></ul>",
			expected: "<p>Start</p><ul><li>Valid</li></ul>",
		},
		{
			name:     "Empty list items with br",
			input:    "<ul><li><br></li><li>Item</li></ul>",
			expected: "<ul><li>Item</li></ul>",
		},
		{
			name:     "Stray br before block elements",
			input:    "</p> <br/><p><span><sup>5 </sup>But if our unrighteousness</span></p>",
			expected: "</p><p><span><sup>5 </sup>But if our unrighteousness</span></p>",
		},
		{
			name:     "Loose br before block elements",
			input:    "Some text <br/><p>New Paragraph</p>",
			expected: "Some text <p>New Paragraph</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanHTML(tt.input)
			if got != tt.expected {
				t.Errorf("cleanHTML() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestBibleAIClient_GetVersions(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bible-versions" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"versions": ["ESV", "NIV"]}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	res, err := client.GetVersions(ctx, nil)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	versions, ok := res["versions"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, versions, 2)
}

func TestBibleAIClient_GetVersions_Error(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	res, err := client.GetVersions(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestBibleAIClient_GetPassage_Fallback(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Missing "verse", has "text"
		fmt.Fprintln(w, `{"text": "John 3:16 (ESV) For God so loved..."}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	res, err := client.GetPassage(ctx, "John 3:16", "ESV")
	assert.NoError(t, err)
	assert.Equal(t, "John 3:16", res["reference"])
	assert.Contains(t, res["text"], "For God so loved")
}

func TestBibleAIClient_Stream_Fallback(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// If stream requested, fail or return JSON (not SSE)
		// Resty checks Content-Type for "text/event-stream"
		// We return application/json to trigger fallback
		w.Header().Set("Content-Type", "application/json")
		// Return a normal V2 response which ChatCompletion expects
		fmt.Fprintln(w, `{"data": {"text": "Fallback Response", "references": []}, "meta": {}}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	// Stream should fall back to ChatCompletion
	ch, provider, err := client.Stream(ctx, "prompt")
	assert.NoError(t, err)
	assert.Equal(t, "bible-ai", provider)

	var chunks []string
	for chunk := range ch {
		chunks = append(chunks, chunk)
	}

	assert.Equal(t, []string{"Fallback Response"}, chunks)
}

func TestBibleAIClient_GetSystemPrompt(t *testing.T) {
	prompts := `{"ask": "You are a helpful assistant."}`
	client := NewRealBibleAIClient("url", "key", prompts)

	p := client.GetSystemPrompt("ask")
	assert.Equal(t, "You are a helpful assistant.", p)

	p2 := client.GetSystemPrompt("unknown")
	assert.Equal(t, "", p2)
}

func TestBibleAIClient_Name(t *testing.T) {
	client := NewRealBibleAIClient("url", "key", "")
	assert.Equal(t, "bible-ai", client.Name())
}

func TestBibleAIClient_Query(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/query" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		// Expect V2 response format
		fmt.Fprintln(w, `{"data": {"text": "Response text", "references": []}, "meta": {}}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	text, provider, err := client.Query(ctx, "prompt", "schema")
	assert.NoError(t, err)
	assert.Equal(t, "Response text", text)
	assert.Equal(t, "bible-ai", provider)
}

func TestBibleAIClient_Stream(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		// Send SSE data
		fmt.Fprintln(w, `data: {"delta": "Hello"}`)
		fmt.Fprintln(w, `data: {"delta": " World"}`)
		fmt.Fprintln(w, `data: [DONE]`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	ch, provider, err := client.Stream(ctx, "prompt")
	assert.NoError(t, err)
	assert.Equal(t, "bible-ai", provider)

	var chunks []string
	for chunk := range ch {
		chunks = append(chunks, chunk)
	}

	assert.Equal(t, []string{"Hello", " World"}, chunks)
}

func TestBibleAIClient_GetPassage(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"verse": "John 3:16 (ESV) For God so loved the world..."}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	res, err := client.GetPassage(ctx, "John 3:16", "ESV")
	assert.NoError(t, err)
	assert.Equal(t, "John 3:16", res["reference"])
	assert.Equal(t, "ESV", res["version"])
	assert.Contains(t, res["text"], "For God so loved")
}

func TestBibleAIClient_GetPassage_Error(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	res, err := client.GetPassage(ctx, "John 3:16", "ESV")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestBibleAIClient_ChatCompletion_Error(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	payload := map[string]interface{}{"prompt": "hello"}
	res, err := client.ChatCompletion(ctx, payload)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestBibleAIClient_EmptyURL(t *testing.T) {
	client := NewRealBibleAIClient("", "key", "")
	ctx := context.Background()

	_, _, err := client.Query(ctx, "prompt", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bible API not configured")

	_, _, err = client.Stream(ctx, "prompt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bible API not configured")

	_, err = client.GetPassage(ctx, "ref", "ver")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bible API not configured")

	_, err = client.ChatCompletion(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bible API not configured")

	_, _, err = client.StreamChatCompletion(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bible API not configured")

	_, err = client.GetVersions(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bible API not configured")
}

func TestBibleAIClient_Stream_Error(t *testing.T) {
	// Simulate connection refused to trigger immediate error in Stream
	client := NewRealBibleAIClient("http://localhost:12345", "key", "") // Invalid port
	ctx := context.Background()

	// This should fail immediately because Stream calls StreamChatCompletion which calls resty.Post which fails
	ch, _, err := client.Stream(ctx, "prompt")
	// Stream doesn't return error immediately if the goroutine handles it, BUT StreamChatCompletion returns channels.
	// Oh, Stream checks `err` from StreamChatCompletion.
	// `StreamChatCompletion` returns nil error because it spawns goroutine.
	// Wait, my `StreamChatCompletion` implementation:
	// `go func() { ... }`
	// returns `outChan, errChan, nil`.
	// So `Stream` gets nil error.
	// Then `Stream` returns `safeOutChan, nil`.
	// So `err` is nil here.
	assert.NoError(t, err)

	// Consume channel to get error?
	// The `Stream` implementation swallows error from errChan or logs it?
	// `Stream` implementation:
	// case err, ok := <-errChan: slog.Error(...) return
	// It basically stops streaming.
	// So `ch` should close.
	count := 0
	for range ch {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestBibleAIClient_Stream_Error_Propagation(t *testing.T) {
	// This test relies on Stream implementation detail where it reads from errChan.
	// Since we can't easily inspect logs, we just verify channel closes without data.
	client := NewRealBibleAIClient("http://localhost:12345", "key", "")
	ch, _, err := client.Stream(context.Background(), "prompt")
	assert.NoError(t, err)
	for range ch {
		t.Fail() // Should produce no data
	}
}

func TestBibleAIClient_Query_Error(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	_, _, err := client.Query(ctx, "prompt", "")
	assert.Error(t, err)
}

func TestNewRealBibleAIClient_InvalidJSON(t *testing.T) {
	// Should log warning but not panic
	client := NewRealBibleAIClient("url", "key", "{invalid-json}")
	assert.NotNil(t, client)
	assert.Empty(t, client.SystemPrompts)
}

func TestBibleAIClient_StreamChatCompletion_EOF(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		// Close immediately
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	out, errs, err := client.StreamChatCompletion(context.Background(), map[string]interface{}{})
	assert.NoError(t, err)

	// Drain
	for range out {
	}
	for range errs {
	}
}

func TestBibleAIClient_StreamChatCompletion_ParseError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintln(w, "data: {invalid-json}")
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	out, _, _ := client.StreamChatCompletion(context.Background(), map[string]interface{}{})

	// Should receive raw string
	str := <-out
	assert.Equal(t, "{invalid-json}", str)
}

func TestBibleAIClient_StreamChatCompletion_V2Delta(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintln(w, `data: {"choices": [{"delta": {"content": "Hello"}}]}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	out, _, _ := client.StreamChatCompletion(context.Background(), map[string]interface{}{})

	str := <-out
	assert.Equal(t, "Hello", str)
}

func TestBibleAIClient_StreamChatCompletion_FallbackText(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintln(w, `data: {"text": "Hello"}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	out, _, _ := client.StreamChatCompletion(context.Background(), map[string]interface{}{})

	str := <-out
	assert.Equal(t, "Hello", str)
}

func TestBibleAIClient_StreamChatCompletion_ImmediateError(t *testing.T) {
	// Connection error
	client := NewRealBibleAIClient("http://invalid-host", "key", "")

	// Since it's async, it returns nil err, but pushes to errChan?
	// Or performFallback is called?
	// If error occurs during request (Post), it calls performFallback.

	out, errs, err := client.StreamChatCompletion(context.Background(), map[string]interface{}{})
	assert.NoError(t, err)

	// Wait for fallback failure (since chat completion also fails)
	select {
	case err := <-errs:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "fallback failed")
	case <-out:
		// might close
	}
}

func TestBibleAIClient_GetPassage_Romans3(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"verse": "Romans 3 (ESV) Then what advantage has the Jew? Or what is the value of circumcision?\n2 Much in every way. To begin with, the Jews were entrusted with the oracles of God.\n3 What if some were unfaithful? Does their faithlessness nullify the faithfulness of God?"}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	res, err := client.GetPassage(ctx, "Romans 3", "ESV")
	assert.NoError(t, err)
	assert.Equal(t, "Romans 3", res["reference"])
	assert.Equal(t, "ESV", res["version"])

	// Check if all verses are retained with their newlines intact
	expectedText := "Then what advantage has the Jew? Or what is the value of circumcision?\n2 Much in every way. To begin with, the Jews were entrusted with the oracles of God.\n3 What if some were unfaithful? Does their faithlessness nullify the faithfulness of God?"
	assert.Equal(t, expectedText, res["text"])
}
