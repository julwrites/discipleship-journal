package services

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("Bible API not configured")
	}

	var result map[string]interface{}
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+c.APIKey).
		SetQueryParam("q", reference).
		SetResult(&result).
		Get(c.APIURL + "/bible/passage")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("Bible API Error: %s", resp.Status())
	}

	return result, nil
}

// ChatCompletion sends a chat completion request to the external API.
func (c *RealBibleAIClient) ChatCompletion(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	if c.APIURL == "" {
		return nil, fmt.Errorf("Bible API not configured")
	}

	var result map[string]interface{}
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+c.APIKey).
		SetBody(payload).
		SetResult(&result).
		Post(c.APIURL + "/chat/completions")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("Bible API Error: %s", resp.Status())
	}

	return result, nil
}
