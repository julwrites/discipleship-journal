package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/validation"
	"firebase.google.com/go/v4/auth"
)

type User struct {
	ID          string                 `json:"id"`
	FirebaseUID string                 `json:"firebase_uid"`
	Email       string                 `json:"email"`
	Username    *string                `json:"username"`
	Settings    map[string]interface{} `json:"settings"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type UpdateUserRequest struct {
	Username *string                `json:"username" validate:"omitempty,min=3,max=30,alphanum"`
	Settings map[string]interface{} `json:"settings" validate:"omitempty"`
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
// @Param user body UpdateUserRequest false "User details (optional)"
// @Success 200 {object} User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateOrUpdateUser(w http.ResponseWriter, r *http.Request) {
	var uid string
	var email string

	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		uid = token.UID
		// Extract email from claims
		if emailClaim, ok := token.Claims["email"].(string); ok {
			email = emailClaim
		}
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		// Mock ID, but we need a uid for the DB query logic below if we were to support it
		// For now just error if not authenticated properly
		if testUserID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// In test mode without firebase token, we might not have email easily unless injected
		// Assume "test@example.com" if missing for test
		email = "test@example.com"
		uid = "test-firebase-uid"
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateUserRequest
	// Manually decode to allow empty body (EOF)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if !errors.Is(err, io.EOF) {
			slog.Error("Failed to decode request body", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// EOF is fine, req will be zero-valued
	} else {
		// Validate if body was present
		if err := validation.ValidateStruct(&req); err != nil {
			slog.Error("Validation failed", "error", err)
			http.Error(w, "Validation error", http.StatusBadRequest)
			return
		}
	}

	// Prepare for Upsert
	var user User
	var settingsBytes []byte

	// 1. Try to find user
	err := h.db.QueryRow(r.Context(), `
		SELECT id, firebase_uid, email, username, settings, created_at, updated_at
		FROM users
		WHERE firebase_uid=$1`, uid).Scan(
		&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &settingsBytes, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			// CREATE (User Not Found)
			var settingsJSON []byte
			if req.Settings != nil {
				settingsJSON, _ = json.Marshal(req.Settings)
			} else {
				// Initialize with empty object
				settingsJSON = []byte("{}")
			}

			err = h.db.QueryRow(r.Context(), `
				INSERT INTO users (firebase_uid, email, username, settings)
				VALUES ($1, $2, $3, $4)
				RETURNING id, firebase_uid, email, username, settings, created_at, updated_at`,
				uid, email, req.Username, settingsJSON).Scan(
				&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &settingsBytes, &user.CreatedAt, &user.UpdatedAt,
			)
			if err != nil {
				slog.Error("Failed to create user", "error", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			// Parse settingsBytes if needed
			if len(settingsBytes) > 0 {
				if err := json.Unmarshal(settingsBytes, &user.Settings); err != nil {
					slog.Error("Failed to unmarshal settings", "error", err)
				}
			}

			// Respond and return immediately - no need to update
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(user); err != nil {
				slog.Error("Failed to encode response", "error", err)
			}
			return
		}

		// DB Error
		slog.Error("Failed to query user", "error", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 2. UPDATE (User Found)
	// Check if update is needed
	shouldUpdate := false
	updateQuery := "UPDATE users SET updated_at = NOW()"
	var args []interface{}
	argIdx := 1

	if req.Username != nil {
		updateQuery += fmt.Sprintf(", username = $%d", argIdx)
		args = append(args, req.Username)
		argIdx++
		shouldUpdate = true
	}
	if req.Settings != nil {
		settingsJSON, _ := json.Marshal(req.Settings)
		updateQuery += fmt.Sprintf(", settings = $%d", argIdx)
		args = append(args, settingsJSON)
		argIdx++
		shouldUpdate = true
	}

	if shouldUpdate {
		updateQuery += fmt.Sprintf(" WHERE firebase_uid = $%d", argIdx)
		args = append(args, uid)
		updateQuery += " RETURNING id, firebase_uid, email, username, settings, created_at, updated_at"

		err = h.db.QueryRow(r.Context(), updateQuery, args...).Scan(
			&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &settingsBytes, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to update user", "error", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	if len(settingsBytes) > 0 {
		if err := json.Unmarshal(settingsBytes, &user.Settings); err != nil {
			slog.Error("Failed to unmarshal settings", "error", err)
		}
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
	var settingsBytes []byte

	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok {
		// Real Firebase auth - query by firebase_uid
		uid := token.UID
		err = h.db.QueryRow(r.Context(), `
			SELECT id, firebase_uid, email, username, settings, created_at, updated_at
			FROM users
			WHERE firebase_uid=$1`, uid).Scan(
			&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &settingsBytes, &user.CreatedAt, &user.UpdatedAt,
		)
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		// Mock auth - query by internal UUID
		err = h.db.QueryRow(r.Context(), `
			SELECT id, firebase_uid, email, username, settings, created_at, updated_at
			FROM users
			WHERE id=$1`, testUserID).Scan(
			&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &settingsBytes, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err != nil {
		slog.Error("User not found", "error", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if len(settingsBytes) > 0 {
		if err := json.Unmarshal(settingsBytes, &user.Settings); err != nil {
			slog.Error("Failed to unmarshal settings", "error", err)
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}
