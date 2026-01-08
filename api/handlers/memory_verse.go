package handlers

import (
	"encoding/json"
	"net/http"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/models"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
)

type MemoryVerseHandler struct {
	service services.MemoryVerseService
}

func NewMemoryVerseHandler(service services.MemoryVerseService) *MemoryVerseHandler {
	return &MemoryVerseHandler{service: service}
}

func (h *MemoryVerseHandler) SearchVerses(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("q")

	verses, err := h.service.SearchVerses(r.Context(), userID, query, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": verses})
}

func (h *MemoryVerseHandler) CreateVerse(w http.ResponseWriter, r *http.Request) {
	firebaseUID := ""
	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		firebaseUID = token.UID
	}

	userID, err := GetUserUUID(r.Context(), firebaseUID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.MemoryVerse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	req.UserID = &userID
	if req.PackName == "" {
		req.PackName = "My Verses"
	}

	verse, err := h.service.CreateVerse(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(verse)
}

func (h *MemoryVerseHandler) GetSystemPacks(w http.ResponseWriter, r *http.Request) {
	packs, err := h.service.GetSystemPacks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": packs})
}
