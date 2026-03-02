package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"discipleship_journal_api/middleware"

	"firebase.google.com/go/v4/auth"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGroupHandler_CreateGroup_Auth_UserNotFound(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Override dbProvider
	oldProvider := dbProvider
	dbProvider = func() DBInterface { return db }
	defer func() { dbProvider = oldProvider }()

	req := httptest.NewRequest("POST", "/api/groups", strings.NewReader(`{"name":"Group"}`))
	// Inject Token but NO TestUserKey
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Expect GetUserUUID to fail
	mockDB.ExpectQuery("SELECT id FROM users").
		WithArgs("firebase-uid").
		WillReturnError(errors.New("db error"))

	handler.CreateGroup(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

func TestGroupShareHandler_ShareItem_InvalidBody(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	_ = mockDB

	_ = db

	handler := NewGroupShareHandler(db, nil)

	req := httptest.NewRequest("POST", "/api/groups/1/shares", strings.NewReader(`{invalid}`))
	w := httptest.NewRecorder()

	handler.ShareItemToGroup(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
