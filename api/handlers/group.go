package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/models"
	"discipleship_journal_api/services"

	"firebase.google.com/go/v4/auth"
	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GroupHandler struct {
	service services.GroupServiceInterface
}

func NewGroupHandler(service services.GroupServiceInterface) *GroupHandler {
	return &GroupHandler{service: service}
}

type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=500"`
	Type        string `json:"type" validate:"omitempty,oneof=group direct"`
}

type GroupResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedBy   string  `json:"created_by"`
	Role        string  `json:"role,omitempty"` // Current user's role
	Type        string  `json:"type"`
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

	// Check for test override first to avoid DB lookup
	if testUserID := r.Context().Value(TestUserKey); testUserID != nil {
		if idStr, ok := testUserID.(string); ok {
			userUUID = idStr
		} else if id, ok := testUserID.(uuid.UUID); ok {
			userUUID = id.String()
		}
	}

	// If not found in test override, try standard auth
	if userUUID == "" {
		if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
			var id uuid.UUID
			id, err = GetUserUUID(r.Context(), token.UID)
			if err != nil {
				http.Error(w, "User not found", http.StatusInternalServerError)
				return
			}
			userUUID = id.String()
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	group, err := h.service.CreateGroup(r.Context(), userUUID, req.Name, &req.Description, req.Type)
	if err != nil {
		slog.Error("Failed to create group", "error", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": group.ID}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// ListMyGroups lists groups the user belongs to
func (h *GroupHandler) ListMyGroups(w http.ResponseWriter, r *http.Request) {
	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	slog.Info("ListMyGroups", "user", userUUID)

	groups, err := h.service.ListUserGroups(r.Context(), userUUID.String())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	slog.Info("ListMyGroups returning", "count", len(groups))

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

	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	groups, err := h.service.SearchGroups(r.Context(), query, userUUID.String())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// JoinGroup allows a user to join a group
func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	err = h.service.JoinGroup(r.Context(), groupID, userUUID.String())
	if err != nil {
		if err == models.ErrAlreadyExists {
			http.Error(w, "Already a member", http.StatusConflict)
		} else {
			http.Error(w, "Failed to join group", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// LeaveGroup allows a user to leave a group
func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	err = h.service.LeaveGroup(r.Context(), groupID, userUUID.String())
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Not a member", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// GetGroupMembers lists members of a group
func (h *GroupHandler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	members, err := h.service.GetGroupMembers(r.Context(), groupID, userUUID.String())
	if err != nil {
		if err.Error() == "access denied" {
			http.Error(w, "Access denied", http.StatusForbidden)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	var response []GroupMemberResponse
	for _, m := range members {
		response = append(response, GroupMemberResponse{
			UserID:      m.UserID,
			DisplayName: m.DisplayName,
			Email:       m.Email,
			Role:        m.Role,
			JoinedAt:    m.JoinedAt.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
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

	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	err = h.service.AddGroupMember(r.Context(), userUUID.String(), groupID, req.UserID)
	if err != nil {
		// Map errors
		if err.Error() == "admin rights required" || err.Error() == "access denied" || err.Error() == "user is not in your connections" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			slog.Error("Failed to add member", "error", err)
			http.Error(w, "Failed to add member", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// GetOrCreateDirectGroup gets or creates a direct message group with a connection
func (h *GroupHandler) GetOrCreateDirectGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PartnerID string `json:"partner_id" validate:"required"`
	}
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	group, created, err := h.service.GetOrCreateDirectGroup(r.Context(), userUUID.String(), req.PartnerID)
	if err != nil {
		if err.Error() == "cannot create direct group with yourself" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err.Error() == "user is not in your connections" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else if err == models.ErrNotFound { // Partner not found
			http.Error(w, "Partner not found", http.StatusNotFound)
		} else {
			slog.Error("Database error", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(map[string]string{"id": group.ID}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// RemoveGroupMember removes a user from a group (Admin only)
func (h *GroupHandler) RemoveGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	targetUserID := chi.URLParam(r, "userId")

	userUUID, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	err = h.service.RemoveGroupMember(r.Context(), userUUID.String(), groupID, targetUserID)
	if err != nil {
		if err.Error() == "admin rights required" || err.Error() == "access denied" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"success": true}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}
