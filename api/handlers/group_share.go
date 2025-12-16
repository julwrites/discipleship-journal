package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"discipleship_journal_api/database"
	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type ShareNoteRequest struct {
	NoteID  string `json:"note_id" validate:"required"`
	Comment string `json:"comment" validate:"max=500"`
}

type SharedNoteResponse struct {
	ID       string                 `json:"id"` // This is the share ID
	GroupID  string                 `json:"group_id"`
	NoteID   string                 `json:"note_id"`
	Title    string                 `json:"title"`
	Content  map[string]interface{} `json:"content"`   // Included for detail view
	SharedBy string                 `json:"shared_by"` // User display name
	SharedAt string                 `json:"shared_at"`
	Comment  string                 `json:"comment"`
}

// ShareNoteToGroup shares a note to a group
func ShareNoteToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	var req ShareNoteRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// 1. Verify membership in group
	var isMember bool
	err = database.DB.QueryRow(r.Context(),
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userUUID).Scan(&isMember)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// 2. Verify ownership of note
	var ownerID string
	err = database.DB.QueryRow(r.Context(),
		"SELECT user_id FROM notes WHERE id = $1",
		req.NoteID).Scan(&ownerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}
	if ownerID != userUUID.String() {
		http.Error(w, "You can only share your own notes", http.StatusForbidden)
		return
	}

	// 3. Create share
	_, err = database.DB.Exec(r.Context(),
		"INSERT INTO group_shares (group_id, note_id, shared_by, comment) VALUES ($1, $2, $3, $4) ON CONFLICT (group_id, note_id) DO UPDATE SET shared_at = NOW(), comment = $4",
		groupID, req.NoteID, userUUID, req.Comment)
	if err != nil {
		slog.Error("Failed to share note", "error", err)
		http.Error(w, "Failed to share note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ListGroupShares lists notes shared with the group
func ListGroupShares(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Verify membership
	var isMember bool
	err = database.DB.QueryRow(r.Context(),
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userUUID).Scan(&isMember)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	rows, err := database.DB.Query(r.Context(),
		`SELECT gs.id, gs.group_id, gs.note_id, n.title, u.display_name, gs.shared_at, gs.comment
		 FROM group_shares gs
		 JOIN notes n ON gs.note_id = n.id
		 JOIN users u ON gs.shared_by = u.id
		 WHERE gs.group_id = $1
		 ORDER BY gs.shared_at DESC`, groupID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var shares []SharedNoteResponse
	for rows.Next() {
		var s SharedNoteResponse
		var sharedAt time.Time
		if err := rows.Scan(&s.ID, &s.GroupID, &s.NoteID, &s.Title, &s.SharedBy, &sharedAt, &s.Comment); err != nil {
			continue
		}
		s.SharedAt = sharedAt.Format(time.RFC3339)
		shares = append(shares, s)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(shares); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// GetSharedNoteDetails fetches details of a shared note
func GetSharedNoteDetails(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	shareID := chi.URLParam(r, "shareId")

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Verify membership
	var isMember bool
	err = database.DB.QueryRow(r.Context(),
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userUUID).Scan(&isMember)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var s SharedNoteResponse
	var sharedAt time.Time
	err = database.DB.QueryRow(r.Context(),
		`SELECT gs.id, gs.group_id, gs.note_id, n.title, n.content, u.display_name, gs.shared_at, gs.comment
		 FROM group_shares gs
		 JOIN notes n ON gs.note_id = n.id
		 JOIN users u ON gs.shared_by = u.id
		 WHERE gs.id = $1 AND gs.group_id = $2`, shareID, groupID).Scan(
		&s.ID, &s.GroupID, &s.NoteID, &s.Title, &s.Content, &s.SharedBy, &sharedAt, &s.Comment)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Shared note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}
	s.SharedAt = sharedAt.Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}
