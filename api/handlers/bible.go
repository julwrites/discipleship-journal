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
// @Param version query string false "Bible Version (e.g. ESV)"
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
	version := r.URL.Query().Get("version")

	result, err := h.Client.GetPassage(r.Context(), ref, version)

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

// GetBibleVersions godoc
// @Summary Get Bible Versions
// @Description Fetch available bible versions from the external API
// @Tags bible
// @Accept json
// @Produce json
// @Param name query string false "Filter by name"
// @Param language query string false "Filter by language"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/bible/versions [get]
func (h *BibleHandler) GetBibleVersions(w http.ResponseWriter, r *http.Request) {
	params := make(map[string]string)
	q := r.URL.Query()
	for k, v := range q {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	result, err := h.Client.GetVersions(r.Context(), params)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch versions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
