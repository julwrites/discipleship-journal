package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"discipleship_journal_api/database"
	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/go-resty/resty/v2"
)

type ChatRequest struct {
	Passage string   `json:"passage" validate:"required,min=5"`
	Themes  []string `json:"themes"`
	Prompt  string   `json:"prompt" validate:"required,min=2"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

// ChatWithAI godoc
// @Summary Chat with AI
// @Description Chat with Bible AI using passage context
// @Tags chat
// @Accept json
// @Produce json
// @Param request body ChatRequest true "Chat Request"
// @Success 200 {object} ChatResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/chat [post]
func ChatWithAI(w http.ResponseWriter, r *http.Request) {
	// BibleAIAPI Endpoint
	apiURL := os.Getenv("BIBLE_API_URL")
	apiKey := os.Getenv("BIBLE_API_KEY")

	// if apiURL == "" {
	// 	// Mock environment logic if needed, or error
	// }

	var req ChatRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	client := resty.New()

	payload := map[string]interface{}{
		"model": "bible-model",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful Bible assistant."},
			{"role": "user", "content": fmt.Sprintf("Context: %s. Themes: %v. Question: %s", req.Passage, req.Themes, req.Prompt)},
		},
	}

	var aiResult map[string]interface{}

	// Only call if URL is present, otherwise simulate
	var answer string
	if apiURL != "" {
		resp, err := client.R().
			SetHeader("Authorization", "Bearer "+apiKey).
			SetBody(payload).
			SetResult(&aiResult).
			Post(apiURL + "/chat/completions")

		if err != nil || resp.IsError() {
			// Handle error or fallback
			answer = "Error contacting AI service."
		} else {
			if choices, ok := aiResult["choices"].([]interface{}); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]interface{}); ok {
					if msg, ok := choice["message"].(map[string]interface{}); ok {
						answer, _ = msg["content"].(string)
					}
				}
			}
		}
	} else {
		// Simulation
		answer = "This is a simulated AI response based on " + req.Passage
	}

	// Create a new Journal Note with the conversation
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID

	userUUID, err := GetUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found for saving note", http.StatusNotFound)
		return
	}

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
		userUUID, noteTitle, noteContent).Scan(&noteID)

	if err != nil {
		http.Error(w, "Failed to save chat note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ChatResponse{Response: answer}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

type AskAIRequest struct {
	Context string `json:"context" validate:"required"`
	Prompt  string `json:"prompt" validate:"required"`
}

// AskAI godoc
// @Summary Ask AI
// @Description Ask AI a question based on general context
// @Tags chat
// @Accept json
// @Produce json
// @Param request body AskAIRequest true "Ask AI Request"
// @Success 200 {object} ChatResponse
// @Failure 400 {object} map[string]string
// @Router /api/ai/ask [post]
func AskAI(w http.ResponseWriter, r *http.Request) {
	apiURL := os.Getenv("BIBLE_API_URL")
	apiKey := os.Getenv("BIBLE_API_KEY")

	var req AskAIRequest
	if !DecodeAndValidate(w, r, &req) {
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
	var answer string

	if apiURL != "" {
		resp, err := client.R().
			SetHeader("Authorization", "Bearer "+apiKey).
			SetBody(payload).
			SetResult(&aiResult).
			Post(apiURL + "/chat/completions")

		if err == nil && !resp.IsError() {
			if choices, ok := aiResult["choices"].([]interface{}); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]interface{}); ok {
					if msg, ok := choice["message"].(map[string]interface{}); ok {
						answer, _ = msg["content"].(string)
					}
				}
			}
		}
	} else {
        // Fallback/Mock
        answer = "AI analysis of note context: " + req.Context[:min(len(req.Context), 20)] + "... -> " + req.Prompt
    }

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ChatResponse{Response: answer}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func min(a, b int) int {
    if a < b { return a }
    return b
}
