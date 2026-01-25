package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
)

type ChatHandler struct {
	Client              services.BibleAIClient
	NoteService         services.NoteServiceInterface
	NotificationService services.NotificationService
	DB                  DBInterface
}

func NewChatHandler(client services.BibleAIClient, noteService services.NoteServiceInterface, notificationService services.NotificationService, db DBInterface) *ChatHandler {
	return &ChatHandler{
		Client:              client,
		NoteService:         noteService,
		NotificationService: notificationService,
		DB:                  db,
	}
}

type ChatRequest struct {
	Passage string   `json:"passage" validate:"required,min=5"`
	Themes  []string `json:"themes"`
	Prompt  string   `json:"prompt" validate:"required,min=2"`
	Version string   `json:"version,omitempty"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

// ChatWithAI godoc
// @Summary Chat with AI
// @Description Chat with Bible AI using passage context
// @Tags chat
// @Accept json
// @Produce text/event-stream
// @Param request body ChatRequest true "Chat Request"
// @Success 200 {string} string "Streamed events"
// @Failure 400 {object} map[string]string
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/chat [post]
func (h *ChatHandler) ChatWithAI(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	h.handleStream(w, r, "chat", req)
}

type AskAIRequest struct {
	Context string `json:"context" validate:"required"`
	Prompt  string `json:"prompt" validate:"required"`
	Version string `json:"version,omitempty"`
}

// AskAI godoc
// @Summary Ask AI
// @Description Ask AI a question based on general context
// @Tags chat
// @Accept json
// @Produce text/event-stream
// @Param request body AskAIRequest true "Ask AI Request"
// @Success 200 {string} string "Streamed events"
// @Failure 400 {object} map[string]string
// @Router /api/ai/ask [post]
func (h *ChatHandler) AskAI(w http.ResponseWriter, r *http.Request) {
	var req AskAIRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	h.handleStream(w, r, "ask", req)
}

func (h *ChatHandler) handleStream(w http.ResponseWriter, r *http.Request, reqType string, reqData interface{}) {
	// 1. Authenticate & Create Note
	var userUUID uuid.UUID
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		uid := token.UID
		if err := h.DB.QueryRow(r.Context(), "SELECT id FROM users WHERE firebase_uid=$1", uid).Scan(&userUUID); err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
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

	var noteTitle string
	var initialContent string
	var payload map[string]interface{}

	if reqType == "chat" {
		req := reqData.(ChatRequest)
		noteTitle = fmt.Sprintf("Chat: %s", req.Prompt)
		if len(noteTitle) > 50 {
			noteTitle = noteTitle[:47] + "..."
		}
		// Initial content with loading indicator
		initialContent = fmt.Sprintf("<p><strong>Passage:</strong> %s</p><p><strong>Themes:</strong> %v</p><p><strong>Q:</strong> %s</p><div class=\"ai-response\"><em>AI is thinking...</em></div>",
			req.Passage, req.Themes, req.Prompt)

		payload = map[string]interface{}{
			"prompt":  req.Prompt,
			"verses":  []string{req.Passage},
			"themes":  req.Themes,
			"type":    "ask",
			"version": req.Version,
		}
	} else {
		req := reqData.(AskAIRequest)
		noteTitle = fmt.Sprintf("Ask AI: %s", req.Prompt)
		if len(noteTitle) > 50 {
			noteTitle = noteTitle[:47] + "..."
		}
		initialContent = fmt.Sprintf("<p><strong>Context:</strong> %s</p><p><strong>Q:</strong> %s</p><div class=\"ai-response\"><em>AI is thinking...</em></div>",
			req.Context, req.Prompt)

		payload = map[string]interface{}{
			"prompt":  req.Prompt,
			"context": req.Context,
			"type":    "ask",
			"version": req.Version,
		}
	}

	contentJSON, _ := json.Marshal(initialContent)
	// Create note with "pending" status
	note, err := h.NoteService.CreateNote(r.Context(), userUUID.String(), noteTitle, contentJSON, "pending")
	if err != nil {
		slog.Error("Failed to create pending note", "error", err)
		http.Error(w, "Failed to initialize operation", http.StatusInternalServerError)
		return
	}

	// 2. Setup Streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable Nginx buffering if any

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send initial event
	if _, err := fmt.Fprintf(w, "event: start\ndata: %s\n\n", note.ID); err != nil {
		return
	}
	flusher.Flush()

	// 3. Launch AI Request in Background (Detached Context)
	// Use buffered channels to prevent goroutine leak on disconnect
	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Create a detached context for the background operation so it survives client disconnect
	bgCtx := context.Background()
	// (In a real app, might want to link this to server shutdown context)

	go func() {
		defer close(resultChan)
		defer close(errChan)

		// Call AI (Synchronous for now)
		aiResult, err := h.Client.ChatCompletion(bgCtx, payload)
		if err != nil {
			errChan <- err
			return
		}

		var answer string
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

		// Update Note with final content
		var finalHTML string
		if reqType == "chat" {
			req := reqData.(ChatRequest)
			finalHTML = fmt.Sprintf("<p><strong>Passage:</strong> %s</p><p><strong>Themes:</strong> %v</p><p><strong>Q:</strong> %s</p><div class=\"ai-response\"><strong>AI:</strong> %s</div>",
				req.Passage, req.Themes, req.Prompt, answer)
		} else {
			req := reqData.(AskAIRequest)
			finalHTML = fmt.Sprintf("<p><strong>Context:</strong> %s</p><p><strong>Q:</strong> %s</p><div class=\"ai-response\"><strong>AI:</strong> %s</div>",
				req.Context, req.Prompt, answer)
		}

		finalJSON, _ := json.Marshal(finalHTML)

		// Update DB with detached context
		if err := h.NoteService.UpdateNote(bgCtx, userUUID.String(), note.ID, noteTitle, finalJSON, "active"); err != nil {
			slog.Error("Failed to update note status", "error", err)
		} else {
			slog.Info("Note updated successfully", "id", note.ID)

			// Send Notification
			if h.NotificationService != nil {
				notificationData := map[string]string{
					"type": "ask_ai_complete",
					"note_id": note.ID,
				}
				if err := h.NotificationService.SendNotification(bgCtx, userUUID.String(), "AI Response Ready", "Your question has been answered.", notificationData); err != nil {
					slog.Warn("Failed to send notification", "error", err)
				}
			}
		}

		resultChan <- answer
	}()

	// 4. Stream Loop
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case answer := <-resultChan:
			// Send the content (simulating chunks could be done here if needed, but we send full for now)
			// Use json.Marshal to safely escape newlines/etc in SSE data
			data, _ := json.Marshal(map[string]string{"response": answer})
			if _, err := fmt.Fprintf(w, "event: chunk\ndata: %s\n\n", data); err != nil {
				return
			}
			if _, err := fmt.Fprintf(w, "event: done\ndata: {}\n\n"); err != nil {
				return
			}
			flusher.Flush()
			return

		case err := <-errChan:
			// Update note status to failed
			if err := h.NoteService.UpdateNote(context.Background(), userUUID.String(), note.ID, noteTitle, contentJSON, "failed"); err != nil {
				slog.Error("Failed to mark note as failed", "error", err)
			}

			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			if _, err := fmt.Fprintf(w, "event: error\ndata: %s\n\n", errMsg); err != nil {
				return
			}
			flusher.Flush()
			return

		case <-ticker.C:
			if _, err := fmt.Fprintf(w, "event: ping\ndata: {}\n\n"); err != nil {
				return
			}
			flusher.Flush()

		case <-r.Context().Done():
			// Client disconnected
			slog.Info("Client disconnected, background processing continuing")
			return
		}
	}
}
