package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"discipleship_journal_api/database"
	"discipleship_journal_api/middleware"
	"discipleship_journal_api/validation"
	"firebase.google.com/go/v4/auth"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID          string       `json:"id"`
	FirebaseUID string       `json:"firebase_uid"`
	Email       string       `json:"email"`
	Username    *string      `json:"username"`
	Settings    UserSettings `json:"settings"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type UserSettings struct {
	BibleVersion string `json:"bible_version"`
}

// Helper to get Internal User UUID from Firebase UID
func GetUserUUID(ctx context.Context, firebaseUID string) (string, error) {
	var id string
	err := database.DB.QueryRow(ctx, "SELECT id FROM users WHERE firebase_uid=$1", firebaseUID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// CreateOrUpdateUser godoc
// @Summary Create or update user
// @Description Sync Firebase user to database
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/users/me [post]
func CreateOrUpdateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID
	email := token.Claims["email"].(string)

	var user User
	// Check if user exists using firebase_uid
	err := database.DB.QueryRow(r.Context(),
		"SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=$1",
		uid).Scan(&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &user.Settings, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Create new user
			// We return the new user details including the generated UUID
			err = database.DB.QueryRow(r.Context(),
				"INSERT INTO users (firebase_uid, email) VALUES ($1, $2) RETURNING id, firebase_uid, email, username, settings, created_at, updated_at",
				uid, email).Scan(&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &user.Settings, &user.CreatedAt, &user.UpdatedAt)

			if err != nil {
				http.Error(w, "Failed to create user: " + err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else {
		// Update email if changed (optional, but good practice)
		if user.Email != email {
			_, err = database.DB.Exec(r.Context(), "UPDATE users SET email=$1, updated_at=NOW() WHERE id=$2", email, user.ID)
			if err != nil {
				// Log error
				fmt.Println("Failed to update email:", err)
			}
			user.Email = email
		}
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

type UpdateUserRequest struct {
	Username     *string `json:"username" validate:"omitempty,min=3,max=30"`
	BibleVersion *string `json:"bible_version" validate:"omitempty,oneof=ESV NIV KJV"`
}

// UpdateUser godoc
// @Summary Update user profile
// @Description Update user settings like username and bible version
// @Tags users
// @Accept json
// @Produce json
// @Param request body UpdateUserRequest true "Update User Request"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/users/me [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.UserContextKey).(*auth.Token)
	uid := token.UID

	userUUID, err := GetUserUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var req UpdateUserRequest
	if !validation.DecodeAndValidate(w, r, &req) {
		return
	}

	if req.Username != nil {
		_, err := database.DB.Exec(r.Context(), "UPDATE users SET username=$1, updated_at=NOW() WHERE id=$2", req.Username, userUUID)
		if err != nil {
			// Handle unique constraint violation
			http.Error(w, "Failed to update username (might be taken)", http.StatusConflict)
			return
		}
	}

	if req.BibleVersion != nil {
		// Need to merge settings carefully. For now, we fetch, merge, update.
		// Or using jsonb_set provided by Postgres.
		// Simple approach: jsonb_set(settings, '{bible_version}', '"ESV"')

		// Note: jsonb_set value must be jsonb. to_jsonb($1::text) might work.
		_, err := database.DB.Exec(r.Context(),
			"UPDATE users SET settings = jsonb_set(COALESCE(settings, '{}'::jsonb), '{bible_version}', to_jsonb($1::text)), updated_at=NOW() WHERE id=$2",
			*req.BibleVersion, userUUID)

		if err != nil {
			http.Error(w, "Failed to update settings: " + err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
