package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	chi "github.com/go-chi/chi/v5"
)

type GroupHandler struct {
	service services.GroupService
}

func NewGroupHandler(service services.GroupService) *GroupHandler {
	return &GroupHandler{service: service}
}

type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// CreateGroup creates a new group
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req CreateGroupRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	groupID, err := h.service.CreateGroup(r.Context(), req.Name, req.Description, userUUID.String())
	if err != nil {
		slog.Error("Failed to create group", "error", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
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

	groups, err := h.service.ListMyGroups(r.Context(), userUUID.String())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
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

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
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
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	err = h.service.JoinGroup(r.Context(), groupID, userUUID.String())
	if err != nil {
		if errors.Is(err, services.ErrAlreadyMember) {
			http.Error(w, "Already a member", http.StatusConflict)
		} else {
			http.Error(w, "Failed to join group", http.StatusInternalServerError)
		}
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

	err = h.service.LeaveGroup(r.Context(), groupID, userUUID.String())
	if err != nil {
		if errors.Is(err, services.ErrNotMember) {
			http.Error(w, "Not a member", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetGroupMembers lists members of a group
func (h *GroupHandler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	members, err := h.service.GetGroupMembers(r.Context(), groupID, userUUID.String())
	if err != nil {
		if errors.Is(err, services.ErrAccessDenied) {
			http.Error(w, "Access denied", http.StatusForbidden)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
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

	err = h.service.AddGroupMember(r.Context(), groupID, req.UserID, userUUID.String())
	if err != nil {
		if errors.Is(err, services.ErrAccessDenied) || errors.Is(err, services.ErrInsufficientPerms) {
			http.Error(w, "Access denied", http.StatusForbidden)
		} else {
			slog.Error("Failed to add member", "error", err)
			http.Error(w, "Failed to add member", http.StatusInternalServerError)
		}
		return
	}

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

	err = h.service.RemoveGroupMember(r.Context(), groupID, targetUserID, userUUID.String())
	if err != nil {
		if errors.Is(err, services.ErrAccessDenied) || errors.Is(err, services.ErrInsufficientPerms) {
			http.Error(w, "Access denied", http.StatusForbidden)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
