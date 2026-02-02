package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/models"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	chi "github.com/go-chi/chi/v5"
)

type Note struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Title     string          `json:"title"`
	Content   json.RawMessage `json:"content,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Tags      []services.Tag  `json:"tags"`
}

type CreateNoteRequest struct {
	Title   string      `json:"title" validate:"max=100"`
	Content interface{} `json:"content"`
	Tags    []string    `json:"tags"`
}

type CreateTagRequest struct {
	Name string `json:"name" validate:"required,max=50"`
}

type NoteHandler struct {
	db          DBInterface
	noteService services.NoteServiceInterface
}

func NewNoteHandler(db DBInterface, noteService services.NoteServiceInterface) *NoteHandler {
	return &NoteHandler{db: db, noteService: noteService}
}

func (h *NoteHandler) getUserUUID(ctx context.Context, firebaseUID string) (string, error) {
	if h.db == nil {
		return "", fmt.Errorf("database connection is nil")
	}
	var id string
	err := h.db.QueryRow(ctx, "SELECT id FROM users WHERE firebase_uid=$1", firebaseUID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

type Meta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

type NotesResponse struct {
	Data []Note `json:"data"`
	Meta Meta   `json:"meta"`
}

// GetNotes godoc
// @Summary Get all notes for user
// @Description Fetch all notes belonging to the authenticated user, with pagination, search, and filtering
// @Tags notes
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param q query string false "Search query"
// @Param startDate query string false "Start Date (RFC3339)"
// @Param endDate query string false "End Date (RFC3339)"
// @Param sortBy query string false "Sort By (updated_at, created_at, title)"
// @Param sortOrder query string false "Sort Order (asc, desc)"
// @Success 200 {object} NotesResponse
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes [get]
func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	if h.noteService == nil {
		slog.Error("NoteService is nil in NoteHandler.GetNotes")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var uid string
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		uid = token.UID
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		// Test environment override
		// If TestUserKey is present, we assume it is the internal UUID
		h.getNotesWithUUID(w, r, testUserID)
		return
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userUUID, err := h.getUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	h.getNotesWithUUID(w, r, userUUID)
}

func (h *NoteHandler) getNotesWithUUID(w http.ResponseWriter, r *http.Request, userUUID string) {
	// Pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Filter parameters
	filter := services.NoteFilter{
		SearchQuery: r.URL.Query().Get("q"),
		Tag:         r.URL.Query().Get("tag"),
		SortBy:      r.URL.Query().Get("sortBy"),
		SortOrder:   r.URL.Query().Get("sortOrder"),
	}

	if startDateStr := r.URL.Query().Get("startDate"); startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filter.StartDate = &t
		}
	}
	if endDateStr := r.URL.Query().Get("endDate"); endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filter.EndDate = &t
		}
	}

	serviceNotes, total, err := h.noteService.GetNotes(r.Context(), userUUID, page, limit, filter)
	if err != nil {
		slog.Error("Failed to fetch notes", "error", err, "user_id", userUUID)
		http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
		return
	}

	// Convert services.Note to handlers.Note
	var notes []Note
	for _, sn := range serviceNotes {
		notes = append(notes, Note{
			ID:        sn.ID,
			UserID:    sn.UserID,
			Title:     sn.Title,
			Content:   sn.Content,
			CreatedAt: sn.CreatedAt,
			UpdatedAt: sn.UpdatedAt,
			Tags:      sn.Tags,
		})
	}

	if notes == nil {
		notes = []Note{}
	}

	totalPages := (total + limit - 1) / limit
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	resp := NotesResponse{
		Data: notes,
		Meta: Meta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteNote godoc
// @Summary Delete a note
// @Description Delete a journal note by ID
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Success 200
// @Failure 403 {string} string "Unauthorized"
// @Failure 404 {string} string "Note not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes/{id} [delete]
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	if h.noteService == nil {
		slog.Error("NoteService is nil in NoteHandler.DeleteNote")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var userUUID string
	var err error
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		userUUID, err = h.getUserUUID(r.Context(), token.UID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		userUUID = testUserID
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	noteID := chi.URLParam(r, "id")
	err = h.noteService.DeleteNote(r.Context(), userUUID, noteID)

	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Note not found or unauthorized", http.StatusNotFound)
		} else {
			slog.Error("Failed to delete note", "error", err, "note_id", noteID)
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// CreateNote godoc
// @Summary Create a new note
// @Description Create a new journal note
// @Tags notes
// @Accept json
// @Produce json
// @Param request body CreateNoteRequest true "Create Note Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes [post]
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	if h.noteService == nil {
		slog.Error("NoteService is nil in NoteHandler.CreateNote")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var userUUID string
	var err error
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		userUUID, err = h.getUserUUID(r.Context(), token.UID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		userUUID = testUserID
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateNoteRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	contentJSON, err := json.Marshal(req.Content)
	if err != nil {
		http.Error(w, "Invalid content", http.StatusBadRequest)
		return
	}

	note, err := h.noteService.CreateNote(r.Context(), userUUID, req.Title, contentJSON, req.Tags)
	if err != nil {
		slog.Error("Failed to create note", "error", err, "user_id", userUUID)
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"id": note.ID}); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateNote godoc
// @Summary Update a note
// @Description Update an existing journal note
// @Tags notes
// @Accept json
// @Produce json
// @Param request body CreateNoteRequest true "Update Note Request"
// @Success 200
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {string} string "Unauthorized"
// @Failure 404 {string} string "Note not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes/{id} [put]
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	if h.noteService == nil {
		slog.Error("NoteService is nil in NoteHandler.UpdateNote")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var userUUID string
	var err error
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		userUUID, err = h.getUserUUID(r.Context(), token.UID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		userUUID = testUserID
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	noteID := chi.URLParam(r, "id")

	var req CreateNoteRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	contentJSON, err := json.Marshal(req.Content)
	if err != nil {
		http.Error(w, "Invalid content", http.StatusBadRequest)
		return
	}

	err = h.noteService.UpdateNote(r.Context(), userUUID, noteID, req.Title, contentJSON, req.Tags)

	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			slog.Error("Failed to update note", "error", err, "note_id", noteID)
			http.Error(w, "Failed to update note", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// GetNote godoc
// @Summary Get a single note
// @Description Fetch a specific note by ID
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Success 200 {object} Note
// @Failure 403 {string} string "Unauthorized"
// @Failure 404 {string} string "Note not found"
// @Router /api/notes/{id} [get]
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	if h.noteService == nil {
		slog.Error("NoteService is nil in NoteHandler.GetNote")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var userUUID string
	var err error
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		userUUID, err = h.getUserUUID(r.Context(), token.UID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		userUUID = testUserID
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	noteID := chi.URLParam(r, "id")

	sn, err := h.noteService.GetNote(r.Context(), userUUID, noteID)

	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			slog.Error("Failed to get note", "error", err, "note_id", noteID)
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	note := Note{
		ID:        sn.ID,
		UserID:    sn.UserID,
		Title:     sn.Title,
		Content:   sn.Content,
		CreatedAt: sn.CreatedAt,
		UpdatedAt: sn.UpdatedAt,
		Tags:      sn.Tags,
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetTags godoc
// @Summary Get all tags for user
// @Description Fetch all tags belonging to the authenticated user
// @Tags tags
// @Accept json
// @Produce json
// @Success 200 {object} []services.Tag
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/tags [get]
func (h *NoteHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	if h.noteService == nil {
		slog.Error("NoteService is nil in NoteHandler.GetTags")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var userUUID string
	var err error
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		userUUID, err = h.getUserUUID(r.Context(), token.UID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		userUUID = testUserID
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tags, err := h.noteService.GetUserTags(r.Context(), userUUID)
	if err != nil {
		slog.Error("Failed to fetch tags", "error", err, "user_id", userUUID)
		http.Error(w, "Failed to fetch tags", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(tags); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new tag
// @Tags tags
// @Accept json
// @Produce json
// @Param request body CreateTagRequest true "Create Tag Request"
// @Success 200 {object} services.Tag
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/tags [post]
func (h *NoteHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	if h.noteService == nil {
		slog.Error("NoteService is nil in NoteHandler.CreateTag")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var userUUID string
	var err error
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		userUUID, err = h.getUserUUID(r.Context(), token.UID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		userUUID = testUserID
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateTagRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	tag, err := h.noteService.CreateTag(r.Context(), userUUID, req.Name)
	if err != nil {
		slog.Error("Failed to create tag", "error", err, "user_id", userUUID)
		http.Error(w, "Failed to create tag", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tag); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteTag godoc
// @Summary Delete a tag
// @Description Delete a tag by ID
// @Tags tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Success 200
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/tags/{id} [delete]
func (h *NoteHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	if h.noteService == nil {
		slog.Error("NoteService is nil in NoteHandler.DeleteTag")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var userUUID string
	var err error
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		userUUID, err = h.getUserUUID(r.Context(), token.UID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		userUUID = testUserID
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tagID := chi.URLParam(r, "id")
	err = h.noteService.DeleteTag(r.Context(), userUUID, tagID)

	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Tag not found", http.StatusNotFound)
		} else {
			slog.Error("Failed to delete tag", "error", err, "tag_id", tagID)
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}
