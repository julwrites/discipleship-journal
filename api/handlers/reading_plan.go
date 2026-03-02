package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"discipleship_journal_api/models"
	"discipleship_journal_api/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ReadingPlanHandler struct {
	service services.ReadingPlanService
}

func NewReadingPlanHandler(service services.ReadingPlanService) *ReadingPlanHandler {
	return &ReadingPlanHandler{service: service}
}

// GetAllPlans godoc
// @Summary      Get all reading plans
// @Description  Get a list of all available reading plans
// @Tags         reading-plans
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api/reading-plans [get]
func (h *ReadingPlanHandler) GetAllPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.service.GetAllPlans(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"data": plans,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Unsubscribe godoc
// @Summary      Unsubscribe from a reading plan
// @Description  Unsubscribe the current user from a reading plan, deleting their progress
// @Tags         reading-plans
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Success      204  {object}  nil
// @Router       /api/reading-plans/{id}/subscribe [delete]
func (h *ReadingPlanHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	planID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	err = h.service.Unsubscribe(r.Context(), userID, planID)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Subscription not found or not active", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnmarkDayComplete godoc
// @Summary      Unmark a day as complete
// @Description  Unmark (undo completion) a specific day in a reading plan
// @Tags         reading-plans
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Param        day_number path int true "Day Number"
// @Success      200  {object}  map[string]bool
// @Router       /api/my-reading-plans/{id}/progress/{day_number} [delete]
func (h *ReadingPlanHandler) UnmarkDayComplete(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	planID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	dayStr := chi.URLParam(r, "day_number")
	dayNumber, err := strconv.Atoi(dayStr)
	if err != nil {
		http.Error(w, "Invalid day number", http.StatusBadRequest)
		return
	}

	err = h.service.UnmarkDayComplete(r.Context(), userID, planID, dayNumber)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "User plan not found or not active", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetPlanProgress godoc
// @Summary      Get reading plan progress
// @Description  Get the list of completed days for a specific plan
// @Tags         reading-plans
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /api/my-reading-plans/{id}/progress [get]
func (h *ReadingPlanHandler) GetPlanProgress(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	planID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	days, err := h.service.GetPlanProgress(r.Context(), userID, planID)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Plan not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if days == nil {
		days = []int{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"completed_days": days,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetPlan godoc
// @Summary      Get a reading plan
// @Description  Get details of a specific reading plan including days
// @Tags         reading-plans
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Success      200  {object}  models.ReadingPlan
// @Router       /api/reading-plans/{id} [get]
func (h *ReadingPlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	plan, err := h.service.GetPlan(r.Context(), id)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Plan not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	days, err := h.service.GetPlanDays(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Combine into a response struct
	response := struct {
		ID          uuid.UUID                `json:"id"`
		Title       string                   `json:"title"`
		Description string                   `json:"description"`
		TotalDays   int                      `json:"total_days"`
		PlanType    string                   `json:"plan_type"`
		CreatedAt   interface{}              `json:"created_at"`
		UpdatedAt   interface{}              `json:"updated_at"`
		Days        []*models.ReadingPlanDay `json:"days"`
	}{
		ID:          plan.ID,
		Title:       plan.Title,
		Description: plan.Description,
		TotalDays:   plan.Days,
		PlanType:    plan.PlanType,
		CreatedAt:   plan.CreatedAt,
		UpdatedAt:   plan.UpdatedAt,
		Days:        days,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Helper to get userID from context (production) or test fallback
func (h *ReadingPlanHandler) getUserID(r *http.Request) (uuid.UUID, error) {
	id, err := GetUserUUIDFromContext(r.Context())
	if err != nil {
		return uuid.Nil, models.ErrNotFound
	}
	return id, nil
}

// Subscribe godoc
// @Summary      Subscribe to a reading plan
// @Description  Subscribe the current user to a reading plan
// @Tags         reading-plans
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Success      200  {object}  models.UserReadingPlan
// @Router       /api/reading-plans/{id}/subscribe [post]
func (h *ReadingPlanHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	planID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	userPlan, err := h.service.Subscribe(r.Context(), userID, planID)
	if err != nil {
		if err == models.ErrAlreadyExists {
			http.Error(w, "Already subscribed", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(userPlan); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetUserPlans godoc
// @Summary      Get user's reading plans
// @Description  Get a list of reading plans the user is subscribed to
// @Tags         reading-plans
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api/my-reading-plans [get]
func (h *ReadingPlanHandler) GetUserPlans(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	plans, err := h.service.GetUserPlans(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"data": plans,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// MarkDayComplete godoc
// @Summary      Mark a day as complete
// @Description  Mark a specific day in a reading plan as complete
// @Tags         reading-plans
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Param        day_number body int true "Day Number"
// @Success      200  {object}  map[string]bool
// @Router       /api/my-reading-plans/{id}/progress [post]
func (h *ReadingPlanHandler) MarkDayComplete(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	planID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	var req struct {
		DayNumber int `json:"day_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.MarkDayComplete(r.Context(), userID, planID, req.DayNumber)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "User plan not found or not active", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
