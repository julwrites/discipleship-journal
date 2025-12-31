package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_CreateOrUpdateUser(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	handler := NewUserHandler(mockDB)
	userUUID := uuid.New()
	firebaseUID := "test-firebase-uid"
	email := "test@example.com"
	reqBody := `{"full_name": "Test User", "avatar_url": "http://example.com/avatar.jpg"}`

	t.Run("Create User (Upsert)", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(reqBody))

		// Mock context with user token
		token := &auth.Token{
			UID: firebaseUID,
			Claims: map[string]interface{}{
				"email": email,
			},
		}
		// Manually inject context as middleware.WithUser is not exported or available
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		// Expect SELECT (handler doesn't actually insert/update, just reads)
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=$1")).
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
				AddRow(userUUID.String(), firebaseUID, email, "Test User", "{}", time.Now(), time.Now()))

		handler.CreateOrUpdateUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var user User
		err := json.NewDecoder(w.Body).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, userUUID.String(), user.ID)
		assert.Equal(t, firebaseUID, user.FirebaseUID)
		assert.Equal(t, "Test User", user.Username)
	})

	t.Run("Update User Handler", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/users", strings.NewReader(reqBody))

		token := &auth.Token{
			UID: firebaseUID,
			Claims: map[string]interface{}{
				"email": email,
			},
		}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		// Expect SELECT (handler doesn't actually insert/update, just reads)
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=$1")).
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
				AddRow(userUUID.String(), firebaseUID, email, "Test User", "{}", time.Now(), time.Now()))

		handler.UpdateUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		handler.CreateOrUpdateUser(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUserHandler_GetMe(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	handler := NewUserHandler(mockDB)
	userUUID := uuid.New()
	firebaseUID := "test-firebase-uid"

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/me", nil)

		token := &auth.Token{
			UID: firebaseUID,
		}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=$1")).
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
				AddRow(userUUID.String(), firebaseUID, "test@example.com", "Test User", "{}", time.Now(), time.Now()))

		handler.GetMe(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var user User
		err := json.NewDecoder(w.Body).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, userUUID.String(), user.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/me", nil)

		token := &auth.Token{
			UID: firebaseUID,
		}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=$1")).
			WithArgs(firebaseUID).
			WillReturnError(pgx.ErrNoRows)

		handler.GetMe(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
