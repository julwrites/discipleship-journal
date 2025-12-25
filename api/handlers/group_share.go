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

type GroupShareHandler struct {
	service services.GroupService
}

func NewGroupShareHandler(service services.GroupService) *GroupShareHandler {
	return &GroupShareHandler{service: service}
}

type ShareNoteRequest struct {
	NoteID  string `json:"note_id" validate:"required"`
	Comment string `json:"comment" validate:"max=500"`
}

// ShareNoteToGroup shares a note to a group
func (h *GroupShareHandler) ShareNoteToGroup(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.ShareNoteToGroup(r.Context(), groupID, req.NoteID, userUUID.String(), req.Comment)
	if err != nil {
		if errors.Is(err, services.ErrAccessDenied) || errors.Is(err, services.ErrNoteNotOwned) {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else if errors.Is(err, services.ErrNoteNotFound) {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			slog.Error("Failed to share note", "error", err)
			http.Error(w, "Failed to share note", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ListGroupShares lists notes shared with the group
func (h *GroupShareHandler) ListGroupShares(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	shares, err := h.service.ListGroupShares(r.Context(), groupID, userUUID.String())
	if err != nil {
		if errors.Is(err, services.ErrAccessDenied) {
			http.Error(w, "Access denied", http.StatusForbidden)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(shares); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// GetSharedNoteDetails fetches details of a shared note
func (h *GroupShareHandler) GetSharedNoteDetails(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	shareID := chi.URLParam(r, "shareId")

	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userUUID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	share, err := h.service.GetSharedNoteDetails(r.Context(), shareID, groupID, userUUID.String())
	if err != nil {
		if errors.Is(err, services.ErrAccessDenied) {
			http.Error(w, "Access denied", http.StatusForbidden)
		} else if errors.Is(err, services.ErrShareNotFound) {
			http.Error(w, "Shared note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(share); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}
