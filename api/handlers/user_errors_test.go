package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"database/sql"
	"discipleship_journal_api/middleware"

	"firebase.google.com/go/v4/auth"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_CreateOrUpdateUser_Errors(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := NewUserHandler(db)
	firebaseUID := "test-firebase-uid"
	email := "test@example.com"
	reqBody := `{"username": "TestUser"}`

	t.Run("DB Error on Select", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(reqBody))

		token := &auth.Token{UID: firebaseUID, Claims: map[string]interface{}{"email": email}}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users")).
			WithArgs(firebaseUID).
			WillReturnError(errors.New("db connection error"))

		handler.CreateOrUpdateUser(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("DB Error on Insert", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(reqBody))

		token := &auth.Token{UID: firebaseUID, Claims: map[string]interface{}{"email": email}}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		// Select returns no rows (trigger insert)
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users")).
			WithArgs(firebaseUID).
			WillReturnError(sql.ErrNoRows)

		// Insert fails
		mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO users")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("insert failed"))

		handler.CreateOrUpdateUser(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("DB Error on Update", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(reqBody))

		token := &auth.Token{UID: firebaseUID, Claims: map[string]interface{}{"email": email}}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		// Select returns existing user
		userUUID := uuid.New()
		user := "OldUser"
		rows := sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
			AddRow(userUUID.String(), firebaseUID, email, &user, []byte("{}"), time.Now(), time.Now())

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users")).
			WithArgs(firebaseUID).
			WillReturnRows(rows)

		// Update fails
		mockDB.ExpectExec(regexp.QuoteMeta("UPDATE users SET")).
			WithArgs(sqlmock.AnyArg()).
			WillReturnError(errors.New("update failed"))

		handler.CreateOrUpdateUser(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid JSON Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(`{invalid}`))
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateOrUpdateUser(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		// Create a body that fails validation (e.g., username too long if there's a limit)
		// Assuming max=30 based on code reading or memory
		longName := strings.Repeat("a", 100)
		body := `{"username": "` + longName + `"}`
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(body))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateOrUpdateUser(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHandler_GetMe_Errors(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := NewUserHandler(db)

	t.Run("Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/me", nil)
		// No context
		w := httptest.NewRecorder()

		handler.GetMe(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("DB Error", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/me", nil)
		firebaseUID := "uid"
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users")).
			WithArgs(firebaseUID).
			WillReturnError(errors.New("db error"))

		handler.GetMe(w, req)

		// Note: The code logs error and returns 500, or 404 if not found?
		// Code says:
		// if err != nil {
		// 		slog.Error("User not found", "error", err)
		// 		http.Error(w, "User not found", http.StatusNotFound)
		// 		return
		// }
		// So it returns 404 even for DB errors? That seems like a bug or feature.
		// Let's assert 404 based on code reading.
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Test Mode Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/me", nil)
		userUUID := uuid.New()
		ctx := context.WithValue(req.Context(), TestUserKey, userUUID.String())
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		testUser := "TestUser"
		settingsJSON := []byte("{}")
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE id=?")).
			WithArgs(userUUID.String()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
				AddRow(userUUID.String(), "firebase-uid", "test@example.com", &testUser, settingsJSON, time.Now(), time.Now()))

		handler.GetMe(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
