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

func TestStreamChatCompletion_Fallback(t *testing.T) {
	tests := []struct {
		name           string
		handler        func(w http.ResponseWriter, r *http.Request)
		expectError    bool
		expectedOutput string
	}{
		{
			name: "Streaming Success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				var req QueryRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}

				// Verify it's a streaming request
				if req.Options == nil || !req.Options.Stream {
					http.Error(w, "expected stream=true", http.StatusBadRequest)
					return
				}

				w.Header().Set("Content-Type", "text/event-stream")
				w.WriteHeader(http.StatusOK)

				// Flush headers
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}

				// V2 Format
				fmt.Fprintf(w, "data: {\"delta\": \"Hello \"}\n\n")
				fmt.Fprintf(w, "data: {\"delta\": \"World\"}\n\n")
				fmt.Fprintf(w, "data: [DONE]\n\n")
			},
			expectError:    false,
			expectedOutput: "Hello World",
		},
		{
			name: "Fallback Triggered (404 on Stream)",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Read body to check if stream is requested
				var req QueryRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}

				if req.Options != nil && req.Options.Stream {
					// Fail the streaming request
					http.Error(w, "Stream not supported", http.StatusNotFound)
					return
				}

				// This must be the fallback request (Stream=false/omitted)
				// Return V2 JSON format
				resp := map[string]interface{}{
					"data": map[string]interface{}{
						"text": "Fallback Response",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(resp)
			},
			expectError:    false,
			expectedOutput: "Fallback Response",
		},
		{
			name: "Fallback Triggered (200 OK but JSON Content-Type)",
			handler: func(w http.ResponseWriter, r *http.Request) {
				var req QueryRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}

				if req.Options != nil && req.Options.Stream {
					// Return JSON instead of stream, simulating API ignoring stream=true
					w.Header().Set("Content-Type", "application/json")
					// Return V2 JSON (or V1 if we were testing mixed, but let's simulate V2 behavior)
					resp := map[string]interface{}{
						"data": map[string]interface{}{
							"text": "Ignored Stream Request",
						},
					}
					_ = json.NewEncoder(w).Encode(resp)
					return
				}

				// Fallback request
				resp := map[string]interface{}{
					"data": map[string]interface{}{
						"text": "Fallback Response from Clean Request",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(resp)
			},
			expectError:    false,
			expectedOutput: "Fallback Response from Clean Request", // Wait, if stream returns valid JSON, performFallback logic might handle it differently?
			// In performFallback: "Streaming failed ... falling back to non-streaming..."
			// But wait, in StreamChatCompletion:
			// "Check if streaming failed (network error or HTTP error) or if response is not event-stream"
			// If 200 OK and not event-stream, it falls back.
			// So it calls performFallback which calls ChatCompletion.
			// So the output comes from the SECOND request (fallback request).
			// My mock handler handles both requests.
			// 1st request (Stream=true) returns JSON. client falls back.
			// 2nd request (Stream=false) returns JSON.
			// So output should be "Fallback Response from Clean Request".
		},
		{
			name: "Total Failure",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Service Unavailable", http.StatusInternalServerError)
			},
			expectError:    true,
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer server.Close()

			client := NewRealBibleAIClient(server.URL, "test-key", "")

			// Use a timeout to prevent hanging tests
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			payload := map[string]interface{}{
				"prompt": "Test Prompt",
			}

			outChan, errChan, err := client.StreamChatCompletion(ctx, payload)
			if err != nil {
				t.Fatalf("StreamChatCompletion setup error: %v", err)
			}

			var result string
			var receivedErr error

			// Collect results
			done := make(chan bool)
			go func() {
				for {
					select {
					case chunk, ok := <-outChan:
						if !ok {
							outChan = nil
						} else {
							result += chunk
						}
					case err, ok := <-errChan:
						if !ok {
							errChan = nil
						} else {
							receivedErr = err
						}
					case <-ctx.Done():
						return
					}
					// Wait for both channels to be closed
					if outChan == nil && errChan == nil {
						close(done)
						return
					}
				}
			}()

			select {
			case <-done:
				// Finished processing
			case <-ctx.Done():
				t.Fatal("Test timed out")
			}

			if tt.expectError {
				if receivedErr == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if receivedErr != nil {
					t.Errorf("Unexpected error: %v", receivedErr)
				}
				if result != tt.expectedOutput {
					t.Errorf("Expected output %q, got %q", tt.expectedOutput, result)
				}
			}
		})
	}
}
