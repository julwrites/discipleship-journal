package handlers

import (
	"encoding/json"
	"net/http"

	"discipleship_journal_api/services"
)

// BibleVersionHandler handles bible version operations.
type BibleVersionHandler struct {
	service *services.BibleVersionService
}

// NewBibleVersionHandler creates a new BibleVersionHandler.
func NewBibleVersionHandler(service *services.BibleVersionService) *BibleVersionHandler {
	return &BibleVersionHandler{service: service}
}

// GetVersions returns a list of all bible versions.
// @Summary Get Bible Versions
// @Description Get a list of all available bible versions.
// @Tags bible-versions
// @Produce json
// @Success 200 {array} services.BibleVersion
// @Failure 500 {object} map[string]string
// @Router /bible-versions [get]
func (h *BibleVersionHandler) GetVersions(w http.ResponseWriter, r *http.Request) {
	versions, err := h.service.GetVersions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		// Log error if logger is available
		_ = err
	}
}
