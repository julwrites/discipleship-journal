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
	// Updated request body to match new struct and validation (alphanum)
	reqBody := `{"username": "TestUser", "settings": {"theme": "dark"}}`

	t.Run("Create User (Upsert - New User)", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(reqBody))

		token := &auth.Token{
			UID: firebaseUID,
			Claims: map[string]interface{}{
				"email": email,
			},
		}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		// 1. SELECT returns No Rows
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=$1")).
			WithArgs(firebaseUID).
			WillReturnError(pgx.ErrNoRows)

		// 2. INSERT
		settingsJSON := []byte(`{"theme":"dark"}`)
		testUser := "TestUser"
		rows := pgxmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
			AddRow(userUUID.String(), firebaseUID, email, &testUser, settingsJSON, time.Now(), time.Now())

		mockDB.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (firebase_uid, email, username, settings)")).
			WithArgs(firebaseUID, email, pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(rows)

		handler.CreateOrUpdateUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var user User
		err := json.NewDecoder(w.Body).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, userUUID.String(), user.ID)
		assert.Equal(t, firebaseUID, user.FirebaseUID)
		assert.NotNil(t, user.Username)
		assert.Equal(t, "TestUser", *user.Username)
	})

	t.Run("Update User Handler (Existing User)", func(t *testing.T) {
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

		// 1. SELECT returns User
		oldUser := "OldName"
		settingsJSON := []byte("{}")
		rows := pgxmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
			AddRow(userUUID.String(), firebaseUID, email, &oldUser, settingsJSON, time.Now(), time.Now())

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=$1")).
			WithArgs(firebaseUID).
			WillReturnRows(rows)

		// 2. UPDATE
		updatedUser := "TestUser"
		updatedSettings := []byte(`{"theme":"dark"}`)
		updatedRows := pgxmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
			AddRow(userUUID.String(), firebaseUID, email, &updatedUser, updatedSettings, time.Now(), time.Now())

		mockDB.ExpectQuery(regexp.QuoteMeta("UPDATE users SET updated_at = NOW(), username = $1, settings = $2 WHERE firebase_uid = $3 RETURNING")).
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), firebaseUID).
			WillReturnRows(updatedRows)

		handler.UpdateUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var user User
		err := json.NewDecoder(w.Body).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, "TestUser", *user.Username)
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

		testUser := "TestUser"
		settingsJSON := []byte("{}")
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=$1")).
			WithArgs(firebaseUID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
				AddRow(userUUID.String(), firebaseUID, "test@example.com", &testUser, settingsJSON, time.Now(), time.Now()))

		handler.GetMe(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var user User
		err := json.NewDecoder(w.Body).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, userUUID.String(), user.ID)
		assert.Equal(t, "TestUser", *user.Username)
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
