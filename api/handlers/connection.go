package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"discipleship_journal_api/database"
	"discipleship_journal_api/middleware"
	"discipleship_journal_api/validation"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type ConnectionRequest struct {
	ReceiverEmail string `json:"receiver_email" validate:"required,email"`
}

type ConnectionResponse struct {
	ID             string `json:"id"`
	RequesterID    string `json:"requester_id"`
	ReceiverID     string `json:"receiver_id"`
	Status         string `json:"status"`
	RequesterEmail string `json:"requester_email,omitempty"`
	ReceiverEmail  string `json:"receiver_email,omitempty"`
}

// SearchUsers searches for users by email or username
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if len(query) < 3 {
		http.Error(w, "Search query too short", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(r.Context(),
		"SELECT id, email, display_name FROM users WHERE email ILIKE $1 OR display_name ILIKE $1 LIMIT 10",
		"%"+query+"%")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]string
	for rows.Next() {
		var id, email, displayName string
		if err := rows.Scan(&id, &email, &displayName); err != nil {
			continue
		}
		users = append(users, map[string]string{
			"id":           id,
			"email":        email,
			"display_name": displayName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

// SendConnectionRequest sends a connection request to another user
func SendConnectionRequest(w http.ResponseWriter, r *http.Request) {
	var req ConnectionRequest
	if !validation.DecodeAndValidate(w, r, &req) {
		return
	}

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	requesterUID := token.UID
	requesterUUID, err := GetUserUUID(r.Context(), requesterUID)
	if err != nil {
		http.Error(w, "Requester not found", http.StatusInternalServerError)
		return
	}

	// Find receiver UUID
	var receiverUUID string
	err = database.DB.QueryRow(r.Context(), "SELECT id FROM users WHERE email = $1", req.ReceiverEmail).Scan(&receiverUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if requesterUUID == receiverUUID {
		http.Error(w, "Cannot connect with yourself", http.StatusBadRequest)
		return
	}

	// Insert connection
	var connID string
	err = database.DB.QueryRow(r.Context(),
		`INSERT INTO connections (requester_id, receiver_id, status)
		 VALUES ($1, $2, 'pending')
		 RETURNING id`, requesterUUID, receiverUUID).Scan(&connID)

	if err != nil {
		// Handle unique constraint violation gracefully
		http.Error(w, "Connection request already exists or error", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": connID, "status": "pending"}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// ListConnections lists all connections for the current user
func ListConnections(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID
	userUUID, err := GetUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	rows, err := database.DB.Query(r.Context(),
		`SELECT c.id, c.requester_id, c.receiver_id, c.status,
		        u1.email as requester_email, u2.email as receiver_email
		 FROM connections c
		 JOIN users u1 ON c.requester_id = u1.id
		 JOIN users u2 ON c.receiver_id = u2.id
		 WHERE c.requester_id = $1 OR c.receiver_id = $1`, userUUID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var connections []ConnectionResponse
	for rows.Next() {
		var c ConnectionResponse
		if err := rows.Scan(&c.ID, &c.RequesterID, &c.ReceiverID, &c.Status, &c.RequesterEmail, &c.ReceiverEmail); err != nil {
			continue
		}
		connections = append(connections, c)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(connections); err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

// RespondToConnectionRequest accepts or rejects a connection request
func RespondToConnectionRequest(w http.ResponseWriter, r *http.Request) {
	connID := chi.URLParam(r, "id")
	action := r.URL.Query().Get("action") // accept or reject

	if action != "accept" && action != "reject" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID
	userUUID, err := GetUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	if action == "reject" {
		_, err = database.DB.Exec(r.Context(),
			"DELETE FROM connections WHERE id = $1 AND (receiver_id = $2 OR requester_id = $2)", connID, userUUID)
	} else {
		_, err = database.DB.Exec(r.Context(),
			"UPDATE connections SET status = 'accepted' WHERE id = $1 AND receiver_id = $2", connID, userUUID)
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
