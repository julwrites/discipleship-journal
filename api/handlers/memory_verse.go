package handlers

import (
	"encoding/json"
	"net/http"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/models"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MemoryVerseHandler struct {
	service services.MemoryVerseService
}

func NewMemoryVerseHandler(service services.MemoryVerseService) *MemoryVerseHandler {
	return &MemoryVerseHandler{service: service}
}

// GetPacks lists packs (system or user)
func (h *MemoryVerseHandler) GetPacks(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	filterType := r.URL.Query().Get("type") // "system" or "user"
	if filterType == "" {
		filterType = "user"
	}

	packs, err := h.service.GetPacks(r.Context(), userID, filterType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": packs}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreatePack creates a new user pack
func (h *MemoryVerseHandler) CreatePack(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.VersePack
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	req.UserID = &userID
	req.IsPublic = false // User packs are private by default

	pack, err := h.service.CreatePack(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pack); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetPackDetails gets verses in a pack
func (h *MemoryVerseHandler) GetPackDetails(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	packIDStr := chi.URLParam(r, "id")
	packID, err := uuid.Parse(packIDStr)
	if err != nil {
		http.Error(w, "Invalid pack ID", http.StatusBadRequest)
		return
	}

	// 1. Get Pack Info
	pack, err := h.service.GetPack(r.Context(), packID, userID)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Pack not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 2. Get Verses
	verses, err := h.service.GetVerses(r.Context(), packID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"pack":   pack,
		"verses": verses,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateVerseInPack adds a verse to a pack
func (h *MemoryVerseHandler) CreateVerseInPack(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	packIDStr := chi.URLParam(r, "id")
	packID, err := uuid.Parse(packIDStr)
	if err != nil {
		http.Error(w, "Invalid pack ID", http.StatusBadRequest)
		return
	}

	// Verify ownership
	pack, err := h.service.GetPack(r.Context(), packID, userID)
	if err != nil {
		http.Error(w, "Pack not found", http.StatusNotFound)
		return
	}
	if pack.UserID == nil || *pack.UserID != userID {
		http.Error(w, "Cannot add verses to this pack", http.StatusForbidden)
		return
	}

	var req models.MemoryVerse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	req.VersePackID = packID

	verse, err := h.service.CreateVerse(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(verse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ClonePack clones a pack to user library
func (h *MemoryVerseHandler) ClonePack(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	packIDStr := chi.URLParam(r, "id")
	packID, err := uuid.Parse(packIDStr)
	if err != nil {
		http.Error(w, "Invalid pack ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Title string `json:"title"`
	}
	// Ignore error if body is empty or malformed, default values will be used
	_ = json.NewDecoder(r.Body).Decode(&req)

	pack, err := h.service.ClonePack(r.Context(), packID, userID, req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pack); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *MemoryVerseHandler) DeletePack(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	packIDStr := chi.URLParam(r, "id")
	packID, err := uuid.Parse(packIDStr)
	if err != nil {
		http.Error(w, "Invalid pack ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeletePack(r.Context(), packID, userID)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Pack not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// SearchVerses (for backward compatibility and global search)
func (h *MemoryVerseHandler) SearchVerses(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("q")

	verses, err := h.service.SearchVerses(r.Context(), userID, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": verses}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateVerse updates a verse
func (h *MemoryVerseHandler) UpdateVerse(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	verseIDStr := chi.URLParam(r, "verseId")
	verseID, err := uuid.Parse(verseIDStr)
	if err != nil {
		http.Error(w, "Invalid verse ID", http.StatusBadRequest)
		return
	}

	var req models.MemoryVerse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	req.ID = verseID

	err = h.service.UpdateVerse(r.Context(), &req, userID)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Verse not found or unauthorized", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// DeleteVerse deletes a verse
func (h *MemoryVerseHandler) DeleteVerse(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	verseIDStr := chi.URLParam(r, "verseId")
	verseID, err := uuid.Parse(verseIDStr)
	if err != nil {
		http.Error(w, "Invalid verse ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteVerse(r.Context(), verseID, userID)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Verse not found or unauthorized", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
