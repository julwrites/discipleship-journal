package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"discipleship_journal_api/database"
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
	Title   string                 `json:"title"`
	Content map[string]interface{} `json:"content"`
}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("user").(*auth.Token)
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

	json.NewEncoder(w).Encode(notes)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("user").(*auth.Token)
	uid := token.UID

	userUUID, err := GetUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
	json.NewEncoder(w).Encode(map[string]string{"id": noteID})
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("user").(*auth.Token)
	uid := token.UID
	noteID := chi.URLParam(r, "id")

	userUUID, err := GetUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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

	if ownerID != userUUID {
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

func GetNote(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("user").(*auth.Token)
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

	if n.UserID != userUUID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(n)
}
