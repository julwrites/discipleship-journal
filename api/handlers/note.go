package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"discipleship_journal_api/database"
	"discipleship_journal_api/middleware"
	"discipleship_journal_api/validation"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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

// GetNotes godoc
// @Summary Get all notes for user
// @Description Fetch all notes belonging to the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Success 200 {array} Note
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes [get]
func GetNotes(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID

	userUUID, err := GetUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	rows, err := database.DB.Query(r.Context(), "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id=$1 ORDER BY updated_at DESC", userUUID)
	if err != nil {
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

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
// @Failure 400 {object} map[string]string
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes [post]
func CreateNote(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID

	userUUID, err := GetUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var req CreateNoteRequest
	if !validation.DecodeAndValidate(w, r, &req) {
		return
	}

	var noteID string
	err = database.DB.QueryRow(r.Context(),
		"INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id",
		userUUID, req.Title, req.Content).Scan(&noteID)

	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"id": noteID}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateNote godoc
// @Summary Update a note
// @Description Update an existing journal note
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Param request body CreateNoteRequest true "Update Note Request"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 403 {string} string "Unauthorized"
// @Failure 404 {string} string "Note not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notes/{id} [put]
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID
	noteID := chi.URLParam(r, "id")

	userUUID, err := GetUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var req CreateNoteRequest
	if !validation.DecodeAndValidate(w, r, &req) {
		return
	}

	// Verify ownership
	var ownerID string
	err = database.DB.QueryRow(r.Context(), "SELECT user_id FROM notes WHERE id=$1", noteID).Scan(&ownerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if ownerID != userUUID.String() {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec(r.Context(),
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
func GetNote(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID
	noteID := chi.URLParam(r, "id")

	userUUID, err := GetUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var n Note
	err = database.DB.QueryRow(r.Context(), "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id=$1", noteID).Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if n.UserID != userUUID.String() {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if err := json.NewEncoder(w).Encode(n); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
