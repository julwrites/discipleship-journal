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

	"database/sql"
	"discipleship_journal_api/middleware"

	"firebase.google.com/go/v4/auth"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_CreateOrUpdateUser(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := NewUserHandler(db)
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
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=?")).
			WithArgs(firebaseUID).
			WillReturnError(sql.ErrNoRows)

		// 2. INSERT
		settingsJSON := []byte(`{"theme":"dark"}`)
		testUser := "TestUser"
		rows := sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
			AddRow(userUUID.String(), firebaseUID, email, &testUser, settingsJSON, time.Now(), time.Now())

		mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO users (firebase_uid, email, username, settings)")).
			WithArgs(firebaseUID, email, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockDB.ExpectQuery("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=\\?").
			WithArgs(firebaseUID).
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
		rows := sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
			AddRow(userUUID.String(), firebaseUID, email, &oldUser, settingsJSON, time.Now(), time.Now())

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=?")).
			WithArgs(firebaseUID).
			WillReturnRows(rows)

		// 2. UPDATE
		updatedUser := "TestUser"
		updatedSettings := []byte(`{"theme":"dark"}`)
		updatedRows := sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
			AddRow(userUUID.String(), firebaseUID, email, &updatedUser, updatedSettings, time.Now(), time.Now())

		mockDB.ExpectExec(regexp.QuoteMeta("UPDATE users SET updated_at = NOW(), username = ?, settings = ? WHERE firebase_uid = ?")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), firebaseUID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockDB.ExpectQuery("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=\\?").
			WithArgs(firebaseUID).
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

	t.Run("Create User (Upsert - Test Mode)", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(reqBody))

		// Inject TestUserKey
		ctx := context.WithValue(req.Context(), TestUserKey, userUUID.String())
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		expectedFirebaseUID := "test-firebase-uid-" + userUUID.String()

		// 1. SELECT returns No Rows (Search by ID)
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE id=?")).
			WithArgs(userUUID.String()).
			WillReturnError(sql.ErrNoRows)

		// 2. INSERT with Explicit ID
		settingsJSON := []byte(`{"theme":"dark"}`)
		testUser := "TestUser"
		rows := sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
			AddRow(userUUID.String(), expectedFirebaseUID, email, &testUser, settingsJSON, time.Now(), time.Now())

		mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO users (id, firebase_uid, email, username, settings) VALUES (?, ?, ?, ?, ?)")).
			WithArgs(userUUID.String(), expectedFirebaseUID, email, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockDB.ExpectQuery("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE id=\\?").
			WithArgs(userUUID.String()).
			WillReturnRows(rows)

		handler.CreateOrUpdateUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var user User
		err := json.NewDecoder(w.Body).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, userUUID.String(), user.ID)
		assert.Equal(t, expectedFirebaseUID, user.FirebaseUID)
	})

	t.Run("Update User (Existing - Test Mode)", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/users", strings.NewReader(reqBody))

		// Inject TestUserKey
		ctx := context.WithValue(req.Context(), TestUserKey, userUUID.String())
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		expectedFirebaseUID := "test-firebase-uid-" + userUUID.String()

		// 1. SELECT returns User (Search by ID)
		oldUser := "OldName"
		settingsJSON := []byte("{}")
		rows := sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
			AddRow(userUUID.String(), expectedFirebaseUID, email, &oldUser, settingsJSON, time.Now(), time.Now())

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE id=?")).
			WithArgs(userUUID.String()).
			WillReturnRows(rows)

		// 2. UPDATE with WHERE id
		updatedUser := "TestUser"
		updatedSettings := []byte(`{"theme":"dark"}`)
		updatedRows := sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
			AddRow(userUUID.String(), expectedFirebaseUID, email, &updatedUser, updatedSettings, time.Now(), time.Now())

		mockDB.ExpectExec(regexp.QuoteMeta("UPDATE users SET updated_at = NOW(), username = ?, settings = ? WHERE id = ?")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userUUID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockDB.ExpectQuery("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE id=\\?").
			WithArgs(userUUID.String()).
			WillReturnRows(updatedRows)

		handler.UpdateUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var user User
		err := json.NewDecoder(w.Body).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, "TestUser", *user.Username)
	})
}

func TestUserHandler_GetMe(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := NewUserHandler(db)
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
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=?")).
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
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

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=?")).
			WithArgs(firebaseUID).
			WillReturnError(sql.ErrNoRows)

		handler.GetMe(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUserHandler_CreateOrUpdateUser_MergesSettings(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := NewUserHandler(db)
	userUUID := uuid.New()
	firebaseUID := "test-firebase-uid"
	email := "test@example.com"

	// Original settings in DB
	originalSettings := []byte(`{"theme":"dark", "notifications": true}`)

	// Request updates only bible_version
	reqBody := `{"settings": {"bible_version": "ESV"}}`

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

	// 1. SELECT returns User with original settings
	oldUser := "OldName"
	rows := sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
		AddRow(userUUID.String(), firebaseUID, email, &oldUser, originalSettings, time.Now(), time.Now())

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=?")).
		WithArgs(firebaseUID).
		WillReturnRows(rows)

	// 2. UPDATE should use merged settings
	// Expected merged settings: theme: dark, notifications: true, bible_version: ESV
	// The order of keys in JSON map is not guaranteed, so we might need a more robust check or rely on stable marshalling if keys are sorted.
	// Go's json.Marshal sorts map keys.
	// {"bible_version":"ESV","notifications":true,"theme":"dark"}
	expectedSettings := []byte(`{"bible_version":"ESV","notifications":true,"theme":"dark"}`)

	updatedUser := "OldName"
	updatedRows := sqlmock.NewRows([]string{"id", "firebase_uid", "email", "username", "settings", "created_at", "updated_at"}).
		AddRow(userUUID.String(), firebaseUID, email, &updatedUser, expectedSettings, time.Now(), time.Now())

	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE users SET updated_at = NOW(), settings = ? WHERE firebase_uid = ?")).
		WithArgs(expectedSettings, firebaseUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockDB.ExpectQuery("SELECT id, firebase_uid, email, username, settings, created_at, updated_at FROM users WHERE firebase_uid=\\?").
		WithArgs(firebaseUID).
		WillReturnRows(updatedRows)

	handler.UpdateUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var user User
	err = json.NewDecoder(w.Body).Decode(&user)
	assert.NoError(t, err)

	// Verify settings in response
	assert.Equal(t, "ESV", user.Settings["bible_version"])
	assert.Equal(t, "dark", user.Settings["theme"])
	assert.Equal(t, true, user.Settings["notifications"])

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
