package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
)

type ChatHandler struct {
	Client      services.BibleAIClient
	NoteService services.NoteServiceInterface
	DB          DBInterface
}

func NewChatHandler(client services.BibleAIClient, noteService services.NoteServiceInterface, db DBInterface) *ChatHandler {
	return &ChatHandler{Client: client, NoteService: noteService, DB: db}
}

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
func (h *ChatHandler) ChatWithAI(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	// Updated payload structure for RealBibleAIClient adapter
	payload := map[string]interface{}{
		"prompt": req.Prompt,
		"verses": []string{req.Passage},
		"themes": req.Themes,
	}

	var answer string
	aiResult, err := h.Client.ChatCompletion(r.Context(), payload)

	if err != nil {
		answer = "Error contacting AI service: " + err.Error()
	} else {
		if choices, ok := aiResult["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if msg, ok := choice["message"].(map[string]interface{}); ok {
					if content, ok := msg["content"].(string); ok {
						answer = content
					}
				}
			}
		} else if content, ok := aiResult["text"].(string); ok {
			answer = content
		} else if content, ok := aiResult["response"].(string); ok {
			answer = content
		}
	}

	// Get user ID from context (supports both real Firebase auth and mock auth)
	var userUUID uuid.UUID
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		// Real Firebase auth
		uid := token.UID
		err = h.DB.QueryRow(r.Context(), "SELECT id FROM users WHERE firebase_uid=$1", uid).Scan(&userUUID)
		if err != nil {
			http.Error(w, "User not found for saving note", http.StatusNotFound)
			return
		}
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		// Mock auth - testUserID is already the internal UUID
		var err error
		userUUID, err = uuid.Parse(testUserID)
		if err != nil {
			http.Error(w, "Invalid test user ID", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	contentJSON, err := json.Marshal(noteContent)
	if err != nil {
		http.Error(w, "Failed to marshal content", http.StatusInternalServerError)
		return
	}

	_, err = h.NoteService.CreateNote(r.Context(), userUUID.String(), noteTitle, contentJSON)

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
func (h *ChatHandler) AskAI(w http.ResponseWriter, r *http.Request) {
	var req AskAIRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	// Updated payload structure for RealBibleAIClient adapter
	payload := map[string]interface{}{
		"prompt":  req.Prompt,
		"context": req.Context,
	}

	var answer string
	aiResult, err := h.Client.ChatCompletion(r.Context(), payload)

	if err != nil {
		answer = "Error contacting AI service: " + err.Error()
	} else {
		if choices, ok := aiResult["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if msg, ok := choice["message"].(map[string]interface{}); ok {
					if content, ok := msg["content"].(string); ok {
						answer = content
					}
				}
			}
		} else if content, ok := aiResult["text"].(string); ok {
			answer = content
		} else if content, ok := aiResult["response"].(string); ok {
			answer = content
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ChatResponse{Response: answer}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
