package handlers

import (
	"encoding/json"
	"errors"
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
	Username *string                `json:"username" validate:"omitempty,max=30"`
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
	var testUserID string
	isTestMode := false

	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		uid = token.UID
		// Extract email from claims
		if emailClaim, ok := token.Claims["email"].(string); ok {
			email = emailClaim
		}
	} else if val, ok := r.Context().Value(TestUserKey).(string); ok {
		testUserID = val
		isTestMode = true
		if testUserID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// In test mode without firebase token, we might not have email easily unless injected
		// Assume "test@example.com" if missing for test
		email = "test@example.com"
		// Generate unique firebase UID for test user
		uid = "test-firebase-uid-" + testUserID
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
	var selectQuery string
	var selectArgs []interface{}

	if isTestMode {
		selectQuery = `
			SELECT id, firebase_uid, email, username, settings, created_at, updated_at
			FROM users
			WHERE id=?`
		selectArgs = []interface{}{testUserID}
	} else {
		selectQuery = `
			SELECT id, firebase_uid, email, username, settings, created_at, updated_at
			FROM users
			WHERE firebase_uid=?`
		selectArgs = []interface{}{uid}
	}

	err := h.db.QueryRowContext(r.Context(), selectQuery, selectArgs...).Scan(
		&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &settingsBytes, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			// CREATE (User Not Found)
			var settingsJSON []byte
			if req.Settings != nil {
				settingsJSON, _ = json.Marshal(req.Settings)
			} else {
				settingsJSON = []byte("{}")
			}

			var usernameVal interface{} = req.Username
			if req.Username != nil && *req.Username == "" {
				usernameVal = nil
			}

			var insertQuery string
			var insertArgs []interface{}

			if isTestMode {
				insertQuery = `
					INSERT INTO users (id, firebase_uid, email, username, settings)
					VALUES (?, ?, ?, ?, ?)`
				insertArgs = []interface{}{testUserID, uid, email, usernameVal, settingsJSON}
			} else {
				insertQuery = `
					INSERT INTO users (firebase_uid, email, username, settings)
					VALUES (?, ?, ?, ?)`
				insertArgs = []interface{}{uid, email, usernameVal, settingsJSON}
			}

			_, err = h.db.ExecContext(r.Context(), insertQuery, insertArgs...)
			if err == nil {
				err = h.db.QueryRowContext(r.Context(), selectQuery, selectArgs...).Scan(
					&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &settingsBytes, &user.CreatedAt, &user.UpdatedAt,
				)
			}
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

	if req.Username != nil {
		updateQuery += ", username = ?"
		var val interface{} = req.Username
		if *req.Username == "" {
			val = nil
		}
		args = append(args, val)
		shouldUpdate = true
	}
	if req.Settings != nil {
		// Merge with existing settings
		var currentSettings map[string]interface{}
		if len(settingsBytes) > 0 {
			if err := json.Unmarshal(settingsBytes, &currentSettings); err != nil {
				slog.Error("Failed to unmarshal existing settings during update", "error", err)
				// Fallback to empty map
				currentSettings = make(map[string]interface{})
			}
		} else {
			currentSettings = make(map[string]interface{})
		}

		for k, v := range req.Settings {
			currentSettings[k] = v
		}

		settingsJSON, _ := json.Marshal(currentSettings)
		updateQuery += ", settings = ?"
		args = append(args, settingsJSON)
		shouldUpdate = true
	}

	if shouldUpdate {
		// Use ID for update if test mode, or firebase_uid if production
		if isTestMode {
			updateQuery += " WHERE id = ?"
			args = append(args, testUserID)
		} else {
			updateQuery += " WHERE firebase_uid = ?"
			args = append(args, uid)
		}
		_, err := h.db.ExecContext(r.Context(), updateQuery, args...)
		if err == nil {
			err = h.db.QueryRowContext(r.Context(), selectQuery, selectArgs...).Scan(
				&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &settingsBytes, &user.CreatedAt, &user.UpdatedAt,
			)
		}
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

	if token, ok := r.Context().Value(middleware.UserContextKey).(*auth.Token); ok && token != nil {
		// Real Firebase auth - query by firebase_uid
		uid := token.UID
		err = h.db.QueryRowContext(r.Context(), `
			SELECT id, firebase_uid, email, username, settings, created_at, updated_at
			FROM users
			WHERE firebase_uid=?`, uid).Scan(
			&user.ID, &user.FirebaseUID, &user.Email, &user.Username, &settingsBytes, &user.CreatedAt, &user.UpdatedAt,
		)
	} else if testUserID, ok := r.Context().Value(TestUserKey).(string); ok {
		// Mock auth - query by internal UUID
		err = h.db.QueryRowContext(r.Context(), `
			SELECT id, firebase_uid, email, username, settings, created_at, updated_at
			FROM users
			WHERE id=?`, testUserID).Scan(
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
