package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type User struct {
	ID          string    `json:"id"`
	FirebaseUID string    `json:"firebase_uid"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserSettings struct {
	BibleVersion string `json:"bible_version"`
}

type UpdateUserRequest struct {
	FullName  string `json:"full_name" validate:"required,min=2,max=100"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

// CreateOrUpdateUser handles user creation (or upsert)
// @Summary Create or update a user
// @Description Register a new user or update existing one based on Firebase UID
// @Tags users
// @Accept json
// @Produce json
// @Param user body UpdateUserRequest true "User details"
// @Success 200 {object} User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func CreateOrUpdateUser(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	// Implementation placeholder
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(User{
		ID:        "123",
		Email:     "test@example.com",
		FullName:  req.FullName,
		AvatarURL: req.AvatarURL,
	}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// UpdateUser handles user updates
// @Summary Update a user
// @Description Update existing user details
// @Tags users
// @Accept json
// @Produce json
// @Param user body UpdateUserRequest true "User details"
// @Success 200 {object} User
// @Failure 400 {object} map[string]string
// @Router /users [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	// Implementation placeholder
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(User{
		ID:        "123",
		Email:     "test@example.com",
		FullName:  req.FullName,
		AvatarURL: req.AvatarURL,
	}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

// GetMe returns the current user profile
// @Summary Get current user
// @Description Get the profile of the authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} User
// @Router /users/me [get]
func GetMe(w http.ResponseWriter, r *http.Request) {
	// Implementation placeholder
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(User{
		ID:    "123",
		Email: "test@example.com",
	}); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}
