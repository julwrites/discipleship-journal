package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	chi "github.com/go-chi/chi/v5"
	pgx "github.com/jackc/pgx/v5"
)

type Note struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Title     string                 `json:"title"`
	Content   map[string]interface{} `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type CreateNoteRequest struct {
	Title   string                 `json:"title" validate:"required,min=1,max=100"`
	Content map[string]interface{} `json:"content" validate:"required"`
}

type NoteHandler struct {
	db          DBInterface
	noteService services.NoteServiceInterface
}

func NewNoteHandler(db DBInterface, noteService services.NoteServiceInterface) *NoteHandler {
	return &NoteHandler{db: db, noteService: noteService}
}

func (h *NoteHandler) getUserUUID(ctx context.Context, firebaseUID string) (string, error) {
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
// @Description Fetch all notes belonging to the authenticated user, with pagination and search
// @Tags notes
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param q query string false "Search query"
// @Success 200 {object} NotesResponse
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes [get]
func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID

	userUUID, err := h.getUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	searchQuery := r.URL.Query().Get("q")

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

	offset := (page - 1) * limit

	var rows pgx.Rows
	var qErr error
	var total int

	// Count total notes
	countQuery := "SELECT COUNT(*) FROM notes WHERE user_id=$1"
	if searchQuery != "" {
		countQuery += " AND (title ILIKE $2 OR content::text ILIKE $2)"
		err = h.db.QueryRow(r.Context(), countQuery, userUUID, "%"+searchQuery+"%").Scan(&total)
	} else {
		err = h.db.QueryRow(r.Context(), countQuery, userUUID).Scan(&total)
	}

	if err != nil {
		http.Error(w, "Failed to count notes", http.StatusInternalServerError)
		return
	}

	baseQuery := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id=$1"

	if searchQuery != "" {
		baseQuery += " AND (title ILIKE $2 OR content::text ILIKE $2)"
		baseQuery += fmt.Sprintf(" ORDER BY updated_at DESC LIMIT %d OFFSET %d", limit, offset)
		rows, qErr = h.db.Query(r.Context(), baseQuery, userUUID, "%"+searchQuery+"%")
	} else {
		baseQuery += fmt.Sprintf(" ORDER BY updated_at DESC LIMIT %d OFFSET %d", limit, offset)
		rows, qErr = h.db.Query(r.Context(), baseQuery, userUUID)
	}

	if qErr != nil {
		http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			continue
		}
		notes = append(notes, n)
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
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID
	noteID := chi.URLParam(r, "id")

	userUUID, err := h.getUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = h.noteService.DeleteNote(r.Context(), userUUID, noteID)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Note not found or unauthorized", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
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
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID

	userUUID, err := h.getUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
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

	note, err := h.noteService.CreateNote(r.Context(), userUUID, req.Title, contentJSON)
	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"id": note.ID}); err != nil {
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
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID
	noteID := chi.URLParam(r, "id")

	userUUID, err := h.getUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var req CreateNoteRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	// Verify ownership
	var ownerID string
	err = h.db.QueryRow(r.Context(), "SELECT user_id FROM notes WHERE id=$1", noteID).Scan(&ownerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if ownerID != userUUID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	_, err = h.db.Exec(r.Context(),
		"UPDATE notes SET title=$1, content=$2, updated_at=NOW() WHERE id=$3",
		req.Title, req.Content, noteID)

	if err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID
	noteID := chi.URLParam(r, "id")

	userUUID, err := h.getUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var n Note
	err = h.db.QueryRow(r.Context(), "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id=$1", noteID).Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if n.UserID != userUUID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if err := json.NewEncoder(w).Encode(n); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
