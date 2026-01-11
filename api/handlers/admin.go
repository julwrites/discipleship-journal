package handlers

import (
	"encoding/json"
	"net/http"

	"discipleship_journal_api/services"
)

// AdminHandler handles admin operations.
type AdminHandler struct {
	bibleVersionService *services.BibleVersionService
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(bibleVersionService *services.BibleVersionService) *AdminHandler {
	return &AdminHandler{bibleVersionService: bibleVersionService}
}

// SyncBibleVersions handles the syncing of bible versions.
// @Summary Sync Bible Versions
// @Description Scrapes and syncs bible versions from external source.
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /admin/sync-bible-versions [post]
func (h *AdminHandler) SyncBibleVersions(w http.ResponseWriter, r *http.Request) {
	err := h.bibleVersionService.SyncVersions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Bible versions synced successfully",
	}); err != nil {
		// Log error if logger is available
		_ = err
	}
}
