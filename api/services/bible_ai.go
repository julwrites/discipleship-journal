package services

import (
	"bufio"
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
	LLMClient
	GetPassage(ctx context.Context, reference string, version string) (map[string]interface{}, error)
	ChatCompletion(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error)
	StreamChatCompletion(ctx context.Context, payload map[string]interface{}) (<-chan string, <-chan error, error)
	GetVersions(ctx context.Context, params map[string]string) (map[string]interface{}, error)
	GetSystemPrompt(key string) string
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
	Options *QueryOptions `json:"options,omitempty"`
}

type QueryPayload struct {
	Verses []string `json:"verses,omitempty"`
	Words  []string `json:"words,omitempty"`
	Prompt string   `json:"prompt,omitempty"`
}

type QueryOptions struct {
	Stream bool `json:"stream,omitempty"`
}

type QueryContext struct {
	History []string     `json:"history,omitempty"`
	Schema  string       `json:"schema,omitempty"`
	Verses  []string     `json:"verses,omitempty"`
	Words   []string     `json:"words,omitempty"`
	User    *UserContext `json:"user,omitempty"`
}

type UserContext struct {
	Version    string `json:"version,omitempty"`
	AIProvider string `json:"ai_provider,omitempty"`
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

// LLMClient Interface Implementation

// Name returns the identifier of the provider.
func (c *RealBibleAIClient) Name() string {
	return "bible-ai"
}

// Query performs a blocking request using the new API schema.
func (c *RealBibleAIClient) Query(ctx context.Context, prompt string, schema string) (string, string, error) {
	if c.APIURL == "" {
		return "", "", fmt.Errorf("bible API not configured")
	}

	version := "ESV"
	if v, ok := ctx.Value(BibleVersionKey).(string); ok && v != "" {
		version = v
	}

	aiProvider := ""
	if p, ok := ctx.Value(AIProviderKey).(string); ok && p != "" {
		aiProvider = p
	}

	// Construct request
	reqPayload := QueryRequest{
		Query: QueryPayload{
			Prompt: prompt,
		},
		Context: &QueryContext{
			Schema: schema,
			User: &UserContext{
				Version:    version,
				AIProvider: aiProvider,
			},
		},
	}

	var errorResult ErrorResponse

	// If schema is provided, the result structure might be dynamic, but OQueryResponse wraps "data".
	// The prompt says: "When stream is false, the API returns a JSON object wrapping the result and metadata."
	// "data": { "text": "...", ... }
	// But struct OQueryResponse matches "data" fields partially?
	// Actually OQueryResponse struct is: Text string, References []Reference.
	// The new V2 response:
	// { "data": { "text": "...", "references": [...] }, "meta": ... }
	// So we need a wrapper struct for V2 response.

	type V2Response struct {
		Data OQueryResponse         `json:"data"`
		Meta map[string]interface{} `json:"meta"`
	}

	var v2Result V2Response

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("X-API-KEY", c.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(reqPayload).
		SetResult(&v2Result).
		SetError(&errorResult).
		Post(c.APIURL + "/query")

	if err != nil {
		return "", c.Name(), err
	}

	if resp.IsError() {
		return "", c.Name(), fmt.Errorf("bible API error: %s, message: %s", resp.Status(), errorResult.Error.Message)
	}

	// Unescape the text before returning
	unescapedText := html.UnescapeString(v2Result.Data.Text)
	cleanedText := cleanHTML(unescapedText)

	return cleanedText, c.Name(), nil
}

// Stream performs a streaming request using the new API schema.
func (c *RealBibleAIClient) Stream(ctx context.Context, prompt string) (<-chan string, string, error) {
	// Re-use StreamChatCompletion logic but simplified for LLMClient interface
	// Create payload expected by StreamChatCompletion
	payload := map[string]interface{}{
		"prompt": prompt,
	}

	if v, ok := ctx.Value(BibleVersionKey).(string); ok && v != "" {
		payload["version"] = v
	}
	if p, ok := ctx.Value(AIProviderKey).(string); ok && p != "" {
		payload["ai_provider"] = p
	}

	outChan, errChan, err := c.StreamChatCompletion(ctx, payload)
	if err != nil {
		return nil, c.Name(), err
	}

	// Convert error channel to immediate return is not possible since StreamChatCompletion is async.
	// But LLMClient.Stream returns (<-chan string, string, error).
	// We need to bridge the gap. StreamChatCompletion returns (outChan, errChan, error).
	// We can merge errChan into the returned channel or handle it differently.
	// The LLMClient interface says "Returns a channel for chunks... error (immediate)".
	// If StreamChatCompletion returns immediate error, we return it.
	// If it returns channels, we need to handle runtime errors from errChan?
	// The LLMClient interface definition in the prompt says:
	// "Stream... Returns a channel for chunks, the provider name, and error (immediate)."
	// It doesn't mention an error channel. Typically this means errors during stream are either logged or panic or the channel is just closed.
	// Or maybe the channel should be type that supports error?
	// No, it is `<-chan string`.
	// For now, I will start a goroutine to drain errChan and maybe log it, or close outChan on error.

	// A better approach is to wrap the output channel to handle errors if we can't return them.
	// But strictly implementing the interface:

	safeOutChan := make(chan string)
	go func() {
		defer close(safeOutChan)
		for {
			select {
			case msg, ok := <-outChan:
				if !ok {
					return
				}
				safeOutChan <- msg
			case err, ok := <-errChan:
				if ok {
					log.Printf("Stream error from BibleAI: %v", err)
				}
				return // Stop streaming on error
			case <-ctx.Done():
				return
			}
		}
	}()

	return safeOutChan, c.Name(), nil
}

// GetPassage fetches a bible passage from the external API.
func (c *RealBibleAIClient) GetPassage(
	ctx context.Context,
	reference string,
	version string,
) (map[string]interface{}, error) {
	if c.APIURL == "" {
		return nil, fmt.Errorf("bible API not configured")
	}

	if version == "" {
		version = "ESV"
	}

	reqPayload := QueryRequest{
		Query: QueryPayload{
			Verses: []string{reference},
		},
		Context: &QueryContext{
			User: &UserContext{
				Version: version,
			},
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
	if result.Verse == "" {
		log.Printf("Bible API response status %s, body length %d", resp.Status(), len(resp.Body()))
		// Fallback manual check
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

	verseText := html.UnescapeString(result.Verse)
	finalRef := reference // Default to user input

	// Parse reference from text if possible
	// Example: "John 3:16 (ESV) For God so loved..."
	// Group 1: Ref, Group 2: Version, Group 3: Text
	matches := verseReferenceRegex.FindStringSubmatch(verseText)

	if len(matches) == 4 {
		finalRef = matches[1]
		// version := matches[2] // We could use this too
		verseText = matches[3]
	}

	return map[string]interface{}{
		"verse":     verseText,
		"text":      verseText,
		"reference": finalRef,
		"version":   version,
	}, nil
}

// ChatCompletion sends a chat completion request to the external API.
func (c *RealBibleAIClient) ChatCompletion(
	ctx context.Context,
	payload map[string]interface{},
) (map[string]interface{}, error) {
	if c.APIURL == "" {
		return nil, fmt.Errorf("bible API not configured")
	}

	reqPayload := c.prepareQueryRequest(payload)

	// V2 Response Wrapper
	type V2Response struct {
		Data OQueryResponse         `json:"data"`
		Meta map[string]interface{} `json:"meta"`
	}
	var v2Result V2Response
	var errorResult ErrorResponse

	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("X-API-KEY", c.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(reqPayload).
		SetResult(&v2Result).
		SetError(&errorResult).
		Post(c.APIURL + "/query")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("bible API error: %s, message: %s", resp.Status(), errorResult.Error.Message)
	}

	// Use Data from V2 response
	result := v2Result.Data

	// Unescape the text before cleaning
	unescapedText := html.UnescapeString(result.Text)
	cleanedText := cleanHTML(unescapedText)

	// Construct return map
	response := map[string]interface{}{
		"text":       cleanedText,
		"references": result.References,
	}

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

// StreamChatCompletion sends a chat completion request to the external API with streaming.
func (c *RealBibleAIClient) StreamChatCompletion(
	ctx context.Context,
	payload map[string]interface{},
) (<-chan string, <-chan error, error) {
	if c.APIURL == "" {
		return nil, nil, fmt.Errorf("bible API not configured")
	}

	outChan := make(chan string)
	errChan := make(chan error, 1) // Buffered

	go func() {
		defer close(outChan)
		defer close(errChan)

		// 1. Attempt Streaming Request
		reqPayload := c.prepareQueryRequest(payload)
		if reqPayload.Options == nil {
			reqPayload.Options = &QueryOptions{}
		}
		reqPayload.Options.Stream = true

		resp, err := c.Client.R().
			SetContext(ctx).
			SetHeader("X-API-KEY", c.APIKey).
			SetHeader("Content-Type", "application/json").
			SetBody(reqPayload).
			SetDoNotParseResponse(true).
			Post(c.APIURL + "/query")

		isEventStream := resp != nil && strings.Contains(resp.Header().Get("Content-Type"), "text/event-stream")

		if err != nil || (resp != nil && resp.IsError()) || !isEventStream {
			if resp != nil && resp.RawResponse != nil && resp.RawResponse.Body != nil {
				resp.RawResponse.Body.Close()
			}
			c.performFallback(ctx, payload, outChan, errChan, err, resp, isEventStream)
			return
		}

		// 3. Process Stream
		defer resp.RawResponse.Body.Close()

		reader := bufio.NewReader(resp.RawResponse.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					errChan <- err
				}
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "data: ") {
				dataStr := strings.TrimPrefix(line, "data: ")
				if dataStr == "[DONE]" {
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
					// Handle raw text fallback for non-JSON SSE
					select {
					case outChan <- dataStr:
					case <-ctx.Done():
						return
					}
					continue
				}

				// V2 Format: {"delta": "..."}
				if delta, ok := data["delta"].(string); ok {
					select {
					case outChan <- delta:
					case <-ctx.Done():
						return
					}
					continue
				}

				// Fallbacks for older formats or other providers
				if text, ok := data["text"].(string); ok && text != "" {
					select {
					case outChan <- text:
					case <-ctx.Done():
						return
					}
				}
				if choices, ok := data["choices"].([]interface{}); ok {
					for _, c := range choices {
						if choice, ok := c.(map[string]interface{}); ok {
							if delta, ok := choice["delta"].(map[string]interface{}); ok {
								if content, ok := delta["content"].(string); ok {
									select {
									case outChan <- content:
									case <-ctx.Done():
										return
									}
								}
							}
						}
					}
				}
			}
		}
	}()

	return outChan, errChan, nil
}

func (c *RealBibleAIClient) performFallback(
	ctx context.Context,
	payload map[string]interface{},
	outChan chan<- string,
	errChan chan<- error,
	originalErr error,
	resp *resty.Response,
	isEventStream bool,
) {
	status := "nil"
	if resp != nil {
		status = resp.Status()
	}
	log.Printf("Streaming failed (err=%v, status=%s, event-stream=%v), falling back to non-streaming...",
		originalErr, status, isEventStream)

	fallbackResp, fallbackErr := c.ChatCompletion(ctx, payload)
	if fallbackErr != nil {
		select {
		case errChan <- fmt.Errorf("fallback failed: %v (original stream error: %v)", fallbackErr, originalErr):
		case <-ctx.Done():
		}
		return
	}

	if text, ok := fallbackResp["text"].(string); ok && text != "" {
		select {
		case outChan <- text:
		case <-ctx.Done():
		}
	}
}

func (c *RealBibleAIClient) prepareQueryRequest(payload map[string]interface{}) QueryRequest {
	prompt, ok := payload["prompt"].(string)
	if !ok {
		prompt = ""
	}

	if ctxText, ok := payload["context"].(string); ok && ctxText != "" {
		prompt += fmt.Sprintf(" Context: %s.", ctxText)
	}

	version, ok := payload["version"].(string)
	if !ok || version == "" {
		version = "ESV"
	}

	// Prepare Context
	aiProvider, _ := payload["ai_provider"].(string)

	queryContext := &QueryContext{
		User: &UserContext{
			Version:    version,
			AIProvider: aiProvider,
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
		r := strings.NewReplacer(
			"{PROMPT}", prompt,
			"{WORDS}", strings.Join(themes, ", "),
			"{PASSAGE}", strings.Join(verses, "\n"),
		)
		finalPrompt = r.Replace(template)
	} else {
		if len(themes) > 0 {
			finalPrompt += fmt.Sprintf(" Focus on themes: %s.", strings.Join(themes, ", "))
		}
	}

	return QueryRequest{
		Query: QueryPayload{
			Prompt: finalPrompt,
		},
		Context: queryContext,
		Options: &QueryOptions{},
	}
}

// GetSystemPrompt retrieves a configured system prompt by key.
func (c *RealBibleAIClient) GetSystemPrompt(key string) string {
	if c.SystemPrompts == nil {
		return ""
	}
	return c.SystemPrompts[key]
}

// GetVersions fetches the list of available bible versions.
func (c *RealBibleAIClient) GetVersions(ctx context.Context, params map[string]string) (map[string]interface{}, error) {
	if c.APIURL == "" {
		return nil, fmt.Errorf("bible API not configured")
	}

	var result map[string]interface{}
	var errorResult ErrorResponse

	req := c.Client.R().
		SetContext(ctx).
		SetHeader("X-API-KEY", c.APIKey).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		SetError(&errorResult)

	for k, v := range params {
		req.SetQueryParam(k, v)
	}

	resp, err := req.Get(c.APIURL + "/bible-versions")

	if err != nil {
		return nil, fmt.Errorf("resty request error: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("bible API error: %s, message: %s", resp.Status(), errorResult.Error.Message)
	}

	return result, nil
}

var (
	emptyParaRegex = regexp.MustCompile(`(?i)<p[^>]*>(\s|&nbsp;|<br\s*/?>)*</p>`)
	emptyLiRegex   = regexp.MustCompile(`(?i)<li[^>]*>(\s|&nbsp;|<br\s*/?>)*</li>`)
	newlineRegex   = regexp.MustCompile(`[\r\n]+`)
	//nolint:lll // Regex is long
	listWhitespaceRegex = regexp.MustCompile(
		`(?i)(</?ul[^>]*>|</?ol[^>]*>|</?li[^>]*>)\s+(</?ul[^>]*>|</?ol[^>]*>|</?li[^>]*>)`,
	)
	// verseReferenceRegex extracts Ref, Version, Text from "Ref (Ver) Text" format.
	// Uses (?s) to allow matching newlines in the text.
	verseReferenceRegex = regexp.MustCompile(`(?s)^([\w\s]+\d+:\d+(?:-\d+)?)\s+\(([^)]+)\)\s+(.*)$`)
)

func cleanHTML(input string) string {
	// Replace newlines with space to avoid word concatenation if used in text
	s := newlineRegex.ReplaceAllString(input, " ")

	// Remove empty paragraphs
	s = emptyParaRegex.ReplaceAllString(s, "")

	// Remove empty list items
	s = emptyLiRegex.ReplaceAllString(s, "")

	// Remove whitespace between list tags
	// Loop to handle consecutive matches (e.g., </li> <li> <li>)
	for {
		cleaned := listWhitespaceRegex.ReplaceAllString(s, "$1$2")
		if cleaned == s {
			break
		}
		s = cleaned
	}

	return strings.TrimSpace(s)
}
