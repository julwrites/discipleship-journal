package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
)

type User struct {
	ID          string    `json:"id"`
	FirebaseUID string    `json:"firebase_uid"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Settings    string    `json:"settings"`
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

type UserHandler struct {
	db DBInterface
}

func NewUserHandler(db DBInterface) *UserHandler {
	return &UserHandler{db: db}
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
func (h *UserHandler) CreateOrUpdateUser(w http.ResponseWriter, r *http.Request) {
	var uid string
	token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	uid = token.UID

	var req UpdateUserRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	var user User
	// Upsert logic using ON CONFLICT
	// Note: firebase_uid should have a unique constraint
	// For now, just get existing user or return error
	// TODO: Update to use correct schema (username, settings instead of full_name, avatar_url)
	err := h.db.QueryRow(r.Context(), `
		SELECT id, firebase_uid, email, username, settings, created_at, updated_at
		FROM users
		WHERE firebase_uid=$1`, uid).Scan(
		&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &user.Settings, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		slog.Error("Failed to create/update user", "error", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
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
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// This is essentially the same as CreateOrUpdateUser but we might want to enforce existence.
	// Reusing logic for now, or we can make it stricter.
	h.CreateOrUpdateUser(w, r)
}

// GetMe returns the current user profile
// @Summary Get current user
// @Description Get the profile of the authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} User
// @Router /users/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	var user User
	var err error

	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		// Real Firebase auth - query by firebase_uid
		uid := token.UID
		err = h.db.QueryRow(r.Context(), `
			SELECT id, firebase_uid, email, username, settings, created_at, updated_at
			FROM users
			WHERE firebase_uid=$1`, uid).Scan(
			&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &user.Settings, &user.CreatedAt, &user.UpdatedAt,
		)
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		// Mock auth - query by internal UUID
		err = h.db.QueryRow(r.Context(), `
			SELECT id, firebase_uid, email, username, settings, created_at, updated_at
			FROM users
			WHERE id=$1`, testUserID).Scan(
			&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &user.Settings, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err != nil {
		// If user not found, maybe return 404 or just create one (auto-provisioning)?
		// Typically GetMe returns 404 if user not registered.
		slog.Error("User not found", "error", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}
