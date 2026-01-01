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

// GetPassage fetches a bible passage from the external API.
func (c *RealBibleAIClient) GetPassage(ctx context.Context, reference string) (map[string]interface{}, error) {
	if c.APIURL == "" {
		return nil, fmt.Errorf("bible API not configured")
	}

	// Construct payload for BibleAIAPI: POST /query with {"query": {"verses": [ref]}}
	payload := map[string]interface{}{
		"query": map[string]interface{}{
			"verses": []string{reference},
		},
	}

	var result map[string]interface{}
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("X-API-KEY", c.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&result).
		Post(c.APIURL + "/query")

	if err != nil {
		return nil, fmt.Errorf("resty request error: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("bible API error: %s, body: %s", resp.Status(), resp.String())
	}

	// resty sometimes fails to parse JSON, so fall back to manual parsing
	if len(result) == 0 {
		body := resp.String()
		var manualResult map[string]interface{}
		if err := json.Unmarshal([]byte(body), &manualResult); err != nil {
			return nil, fmt.Errorf("failed to parse Bible API response: %v", err)
		}
		result = manualResult
	}

	// The API returns {"verse": "John 3:16 (ESV) For God so loved the world..."}
	// Extract the verse text from the response
	if verse, ok := result["verse"].(string); ok {
		return map[string]interface{}{
			"verse": verse,
		}, nil
	}

	// Check for error response
	if errObj, ok := result["error"].(map[string]interface{}); ok {
		return nil, fmt.Errorf("bible API error: %v", errObj)
	}

	// Return raw result if verse not found (for backward compatibility)
	return result, nil
}

// ChatCompletion sends a chat completion request to the external API.
func (c *RealBibleAIClient) ChatCompletion(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	if c.APIURL == "" {
		return nil, fmt.Errorf("bible API not configured")
	}

	// Transform generic payload to BibleAIAPI payload
	// Expected input payload (from ChatHandler):
	// {
	//    "prompt": "User question...",
	//    "verses": ["John 3:16"],
	//    "themes": ["Love"],
	//    "context": "Previous context..." (optional)
	// }
	//
	// Output payload for BibleAIAPI:
	// {
	//   "query": { "prompt": "..." },
	//   "context": { "verses": [...], "user": { "version": "ESV" } }
	// }

	prompt, ok := payload["prompt"].(string)
	if !ok {
		prompt = ""
	}
	verses, ok := payload["verses"].([]string)
	if !ok {
		verses = []string{}
	}

	// If themes are present, append to prompt
	if themes, ok := payload["themes"].([]string); ok && len(themes) > 0 {
		prompt += fmt.Sprintf(" Focus on themes: %s.", strings.Join(themes, ", "))
	}

	// If context/text is present (for general AskAI), append to prompt or handle differently
	// The API supports "words" in context, or we can just shove it in the prompt.
	if ctxText, ok := payload["context"].(string); ok && ctxText != "" {
		prompt += fmt.Sprintf(" Context: %s.", ctxText)
	}

	// Build context object
	contextObj := map[string]interface{}{
		"user": map[string]interface{}{
			"version": "ESV", // Default to ESV
		},
	}

	if len(verses) > 0 {
		contextObj["verses"] = verses
	}

	apiPayload := map[string]interface{}{
		"query": map[string]interface{}{
			"prompt": prompt,
		},
		"context": contextObj,
	}

	var result map[string]interface{}
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("X-API-KEY", c.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(apiPayload).
		SetResult(&result).
		Post(c.APIURL + "/query")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("bible API error: %s", resp.Status())
	}

	// resty sometimes fails to parse JSON, so fall back to manual parsing
	if len(result) == 0 {
		body := resp.String()
		var manualResult map[string]interface{}
		if err := json.Unmarshal([]byte(body), &manualResult); err != nil {
			return nil, fmt.Errorf("failed to parse Bible API response: %v", err)
		}
		result = manualResult
	}

	// BibleAIAPI response structure for LLM queries:
	// {
	//   "text": "Response text...",
	//   "references": [
	//     {"verse": "John 3:16", "url": "..."}
	//   ]
	// }
	//
	// Check for error response first
	if errObj, ok := result["error"].(map[string]interface{}); ok {
		return nil, fmt.Errorf("bible API error: %v", errObj)
	}

	// Extract text from response
	var content string
	if txt, ok := result["text"].(string); ok {
		content = txt
	} else if respStr, ok := result["response"].(string); ok {
		content = respStr
	}

	// Extract references if present
	var references []interface{}
	if refs, ok := result["references"].([]interface{}); ok {
		references = refs
	}

	// Return structured response that handlers can parse
	response := map[string]interface{}{
		"text": content,
	}
	if len(references) > 0 {
		response["references"] = references
	}

	// Also wrap in OpenAI format for backward compatibility with ChatHandler
	if content != "" {
		response["choices"] = []interface{}{
			map[string]interface{}{
				"message": map[string]interface{}{
					"content": content,
				},
			},
		}
	}

	return response, nil
}
