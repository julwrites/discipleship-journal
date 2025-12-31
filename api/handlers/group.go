package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	chi "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
)

type GroupHandler struct {
	db                  DBInterface
	notificationService services.NotificationService
}

func NewGroupHandler(db DBInterface, notificationService services.NotificationService) *GroupHandler {
	return &GroupHandler{db: db, notificationService: notificationService}
}

type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type GroupResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
	Role        string `json:"role,omitempty"` // Current user's role
}

type GroupMemberResponse struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	JoinedAt    string `json:"joined_at"`
}

// CreateGroup creates a new group
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req CreateGroupRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

    var userUUID string
    var err error
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
        var id uuid.UUID
		id, err = GetUserUUID(r.Context(), token.UID)
        if err != nil {
		    http.Error(w, "User not found", http.StatusInternalServerError)
		    return
	    }
        userUUID = id.String()
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		userUUID = testUserID
	} else {
         http.Error(w, "Unauthorized", http.StatusUnauthorized)
         return
    }

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := tx.Rollback(r.Context()); err != nil && err != pgx.ErrTxClosed {
			slog.Error("Failed to rollback transaction", "error", err)
		}
	}()

	var groupID string
	err = tx.QueryRow(r.Context(),
		"INSERT INTO groups (name, description, created_by) VALUES ($1, $2, $3) RETURNING id",
		req.Name, req.Description, userUUID).Scan(&groupID)
	if err != nil {
		slog.Error("Failed to create group", "error", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	// Add creator as admin
	_, err = tx.Exec(r.Context(),
		"INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'admin')",
		groupID, userUUID)
	if err != nil {
		slog.Error("Failed to add member", "error", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": groupID}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// ListMyGroups lists groups the user belongs to
func (h *GroupHandler) ListMyGroups(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT g.id, g.name, g.description, g.created_by, gm.role
		 FROM groups g
		 JOIN group_members gm ON g.id = gm.group_id
		 WHERE gm.user_id = $1`, userUUID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var groups []GroupResponse
	for rows.Next() {
		var g GroupResponse
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.CreatedBy, &g.Role); err != nil {
			continue
		}
		groups = append(groups, g)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// SearchGroups searches for groups by name
func (h *GroupHandler) SearchGroups(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if len(query) < 3 {
		http.Error(w, "Search query too short", http.StatusBadRequest)
		return
	}

	// List groups and check if current user is a member
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT g.id, g.name, g.description, g.created_by,
		 COALESCE((SELECT role FROM group_members WHERE group_id = g.id AND user_id = $2), '') as role
		 FROM groups g
		 WHERE g.name ILIKE $1 LIMIT 20`, "%"+query+"%", userUUID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var groups []GroupResponse
	for rows.Next() {
		var g GroupResponse
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.CreatedBy, &g.Role); err != nil {
			continue
		}
		groups = append(groups, g)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// JoinGroup allows a user to join a group
func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Check if already a member
	var exists bool
	err = h.db.QueryRow(r.Context(),
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userUUID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Already a member", http.StatusConflict)
		return
	}

	_, err = h.db.Exec(r.Context(),
		"INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'member')",
		groupID, userUUID)
	if err != nil {
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// LeaveGroup allows a user to leave a group
func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Prevent last admin from leaving? (Optional, skipping for MVP)

	result, err := h.db.Exec(r.Context(),
		"DELETE FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, userUUID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Not a member", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetGroupMembers lists members of a group
func (h *GroupHandler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")

	// Check if user is a member of the group (or group is public? Assuming public read of members for now)
	// For privacy, maybe only members can see members. Let's enforce membership.
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	var isMember bool
	err = h.db.QueryRow(r.Context(),
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

	rows, err := h.db.Query(r.Context(),
		`SELECT gm.user_id, u.display_name, u.email, gm.role, gm.joined_at
		 FROM group_members gm
		 JOIN users u ON gm.user_id = u.id
		 WHERE gm.group_id = $1`, groupID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []GroupMemberResponse
	for rows.Next() {
		var m GroupMemberResponse
		var joinedAt time.Time
		if err := rows.Scan(&m.UserID, &m.DisplayName, &m.Email, &m.Role, &joinedAt); err != nil {
			continue
		}
		m.JoinedAt = joinedAt.Format(time.RFC3339)
		members = append(members, m)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(members); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

type AddMemberRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// AddGroupMember adds a user to a group (Admin only)
func (h *GroupHandler) AddGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	var req AddMemberRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Verify admin role
	var role string
	err = h.db.QueryRow(r.Context(),
		"SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, userUUID).Scan(&role)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Access denied", http.StatusForbidden)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if role != "admin" {
		http.Error(w, "Admin rights required", http.StatusForbidden)
		return
	}

	// Check if user to add exists and is not already member
	// (Assuming req.UserID is the internal UUID, we should probably check existence)
	_, err = h.db.Exec(r.Context(),
		"INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'member') ON CONFLICT DO NOTHING",
		groupID, req.UserID)
	if err != nil {
		slog.Error("Failed to add member", "error", err)
		http.Error(w, "Failed to add member", http.StatusInternalServerError)
		return
	}

	// Get group name synchronously
	var groupName string
	if err := h.db.QueryRow(r.Context(), "SELECT name FROM groups WHERE id = $1", groupID).Scan(&groupName); err != nil {
		groupName = "a group"
	}

	// Send notification
	go func() {
		ctx := context.Background()

		err := h.notificationService.SendNotification(ctx, req.UserID, "Group Invitation", "You have been added to "+groupName, map[string]string{
			"type": "group_invite",
			"id":   groupID,
		})
		if err != nil {
			slog.Error("Failed to send notification", "error", err)
		}
	}()

	w.WriteHeader(http.StatusCreated)
}

// RemoveGroupMember removes a user from a group (Admin only)
func (h *GroupHandler) RemoveGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	targetUserID := chi.URLParam(r, "userId")

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Verify admin role
	var role string
	err = h.db.QueryRow(r.Context(),
		"SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, userUUID).Scan(&role)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden) // Or DB error
		return
	}

	if role != "admin" {
		http.Error(w, "Admin rights required", http.StatusForbidden)
		return
	}

	// Prevent removing self via this endpoint? Or allow it?
	// Usually admins leave via LeaveGroup, but removing self is edge case.
	// If admin removes self, who is admin?
	// For now, let's allow removing any member.

	_, err = h.db.Exec(r.Context(),
		"DELETE FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, targetUserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
