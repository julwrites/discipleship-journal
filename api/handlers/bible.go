package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

type PassageRequest struct {
	Reference string `json:"reference"`
}

type PassageResponse struct {
	Reference string `json:"reference"`
	Text      string `json:"text"`
	// Add other fields from BibleAIAPI response if needed
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
func GetBiblePassage(w http.ResponseWriter, r *http.Request) {
	// BibleAIAPI Endpoint
	apiURL := os.Getenv("BIBLE_API_URL")
	apiKey := os.Getenv("BIBLE_API_KEY")

	if apiURL == "" {
		http.Error(w, "Bible API not configured", http.StatusInternalServerError)
		return
	}

	// Read query param
	ref := r.URL.Query().Get("ref")
	if ref == "" {
		http.Error(w, "Reference is required", http.StatusBadRequest)
		return
	}

	client := resty.New()

	// Assuming BibleAIAPI structure. If it's a GET request:
	// Adjust endpoint path based on actual BibleAIAPI docs
	var result map[string]interface{}
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey). // Or whatever auth method it uses
		SetQueryParam("q", ref). // Or "reference" or path param
		SetResult(&result).
		Get(apiURL + "/bible/passage") // Adjust path

	if err != nil {
		http.Error(w, "Failed to call Bible API", http.StatusInternalServerError)
		return
	}

	if resp.IsError() {
		http.Error(w, fmt.Sprintf("Bible API Error: %s", resp.Status()), resp.StatusCode())
		return
	}

	// Transform result if necessary, or just proxy it
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
