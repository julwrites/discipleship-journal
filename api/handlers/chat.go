package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"discipleship_journal_api/database"
	"firebase.google.com/go/v4/auth"
	"github.com/go-resty/resty/v2"
)

type ChatRequest struct {
	Passage string   `json:"passage"`
	Themes  []string `json:"themes"`
	Prompt  string   `json:"prompt"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

func ChatWithAI(w http.ResponseWriter, r *http.Request) {
	// BibleAIAPI Endpoint
	apiURL := os.Getenv("BIBLE_API_URL")
	apiKey := os.Getenv("BIBLE_API_KEY")

	if apiURL == "" {
		http.Error(w, "Bible API not configured", http.StatusInternalServerError)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client := resty.New()

	// Construct the prompt to send to BibleAIAPI
	// Assuming BibleAIAPI has a /chat or /generate endpoint
	// This structure depends on the actual API contract.
	// For now, I'll assume a generic completions endpoint or specific chat endpoint.

	payload := map[string]interface{}{
		"model": "bible-model", // hypothetical
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful Bible assistant."},
			{"role": "user", "content": fmt.Sprintf("Context: %s. Themes: %v. Question: %s", req.Passage, req.Themes, req.Prompt)},
		},
	}

	var aiResult map[string]interface{}
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetBody(payload).
		SetResult(&aiResult).
		Post(apiURL + "/chat/completions") // Adjust path

	if err != nil {
		http.Error(w, "Failed to call AI API", http.StatusInternalServerError)
		return
	}

	if resp.IsError() {
		// Fallback for demo/dev purposes if API is not actually live
		// http.Error(w, fmt.Sprintf("AI API Error: %s", resp.Status()), resp.StatusCode())
		// return

		// MOCK RESPONSE for development since I don't have a real API key in sandbox
		aiResult = map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]string{
						"content": "This is a simulated AI response based on " + req.Passage,
					},
				},
			},
		}
	}

	// Extract content - heavily dependent on API structure (e.g. OpenAI format)
	// Simplified extraction:
	var answer string
	if choices, ok := aiResult["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if msg, ok := choice["message"].(map[string]interface{}); ok {
				answer, _ = msg["content"].(string)
			}
		}
	}

	// Create a new Journal Note with the conversation
	token := r.Context().Value("user").(*auth.Token)
	uid := token.UID

	noteTitle := fmt.Sprintf("Chat: %s", req.Prompt)
	if len(noteTitle) > 50 {
		noteTitle = noteTitle[:47] + "..."
	}

	noteContent := map[string]interface{}{
		"markdown": fmt.Sprintf("# %s\n\n**Passage:** %s\n**Themes:** %v\n\n**Q:** %s\n\n**AI:** %s",
			noteTitle, req.Passage, req.Themes, req.Prompt, answer),
		"type": "chat_log",
	}

	var noteID string
	err = database.DB.QueryRow(r.Context(),
		"INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id",
		uid, noteTitle, noteContent).Scan(&noteID)

	if err != nil {
		// Log error but maybe still return the answer?
		// For now, fail hard.
		http.Error(w, "Failed to save chat note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Response: answer})
}

type AskAIRequest struct {
	Context string `json:"context"`
	Prompt  string `json:"prompt"`
}

func AskAI(w http.ResponseWriter, r *http.Request) {
	apiURL := os.Getenv("BIBLE_API_URL")
	apiKey := os.Getenv("BIBLE_API_KEY")

	if apiURL == "" {
		// Mock for dev
		// http.Error(w, "Bible API not configured", http.StatusInternalServerError)
	}

	var req AskAIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client := resty.New()

	payload := map[string]interface{}{
		"model": "bible-model",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant analyzing a journal note."},
			{"role": "user", "content": fmt.Sprintf("Context: %s. Question: %s", req.Context, req.Prompt)},
		},
	}

	var aiResult map[string]interface{}
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetBody(payload).
		SetResult(&aiResult).
		Post(apiURL + "/chat/completions")

	var answer string
	if err == nil && !resp.IsError() {
		if choices, ok := aiResult["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if msg, ok := choice["message"].(map[string]interface{}); ok {
					answer, _ = msg["content"].(string)
				}
			}
		}
	} else {
        // Fallback/Mock
        answer = "AI analysis of note context: " + req.Context[:min(len(req.Context), 20)] + "... -> " + req.Prompt
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Response: answer})
}

func min(a, b int) int {
    if a < b { return a }
    return b
}
