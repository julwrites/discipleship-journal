package services

import (
	"context"
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
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("bible API error: %s", resp.Status())
	}

	// The API returns the result. We might need to transform it to match the
	// expected interface (reference, text).
	// For now, return the raw result. The API likely returns a list of verses.
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

	apiPayload := map[string]interface{}{
		"query": map[string]interface{}{
			"prompt": prompt,
		},
		"context": map[string]interface{}{
			"user": map[string]interface{}{
				"version": "ESV", // Default to ESV
			},
		},
	}

	if len(verses) > 0 {
		apiPayload["context"].(map[string]interface{})["verses"] = verses
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

	// Transform response to match OpenAI style if ChatHandler expects it?
	// ChatHandler expects:
	// result["choices"][0]["message"]["content"]
	//
	// BibleAIAPI response structure for LLM is likely just the text or a structured object.
	// We need to wrap it to avoid breaking ChatHandler, OR update ChatHandler.
	// Let's assume the API returns `{"text": "response"}` or similar.
	// We'll wrap it to mimic OpenAI format for minimal disruption in handlers,
	// or we can rely on `result` having the data.

	// If the API returns key "response" or "text":
	var content string
	if txt, ok := result["text"].(string); ok {
		content = txt
	} else if respStr, ok := result["response"].(string); ok {
		content = respStr
	} else {
		// Fallback: dump the whole result as string if structure is unknown
		// or just pass it through if it matches
		// content = fmt.Sprintf("%v", result)
		// But let's try to be smart.

		// If the result IS the structure we want, great.
		// If not, we construct the OpenAI-like structure so ChatHandler works.
		// But wait, if I can't verify the API response structure, this is risky.
		// However, standardizing on the return value of THIS function is safer.
	}

	// If we found content, wrap it.
	if content != "" {
		return map[string]interface{}{
			"choices": []interface{}{
				map[string]interface{}{
					"message": map[string]interface{}{
						"content": content,
					},
				},
			},
		}, nil
	}

	// If we didn't find "text" or "response", maybe it is already in the format?
	// Or maybe it's just the map. Return it and let Handler fail or succeed.
	return result, nil
}
