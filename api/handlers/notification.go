package handlers

import (
	"encoding/json"
	"net/http"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"

	"firebase.google.com/go/v4/auth"
)

type NotificationHandler struct {
	service services.NotificationService
}

func NewNotificationHandler(service services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

type RegisterDeviceRequest struct {
	Token      string `json:"token" validate:"required"`
	DeviceType string `json:"device_type"`
}

// RegisterDevice registers an FCM token for the authenticated user.
// @Summary Register FCM token
// @Description Registers a Firebase Cloud Messaging token for the user to receive push notifications.
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body RegisterDeviceRequest true "Device Registration"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/notifications/register [post]
// @Security BearerAuth
func (h *NotificationHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	userID, err := GetUserUUID(r.Context(), token.UID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req RegisterDeviceRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	if err := h.service.RegisterDevice(r.Context(), userID.String(), req.Token, req.DeviceType); err != nil {
		http.Error(w, "Failed to register device", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
}
