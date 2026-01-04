package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// BibleAIClient defines the interface for interacting with the Bible AI API.
type BibleAIClient interface {
	GetPassage(ctx context.Context, reference string) (map[string]interface{}, error)
	ChatCompletion(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error)
}

// RealBibleAIClient is the production implementation using the real API.
type RealBibleAIClient struct {
	APIURL string
	APIKey string
	Client *resty.Client
}

// NewRealBibleAIClient creates a new instance of RealBibleAIClient.
func NewRealBibleAIClient(apiURL, apiKey string) *RealBibleAIClient {
	return &RealBibleAIClient{
		APIURL: apiURL,
		APIKey: apiKey,
		Client: resty.New(),
	}
}

// Structures based on OpenAPI Spec

type QueryRequest struct {
	Query   QueryPayload  `json:"query"`
	Context *QueryContext `json:"context,omitempty"`
}

type QueryPayload struct {
	Verses []string `json:"verses,omitempty"`
	Words  []string `json:"words,omitempty"`
	Prompt string   `json:"prompt,omitempty"`
}

type QueryContext struct {
	History []string     `json:"history,omitempty"`
	Schema  string       `json:"schema,omitempty"`
	Verses  []string     `json:"verses,omitempty"`
	Words   []string     `json:"words,omitempty"`
	User    *UserContext `json:"user,omitempty"`
}

type UserContext struct {
	Version string `json:"version,omitempty"`
}

type VerseResponse struct {
	Verse string `json:"verse"`
}

type Reference struct {
	Verse string `json:"verse"`
	URL   string `json:"url"`
}

type OQueryResponse struct {
	Text       string      `json:"text"`
	References []Reference `json:"references"`
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// GetPassage fetches a bible passage from the external API.
func (c *RealBibleAIClient) GetPassage(ctx context.Context, reference string) (map[string]interface{}, error) {
	if c.APIURL == "" {
		return nil, fmt.Errorf("bible API not configured")
	}

	reqPayload := QueryRequest{
		Query: QueryPayload{
			Verses: []string{reference},
		},
	}

	var result VerseResponse
	var errorResult ErrorResponse

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("X-API-KEY", c.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(reqPayload).
		SetResult(&result).
		SetError(&errorResult).
		Post(c.APIURL + "/query")

	if err != nil {
		return nil, fmt.Errorf("resty request error: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("bible API error: %s, message: %s", resp.Status(), errorResult.Error.Message)
	}

	// Helper for parsing if resty failed to unmarshal into result automatically
	// (Though SetResult usually handles it, sometimes API returns different structure)
	if result.Verse == "" {
		// Fallback manual check in case it didn't unmarshal
		var raw map[string]interface{}
		_ = json.Unmarshal(resp.Body(), &raw)
		if v, ok := raw["verse"].(string); ok {
			result.Verse = v
		}
	}

	return map[string]interface{}{
		"verse": result.Verse,
	}, nil
}

// ChatCompletion sends a chat completion request to the external API.
func (c *RealBibleAIClient) ChatCompletion(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	if c.APIURL == "" {
		return nil, fmt.Errorf("bible API not configured")
	}

	prompt, ok := payload["prompt"].(string)
	if !ok {
		prompt = ""
	}

	// Prepare Context
	queryContext := &QueryContext{
		User: &UserContext{
			Version: "ESV",
		},
	}

	// Map 'verses' from payload to context.verses
	if verses, ok := payload["verses"].([]string); ok && len(verses) > 0 {
		queryContext.Verses = verses
	}

	// Map 'themes' from payload to context.words
	if themes, ok := payload["themes"].([]string); ok && len(themes) > 0 {
		queryContext.Words = themes
		// Also optionally append to prompt to ensure LLM focuses on them
		prompt += fmt.Sprintf(" Focus on themes: %s.", strings.Join(themes, ", "))
	}

	// Handle 'context' string (generic text) - API doesn't have a field for this, so append to prompt
	if ctxText, ok := payload["context"].(string); ok && ctxText != "" {
		prompt += fmt.Sprintf(" Context: %s.", ctxText)
	}

	reqPayload := QueryRequest{
		Query: QueryPayload{
			Prompt: prompt,
		},
		Context: queryContext,
	}

	var result OQueryResponse
	var errorResult ErrorResponse

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("X-API-KEY", c.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(reqPayload).
		SetResult(&result).
		SetError(&errorResult).
		Post(c.APIURL + "/query")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("bible API error: %s, message: %s", resp.Status(), errorResult.Error.Message)
	}

	// Construct return map
	response := map[string]interface{}{
		"text":       result.Text,
		"references": result.References,
	}

	// Backward compatibility for ChatHandler which expects OpenAI format
	if result.Text != "" {
		response["choices"] = []interface{}{
			map[string]interface{}{
				"message": map[string]interface{}{
					"content": result.Text,
				},
			},
		}
	}

	return response, nil
}
