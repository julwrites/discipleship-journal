package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"discipleship_journal_api/database"
	"firebase.google.com/go/v4/auth"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  *string   `json:"username"`
	Settings  UserSettings `json:"settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserSettings struct {
	BibleVersion string `json:"bible_version"`
}

// CreateOrUpdateUser is called after login to sync Firebase user to our DB
func CreateOrUpdateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("user").(*auth.Token)
	uid := token.UID
	email := token.Claims["email"].(string)

	var user User
	// Check if user exists
	err := database.DB.QueryRow(r.Context(), "SELECT id, email, username, settings FROM users WHERE id=$1", uid).Scan(&user.ID, &user.Email, &user.Username, &user.Settings)

	if err != nil {
		// Create new user
		_, err := database.DB.Exec(r.Context(), "INSERT INTO users (id, email) VALUES ($1, $2)", uid, email)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		user.ID = uid
		user.Email = email
	}

	json.NewEncoder(w).Encode(user)
}

type UpdateUserRequest struct {
	Username     *string `json:"username"`
	BibleVersion *string `json:"bible_version"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("user").(*auth.Token)
	uid := token.UID

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username != nil {
		_, err := database.DB.Exec(r.Context(), "UPDATE users SET username=$1, updated_at=NOW() WHERE id=$2", req.Username, uid)
		if err != nil {
			// Handle unique constraint violation
			http.Error(w, "Failed to update username (might be taken)", http.StatusConflict)
			return
		}
	}

    // Simple settings update logic (could be more robust merging JSON)
	if req.BibleVersion != nil {
        settingsJSON := map[string]string{"bible_version": *req.BibleVersion}
		_, err := database.DB.Exec(r.Context(), "UPDATE users SET settings=$1, updated_at=NOW() WHERE id=$2", settingsJSON, uid)
		if err != nil {
			http.Error(w, "Failed to update settings", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
