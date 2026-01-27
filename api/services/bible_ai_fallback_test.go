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
				if !req.Query.Stream {
					http.Error(w, "expected stream=true", http.StatusBadRequest)
					return
				}

				w.Header().Set("Content-Type", "text/event-stream")
				w.WriteHeader(http.StatusOK)

				// Flush headers
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}

				fmt.Fprintf(w, "data: {\"text\": \"Hello \"}\n\n")
				fmt.Fprintf(w, "data: {\"text\": \"World\"}\n\n")
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

				if req.Query.Stream {
					// Fail the streaming request
					http.Error(w, "Stream not supported", http.StatusNotFound)
					return
				}

				// This must be the fallback request (Stream=false/omitted)
				resp := OQueryResponse{
					Text: "Fallback Response",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
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

				if req.Query.Stream {
					// Return JSON instead of stream, simulating API ignoring stream=true
					w.Header().Set("Content-Type", "application/json")
					// We return some dummy JSON that might be returned if stream was ignored
					resp := OQueryResponse{
						Text: "Ignored Stream Request",
					}
					json.NewEncoder(w).Encode(resp)
					return
				}

				// Fallback request
				resp := OQueryResponse{
					Text: "Fallback Response from Clean Request",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			},
			expectError:    false,
			expectedOutput: "Fallback Response from Clean Request",
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
