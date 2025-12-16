package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"discipleship_journal_api/services"
)

type BibleHandler struct {
	Client services.BibleAIClient
}

func NewBibleHandler(client services.BibleAIClient) *BibleHandler {
	return &BibleHandler{Client: client}
}

// GetBiblePassage godoc
// @Summary Get Bible Passage
// @Description Fetch a bible passage from the external API
// @Tags bible
// @Accept json
// @Produce json
// @Param ref query string true "Bible Reference (e.g. John 3:16)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Reference is required"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/bible/passage [get]
func (h *BibleHandler) GetBiblePassage(w http.ResponseWriter, r *http.Request) {
	// Read query param
	ref := r.URL.Query().Get("ref")
	if ref == "" {
		http.Error(w, "Reference is required", http.StatusBadRequest)
		return
	}

	result, err := h.Client.GetPassage(r.Context(), ref)

	if err != nil {
		// Log error if logger is available, or just send 500
		http.Error(w, fmt.Sprintf("Failed to call Bible API: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
