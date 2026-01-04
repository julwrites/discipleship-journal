package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"

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

// parsePassageFromHTML converts HTML from Bible API to markdown format.
// It handles common tags found in Bible passages: p, span, b, i, sup, br, h1-h4.
func parsePassageFromHTML(htmlStr string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	var result strings.Builder
	var parseNode func(*html.Node)

	parseNode = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			result.WriteString(n.Data)
		case html.ElementNode:
			switch n.Data {
			case "b", "strong":
				result.WriteString("**")
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					parseNode(c)
				}
				result.WriteString("**")
			case "i", "em":
				result.WriteString("*")
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					parseNode(c)
				}
				result.WriteString("*")
			case "sup":
				result.WriteString("^")
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					parseNode(c)
				}
				result.WriteString("^")
			case "br":
				result.WriteString("\n")
			case "p":
				// Add newline before paragraph (except first)
				if result.Len() > 0 {
					result.WriteString("\n\n")
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					parseNode(c)
				}
			case "h1", "h2", "h3", "h4":
				result.WriteString("\n\n**")
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					parseNode(c)
				}
				result.WriteString("**\n")
			case "span":
				// Span just passes through
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					parseNode(c)
				}
			default:
				// Unknown element, process children
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					parseNode(c)
				}
			}
		}
	}

	// Find body or start from root
	var start *html.Node
	if doc.Type == html.DocumentNode && doc.FirstChild != nil && doc.FirstChild.NextSibling != nil {
		// Skip doctype and html, go to body
		htmlNode := doc.FirstChild.NextSibling
		if htmlNode.FirstChild != nil {
			start = htmlNode.FirstChild.NextSibling // body
		}
	}
	if start == nil {
		start = doc
	}

	parseNode(start)

	// Clean up extra whitespace
	text := result.String()
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "  ", " ")
	text = strings.ReplaceAll(text, "\n\n\n", "\n\n")

	return text, nil
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

	// Convert HTML to markdown for display and insertion
	verseText := result.Verse
	if verseText != "" {
		markdown, err := parsePassageFromHTML(verseText)
		if err != nil {
			log.Printf("Failed to convert HTML to markdown: %v. Using original text.", err)
		} else {
			verseText = markdown
		}
	}

	return map[string]interface{}{
		"verse": verseText,
		"text":  verseText,
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
