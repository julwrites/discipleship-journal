package services

import (
	"context"
	"fmt"
	"sync"
)

type contextKey string

const (
	AIProviderKey   contextKey = "ai_provider"
	BibleVersionKey contextKey = "bible_version"
)

// LLMClient defines the interface for interacting with an LLM provider.
type LLMClient interface {
	// Query performs a blocking request.
	// Returns the raw response string, the provider name, and error.
	Query(ctx context.Context, prompt string, schema string) (response string, providerName string, err error)

	// Stream performs a streaming request.
	// Returns a channel for chunks, the provider name, and error (immediate).
	// Note: The channel closes when the stream is done.
	Stream(ctx context.Context, prompt string) (<-chan string, string, error)

	// Name returns the identifier of the provider (e.g., "openai").
	Name() string
}

// FallbackClient manages multiple LLM clients and attempts them in order.
type FallbackClient struct {
	clients           map[string]LLMClient
	priorityOrder     []string
	preferredProvider string
	mu                sync.RWMutex
}

// NewFallbackClient creates a new FallbackClient.
func NewFallbackClient(clients []LLMClient, priorityOrder []string) *FallbackClient {
	clientMap := make(map[string]LLMClient)
	for _, c := range clients {
		clientMap[c.Name()] = c
	}

	return &FallbackClient{
		clients:       clientMap,
		priorityOrder: priorityOrder,
	}
}

// SetPreferredProvider sets the preferred provider for subsequent requests.
func (fc *FallbackClient) SetPreferredProvider(providerName string) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.preferredProvider = providerName
}

// Query attempts to query using the preferred provider, then falls back to priority order.
func (fc *FallbackClient) Query(ctx context.Context, prompt string, schema string) (string, string, error) {
	fc.mu.RLock()
	preferred := fc.preferredProvider
	fc.mu.RUnlock()

	if p, ok := ctx.Value(AIProviderKey).(string); ok && p != "" {
		preferred = p
	}

	var order []string
	if preferred != "" {
		order = append(order, preferred)
	}
	// Append others, avoiding duplicates if preferred is in priorityOrder
	seen := make(map[string]bool)
	if preferred != "" {
		seen[preferred] = true
	}
	for _, name := range fc.priorityOrder {
		if !seen[name] {
			order = append(order, name)
		}
	}

	var lastErr error

	for _, name := range order {
		client, ok := fc.clients[name]
		if !ok {
			continue
		}

		resp, providerName, err := client.Query(ctx, prompt, schema)
		if err == nil {
			return resp, providerName, nil
		}
		lastErr = err
	}

	return "", "", fmt.Errorf("all providers failed, last error: %v", lastErr)
}

// Stream attempts to stream using the preferred provider, then falls back to priority order.
func (fc *FallbackClient) Stream(ctx context.Context, prompt string) (<-chan string, string, error) {
	fc.mu.RLock()
	preferred := fc.preferredProvider
	fc.mu.RUnlock()

	if p, ok := ctx.Value(AIProviderKey).(string); ok && p != "" {
		preferred = p
	}

	var order []string
	if preferred != "" {
		order = append(order, preferred)
	}
	seen := make(map[string]bool)
	if preferred != "" {
		seen[preferred] = true
	}
	for _, name := range fc.priorityOrder {
		if !seen[name] {
			order = append(order, name)
		}
	}

	var lastErr error

	for _, name := range order {
		client, ok := fc.clients[name]
		if !ok {
			continue
		}

		outChan, providerName, err := client.Stream(ctx, prompt)
		if err == nil {
			return outChan, providerName, nil
		}
		lastErr = err
	}

	return nil, "", fmt.Errorf("all providers failed, last error: %v", lastErr)
}

// Name returns the name of the fallback client.
func (fc *FallbackClient) Name() string {
	return "fallback"
}
