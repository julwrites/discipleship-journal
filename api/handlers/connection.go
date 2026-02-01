package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"discipleship_journal_api/services"
	chi "github.com/go-chi/chi/v5"
	pgx "github.com/jackc/pgx/v5"
)

type ConnectionRequest struct {
	ReceiverEmail string `json:"receiver_email" validate:"omitempty,email"`
	ReceiverID    string `json:"receiver_id" validate:"omitempty,uuid"`
}

type ConnectionResponse struct {
	ID                string  `json:"id"`
	RequesterID       string  `json:"requester_id"`
	ReceiverID        string  `json:"receiver_id"`
	Status            string  `json:"status"`
	RequesterEmail    string  `json:"requester_email,omitempty"`
	ReceiverEmail     string  `json:"receiver_email,omitempty"`
	RequesterUsername *string `json:"requester_username,omitempty"`
	ReceiverUsername  *string `json:"receiver_username,omitempty"`
}

type ConnectionHandler struct {
	db                  DBInterface
	notificationService services.NotificationService
}

func NewConnectionHandler(db DBInterface, notificationService services.NotificationService) *ConnectionHandler {
	return &ConnectionHandler{db: db, notificationService: notificationService}
}

// SearchUsers searches for users by email or username
func (h *ConnectionHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if len(query) < 3 {
		http.Error(w, "Search query too short", http.StatusBadRequest)
		return
	}

	requesterUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Strict match on email or username
	rows, err := h.db.Query(r.Context(),
		`SELECT u.id, u.email, u.username,
		 EXISTS(SELECT 1 FROM connections c WHERE ((c.requester_id = u.id AND c.receiver_id = $2) OR (c.receiver_id = u.id AND c.requester_id = $2)) AND c.status = 'accepted') as is_connected
		 FROM users u
		 WHERE (email ILIKE $1 OR username ILIKE $1) AND u.id != $2 LIMIT 20`,
		query, requesterUUID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var id, email string
		var username *string
		var isConnected bool
		if err := rows.Scan(&id, &email, &username, &isConnected); err != nil {
			continue
		}

		user := map[string]interface{}{
			"id":           id,
			"is_connected": isConnected,
		}
		if username != nil && *username != "" {
			user["username"] = *username
		}

		// Reveal email if connected OR if the query was the email itself
		if isConnected || strings.EqualFold(email, query) {
			user["email"] = email
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

// SendConnectionRequest sends a connection request to another user
func (h *ConnectionHandler) SendConnectionRequest(w http.ResponseWriter, r *http.Request) {
	var req ConnectionRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	requesterUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if req.ReceiverEmail == "" && req.ReceiverID == "" {
		http.Error(w, "Receiver email or ID required", http.StatusBadRequest)
		return
	}

	// Find receiver UUID
	var receiverUUID string
	if req.ReceiverID != "" {
		// Verify exists
		err = h.db.QueryRow(r.Context(), "SELECT id FROM users WHERE id = $1", req.ReceiverID).Scan(&receiverUUID)
	} else {
		err = h.db.QueryRow(r.Context(), "SELECT id FROM users WHERE email = $1", req.ReceiverEmail).Scan(&receiverUUID)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if requesterUUID.String() == receiverUUID {
		http.Error(w, "Cannot connect with yourself", http.StatusBadRequest)
		return
	}

	// Insert connection
	var connID string
	err = h.db.QueryRow(r.Context(),
		`INSERT INTO connections (requester_id, receiver_id, status)
		 VALUES ($1, $2, 'pending')
		 RETURNING id`, requesterUUID, receiverUUID).Scan(&connID)

	if err != nil {
		// Handle unique constraint violation gracefully
		http.Error(w, "Connection request already exists or error", http.StatusConflict)
		return
	}

	// Fetch requester name synchronously to avoid race conditions in tests and ensure data availability
	var requesterName string
	var requesterUsername *string
	if err := h.db.QueryRow(r.Context(), "SELECT username FROM users WHERE id = $1", requesterUUID).Scan(&requesterUsername); err != nil {
		requesterName = "Someone"
	} else if requesterUsername != nil && *requesterUsername != "" {
		requesterName = *requesterUsername
	} else {
		requesterName = "Someone"
	}

	// Send notification
	go func() {
		// Use a detached context for background task
		ctx := context.Background()

		err := h.notificationService.SendNotification(ctx, receiverUUID, "New Connection Request", requesterName+" wants to connect with you.", map[string]string{
			"type": "connection_request",
			"id":   connID,
		})
		if err != nil {
			slog.Error("Failed to send notification", "error", err)
		}
	}()

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": connID, "status": "pending"}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// ListConnections lists all connections for the current user
func (h *ConnectionHandler) ListConnections(w http.ResponseWriter, r *http.Request) {
	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT c.id, c.requester_id, c.receiver_id, c.status,
		        u1.email as requester_email, u2.email as receiver_email,
		        u1.username as requester_username, u2.username as receiver_username
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
		if err := rows.Scan(&c.ID, &c.RequesterID, &c.ReceiverID, &c.Status, &c.RequesterEmail, &c.ReceiverEmail, &c.RequesterUsername, &c.ReceiverUsername); err != nil {
			continue
		}
		connections = append(connections, c)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(connections); err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

// AcceptConnectionRequest accepts a connection request
func (h *ConnectionHandler) AcceptConnectionRequest(w http.ResponseWriter, r *http.Request) {
	connID := chi.URLParam(r, "id")

	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commandTag, err := h.db.Exec(r.Context(),
		"UPDATE connections SET status = 'accepted' WHERE id = $1 AND receiver_id = $2", connID, userUUID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Connection request not found or not for you", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// DeleteConnectionRequest rejects or deletes a connection
func (h *ConnectionHandler) DeleteConnectionRequest(w http.ResponseWriter, r *http.Request) {
	connID := chi.URLParam(r, "id")

	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commandTag, err := h.db.Exec(r.Context(),
		"DELETE FROM connections WHERE id = $1 AND (receiver_id = $2 OR requester_id = $2)", connID, userUUID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Connection not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}
