package services

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
	expectedText := `Then what advantage has the Jew? Or what is the value of circumcision?
2 Much in every way. To begin with, the Jews were entrusted with the oracles of God.
3 What if some were unfaithful? Does their faithlessness nullify the faithfulness of God?
4 By no means! Let God be true though every one were a liar, as it is written, “That you may be justified in your words, and prevail when you are judged.”
5 But if our unrighteousness serves to show the righteousness of God, what shall we say? That God is unrighteous to inflict wrath on us? (I speak in a human way.)
6 By no means! For then how could God judge the world?
7 But if through my lie God's truth abounds to his glory, why am I still being condemned as a sinner?
8 And why not do evil that good may come?—as some people slanderously charge us with saying. Their condemnation is just.
9 What then? Are we Jews any better off? No, not at all. For we have already charged that all, both Jews and Greeks, are under sin,
10 as it is written: “None is righteous, no, not one;
11 no one understands; no one seeks for God.
12 All have turned aside; together they have become worthless; no one does good, not even one.”
13 “Their throat is an open grave; they use their tongues to deceive.” “The venom of asps is under their lips.”
14 “Their mouth is full of curses and bitterness.”
15 “Their feet are swift to shed blood;
16 in their paths are ruin and misery,
17 and the way of peace they have not known.”
18 “There is no fear of God before their eyes.”
19 Now we know that whatever the law says it speaks to those who are under the law, so that every mouth may be stopped, and the whole world may be held accountable to God.
20 For by works of the law no human being will be justified in his sight, since through the law comes knowledge of sin.
21 But now the righteousness of God has been manifested apart from the law, although the Law and the Prophets bear witness to it—
22 the righteousness of God through faith in Jesus Christ for all who believe. For there is no distinction:
23 for all have sinned and fall short of the glory of God,
24 and are justified by his grace as a gift, through the redemption that is in Christ Jesus,
25 whom God put forward as a propitiation by his blood, to be received by faith. This was to show God's righteousness, because in his divine forbearance he had passed over former sins.
26 It was to show his righteousness at the present time, so that he might be just and the justifier of the one who has faith in Jesus.
27 Then what becomes of our boasting? It is excluded. By what kind of law? By a law of works? No, but by the law of faith.
28 For we hold that one is justified by faith apart from works of the law.
29 Or is God the God of Jews only? Is he not the God of Gentiles also? Yes, of Gentiles also,
30 since God is one—who will justify the circumcised by faith and the uncircumcised through faith.
31 Do we then overthrow the law by this faith? By no means! On the contrary, we uphold the law.`

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		responseString := fmt.Sprintf(`{"verse": "Romans 3 (ESV) %s"}`, expectedText)

		// We replace newline with \n for JSON format validity
		responseStringJSON := strings.ReplaceAll(responseString, "\n", "\\n")
		fmt.Fprintln(w, responseStringJSON)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewRealBibleAIClient(server.URL, "key", "")
	ctx := context.Background()

	res, err := client.GetPassage(ctx, "Romans 3", "ESV")
	assert.NoError(t, err)
	assert.Equal(t, "Romans 3", res["reference"])
	assert.Equal(t, "ESV", res["version"])

	// Check if all verses are retained and their newlines are normalized to spaces
	expectedNormalized := strings.ReplaceAll(expectedText, "\n", " ")
	assert.Equal(t, expectedNormalized, res["text"])
}
