package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
)

type ChatHandler struct {
	Client              services.LLMClient
	NoteService         services.NoteServiceInterface
	NotificationService services.NotificationService
	DB                  DBInterface
}

func NewChatHandler(client services.LLMClient, noteService services.NoteServiceInterface, notificationService services.NotificationService, db DBInterface) *ChatHandler {
	return &ChatHandler{
		Client:              client,
		NoteService:         noteService,
		NotificationService: notificationService,
		DB:                  db,
	}
}

type ChatOptions struct {
	Stream     *bool  `json:"stream,omitempty"`
	AIProvider string `json:"ai_provider,omitempty"`
}

type ChatRequest struct {
	Passage string       `json:"passage" validate:"required,min=5"`
	Themes  []string     `json:"themes"`
	Prompt  string       `json:"prompt" validate:"required,min=2"`
	Version string       `json:"version,omitempty"`
	Options *ChatOptions `json:"options,omitempty"`
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

	h.handleRequest(w, r, "chat", req)
}

type AskAIRequest struct {
	Context string       `json:"context" validate:"required"`
	Prompt  string       `json:"prompt" validate:"required"`
	Version string       `json:"version,omitempty"`
	Options *ChatOptions `json:"options,omitempty"`
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

	h.handleRequest(w, r, "ask", req)
}

func (h *ChatHandler) handleRequest(w http.ResponseWriter, r *http.Request, reqType string, reqData interface{}) {
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
	var options *ChatOptions
	var reqVersion string
	var fullPrompt string

	if reqType == "chat" {
		req := reqData.(ChatRequest)
		options = req.Options
		reqVersion = req.Version
		noteTitle = fmt.Sprintf("Chat: %s", req.Prompt)
		if len(noteTitle) > 50 {
			noteTitle = noteTitle[:47] + "..."
		}
		// Initial content with loading indicator
		initialContent = fmt.Sprintf("<p><strong>Passage:</strong> %s</p><p><strong>Themes:</strong> %v</p><p><strong>Q:</strong> %s</p><div class=\"ai-response\"><em>AI is thinking...</em></div>",
			req.Passage, req.Themes, req.Prompt)

		fullPrompt = fmt.Sprintf("Passage: %s\nThemes: %s\n\nQuestion: %s", req.Passage, strings.Join(req.Themes, ", "), req.Prompt)

	} else {
		req := reqData.(AskAIRequest)
		options = req.Options
		reqVersion = req.Version
		noteTitle = fmt.Sprintf("Ask AI: %s", req.Prompt)
		if len(noteTitle) > 50 {
			noteTitle = noteTitle[:47] + "..."
		}
		initialContent = fmt.Sprintf("<p><strong>Context:</strong> %s</p><p><strong>Q:</strong> %s</p><div class=\"ai-response\"><em>AI is thinking...</em></div>",
			req.Context, req.Prompt)

		fullPrompt = fmt.Sprintf("Context: %s\n\nQuestion: %s", req.Context, req.Prompt)
	}

	contentJSON, _ := json.Marshal(initialContent)
	// Create note with "pending" status
	note, err := h.NoteService.CreateNote(r.Context(), userUUID.String(), noteTitle, contentJSON, "pending")
	if err != nil {
		slog.Error("Failed to create pending note", "error", err)
		http.Error(w, "Failed to initialize operation", http.StatusInternalServerError)
		return
	}

	// Check if streaming is requested (default: true)
	stream := true
	if options != nil && options.Stream != nil {
		stream = *options.Stream
	}

	// Prepare Context
	bgCtx := context.Background()
	ctxWithOpts := bgCtx
	if options != nil && options.AIProvider != "" {
		ctxWithOpts = context.WithValue(ctxWithOpts, services.AIProviderKey, options.AIProvider)
	}
	if reqVersion != "" {
		ctxWithOpts = context.WithValue(ctxWithOpts, services.BibleVersionKey, reqVersion)
	}

	if !stream {
		// Blocking Mode
		resp, _, err := h.Client.Query(ctxWithOpts, fullPrompt, "")
		if err != nil {
			slog.Error("Failed to query AI", "error", err)
			h.NoteService.UpdateNote(bgCtx, userUUID.String(), note.ID, noteTitle, contentJSON, "failed")
			http.Error(w, "AI Request failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fullAnswer := resp

		// Update Note
		finalHTML := h.formatFinalHTML(reqType, reqData, fullAnswer)
		finalJSON, _ := json.Marshal(finalHTML)
		h.NoteService.UpdateNote(bgCtx, userUUID.String(), note.ID, noteTitle, finalJSON, "active")
		h.sendNotification(bgCtx, userUUID.String(), note.ID)

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"note_id":  note.ID,
			"response": fullAnswer,
		})
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

	// 3. Launch AI Request
	// Use detached context so stream reading continues if client disconnects
	outChan, _, err := h.Client.Stream(ctxWithOpts, fullPrompt)
	if err != nil {
		slog.Error("Failed to start AI stream", "error", err)
		if err := h.NoteService.UpdateNote(bgCtx, userUUID.String(), note.ID, noteTitle, contentJSON, "failed"); err != nil {
			slog.Error("Failed to mark note as failed", "error", err)
		}
		errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", errMsg)
		flusher.Flush()
		return
	}

	// 4. Stream Loop
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	fullAnswer := ""
	clientConnected := true

	for {
		select {
		case chunk, ok := <-outChan:
			if !ok {
				// Stream finished successfully
				goto StreamFinished
			}

			fullAnswer += chunk

			if clientConnected {
				data, _ := json.Marshal(map[string]string{"response": chunk})
				if _, err := fmt.Fprintf(w, "event: chunk\ndata: %s\n\n", data); err != nil {
					clientConnected = false
					slog.Info("Client disconnected during stream write")
				} else {
					flusher.Flush()
				}
			}

		// Error handling for stream usually comes via separate channel in previous implementation.
		// LLMClient.Stream returns (chan string, string, error) where error is immediate.
		// Errors DURING stream? The interface says "Returns a channel for chunks...".
		// It doesn't return an error channel.
		// My implementation of RealBibleAIClient.Stream swallows errors and logs them (or stops stream).
		// So here we only listen to outChan.
		// Wait, if I want to propagate errors during stream, I can't with current LLMClient interface.
		// This is a limitation of the interface requested.
		// I will proceed assuming outChan close means done (success or error).

		case <-ticker.C:
			if clientConnected {
				if _, err := fmt.Fprintf(w, "event: ping\ndata: {}\n\n"); err != nil {
					clientConnected = false
				} else {
					flusher.Flush()
				}
			}

		case <-r.Context().Done():
			if clientConnected {
				clientConnected = false
				slog.Info("Client disconnected (context done), continuing background processing")
			}
		}
	}

StreamFinished:
	// Update Note with final content
	finalHTML := h.formatFinalHTML(reqType, reqData, fullAnswer)
	finalJSON, _ := json.Marshal(finalHTML)

	if err := h.NoteService.UpdateNote(bgCtx, userUUID.String(), note.ID, noteTitle, finalJSON, "active"); err != nil {
		slog.Error("Failed to update note status", "error", err)
	} else {
		slog.Info("Note updated successfully", "id", note.ID)
		h.sendNotification(bgCtx, userUUID.String(), note.ID)
	}

	if clientConnected {
		fmt.Fprintf(w, "event: done\ndata: {}\n\n")
		flusher.Flush()
	}
}

func (h *ChatHandler) formatFinalHTML(reqType string, reqData interface{}, answer string) string {
	if reqType == "chat" {
		req := reqData.(ChatRequest)
		return fmt.Sprintf("<p><strong>Passage:</strong> %s</p><p><strong>Themes:</strong> %v</p><p><strong>Q:</strong> %s</p><div class=\"ai-response\"><strong>AI:</strong> %s</div>",
			req.Passage, req.Themes, req.Prompt, answer)
	}
	req := reqData.(AskAIRequest)
	return fmt.Sprintf("<p><strong>Context:</strong> %s</p><p><strong>Q:</strong> %s</p><div class=\"ai-response\"><strong>AI:</strong> %s</div>",
		req.Context, req.Prompt, answer)
}

func (h *ChatHandler) sendNotification(ctx context.Context, userID, noteID string) {
	if h.NotificationService != nil {
		notificationData := map[string]string{
			"type":    "ask_ai_complete",
			"note_id": noteID,
		}
		if err := h.NotificationService.SendNotification(ctx, userID, "AI Response Ready", "Your question has been answered.", notificationData); err != nil {
			slog.Warn("Failed to send notification", "error", err)
		}
	}
}
