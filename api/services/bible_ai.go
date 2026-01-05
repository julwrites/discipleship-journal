package services

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"regexp"
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
	APIURL        string
	APIKey        string
	Client        *resty.Client
	SystemPrompts map[string]string
}

// NewRealBibleAIClient creates a new instance of RealBibleAIClient.
func NewRealBibleAIClient(apiURL, apiKey, systemPromptsJSON string) *RealBibleAIClient {
	prompts := make(map[string]string)
	if systemPromptsJSON != "" {
		if err := json.Unmarshal([]byte(systemPromptsJSON), &prompts); err != nil {
			log.Printf("Failed to parse system prompts JSON: %v", err)
		}
	}

	return &RealBibleAIClient{
		APIURL:        apiURL,
		APIKey:        apiKey,
		Client:        resty.New(),
		SystemPrompts: prompts,
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
	Verse     string `json:"verse"`
	Reference string `json:"reference,omitempty"`
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
		log.Printf("Bible API response status %s, body length %d", resp.Status(), len(resp.Body()))
		// Fallback manual check in case it didn't unmarshal
		var raw map[string]interface{}
		_ = json.Unmarshal(resp.Body(), &raw)
		if v, ok := raw["verse"].(string); ok {
			result.Verse = v
		} else if t, ok := raw["text"].(string); ok {
			result.Verse = t
		} else {
			log.Printf("Bible API response body: %s", resp.Body())
		}
	}

	// Return raw HTML from the API, unescaping it in case it was escaped in the JSON response
	verseText := html.UnescapeString(result.Verse)

	return map[string]interface{}{
		"verse":     verseText,
		"text":      verseText,
		"reference": result.Reference,
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

	// Handle optional context by appending it to the prompt variable
	// This ensures it is included even if the template doesn't explicitly use {CONTEXT}
	if ctxText, ok := payload["context"].(string); ok && ctxText != "" {
		prompt += fmt.Sprintf(" Context: %s.", ctxText)
	}

	// Prepare Context
	queryContext := &QueryContext{
		User: &UserContext{
			Version: "ESV",
		},
	}

	var themes []string
	if ts, ok := payload["themes"].([]string); ok {
		themes = ts
	}
	if len(themes) > 0 {
		queryContext.Words = themes
	}

	var verses []string
	if vs, ok := payload["verses"].([]string); ok {
		verses = vs
	}
	if len(verses) > 0 {
		queryContext.Verses = verses
	}

	// Select prompt template
	promptType, _ := payload["type"].(string)
	if promptType == "" {
		promptType = "ask"
	}

	template := c.SystemPrompts[promptType]

	finalPrompt := prompt
	if template != "" {
		// Replace tags: {PROMPT}, {WORDS}, {PASSAGE}
		r := strings.NewReplacer(
			"{PROMPT}", prompt,
			"{WORDS}", strings.Join(themes, ", "),
			"{PASSAGE}", strings.Join(verses, "\n"),
		)
		finalPrompt = r.Replace(template)
	} else {
		// Fallback/Legacy behavior: append context manually if prompt template not found
		// (though we already appended context to prompt above, so just add themes)
		if len(themes) > 0 {
			finalPrompt += fmt.Sprintf(" Focus on themes: %s.", strings.Join(themes, ", "))
		}
	}

	reqPayload := QueryRequest{
		Query: QueryPayload{
			Prompt: finalPrompt,
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

	// Unescape the text before cleaning
	unescapedText := html.UnescapeString(result.Text)
	cleanedText := cleanHTML(unescapedText)

	// Construct return map
	response := map[string]interface{}{
		"text":       cleanedText,
		"references": result.References,
	}

	// Backward compatibility for ChatHandler which expects OpenAI format
	if cleanedText != "" {
		response["choices"] = []interface{}{
			map[string]interface{}{
				"message": map[string]interface{}{
					"content": cleanedText,
				},
			},
		}
	}

	return response, nil
}

var (
	emptyParaRegex = regexp.MustCompile(`(?i)<p[^>]*>(\s|&nbsp;|<br\s*/?>)*</p>`)
	newlineRegex   = regexp.MustCompile(`[\r\n]+`)
)

func cleanHTML(input string) string {
	// Replace newlines with space to avoid word concatenation if used in text
	s := newlineRegex.ReplaceAllString(input, " ")

	// Remove empty paragraphs
	s = emptyParaRegex.ReplaceAllString(s, "")

	return strings.TrimSpace(s)
}
