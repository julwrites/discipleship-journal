package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/database"
	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateGroup(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	database.DB = mock

	// Mock User
	userID := "test-firebase-uid"
	userUUID := "00000000-0000-0000-0000-000000000001"

	// Mock Request
	reqBody := CreateGroupRequest{
		Name:        "Test Group",
		Description: "A test group",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/groups", bytes.NewBuffer(body))

	// Inject Auth Context
	token := &auth.Token{UID: userID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	// Mock DB Expectations
	// 1. GetUserUUID
	mock.ExpectQuery("SELECT id FROM users WHERE firebase_uid").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// 2. Transaction Begin
	mock.ExpectBegin()

	// 3. Insert Group
	mock.ExpectQuery("INSERT INTO groups").
		WithArgs(reqBody.Name, reqBody.Description, userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("group-123"))

	// 4. Insert Member (Admin)
	mock.ExpectExec("INSERT INTO group_members").
		WithArgs("group-123", userUUID).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// 5. Commit
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	CreateGroup(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify JSON response
	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "group-123", resp["id"])

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListMyGroups(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	database.DB = mock

	userID := "test-firebase-uid"
	userUUID := "00000000-0000-0000-0000-000000000001"

	req := httptest.NewRequest("GET", "/groups", nil)
	token := &auth.Token{UID: userID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	mock.ExpectQuery("SELECT id FROM users WHERE firebase_uid").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	mock.ExpectQuery("SELECT g.id, g.name").
		WithArgs(userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_by", "role"}).
			AddRow("g1", "Group 1", "Desc 1", "creator1", "admin").
			AddRow("g2", "Group 2", "Desc 2", "creator2", "member"))

	w := httptest.NewRecorder()
	ListMyGroups(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var groups []GroupResponse
	err = json.Unmarshal(w.Body.Bytes(), &groups)
	assert.NoError(t, err)
	assert.Len(t, groups, 2)
	assert.Equal(t, "Group 1", groups[0].Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestJoinGroup(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	database.DB = mock

	userID := "test-uid"
	userUUID := "uuid-1"
	groupID := "group-1"

	req := httptest.NewRequest("POST", "/groups/"+groupID+"/join", nil)
	// Add chi URL param mocking if needed, or rely on chi.URLParam reading from context if we use chi in test.
	// Since the handler uses chi.URLParam, we need to mock it.
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", groupID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	token := &auth.Token{UID: userID}
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, token))

	// GetUserUUID
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// Check exists
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(groupID, userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	// Insert
	mock.ExpectExec("INSERT INTO group_members").
		WithArgs(groupID, userUUID).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	w := httptest.NewRecorder()
	JoinGroup(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
