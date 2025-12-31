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
	var uid, email string
	token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	uid = token.UID
	if e, found := token.Claims["email"]; found {
		email = e.(string)
	}

	var req UpdateUserRequest
	if !DecodeAndValidate(w, r, &req) {
		return
	}

	var user User
	// Upsert logic using ON CONFLICT
	// Note: firebase_uid should have a unique constraint
	err := h.db.QueryRow(r.Context(), `
		INSERT INTO users (firebase_uid, email, full_name, avatar_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (firebase_uid) DO UPDATE
		SET full_name = EXCLUDED.full_name,
			avatar_url = EXCLUDED.avatar_url,
			email = EXCLUDED.email,
			updated_at = NOW()
		RETURNING id, firebase_uid, email, full_name, avatar_url, created_at, updated_at`,
		uid, email, req.FullName, req.AvatarURL).Scan(
		&user.ID, &user.FirebaseUID, &user.Email, &user.FullName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
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
	token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	uid := token.UID

	var user User
	err := h.db.QueryRow(r.Context(), `
		SELECT id, firebase_uid, email, full_name, avatar_url, created_at, updated_at
		FROM users
		WHERE firebase_uid=$1`, uid).Scan(
		&user.ID, &user.FirebaseUID, &user.Email, &user.FullName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		// If user not found, maybe return 404 or just create one (auto-provisioning)?
		// Typically GetMe returns 404 if user not registered.
		slog.Error("User not found", "error", err, "uid", uid)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}
