package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestSendConnectionRequest(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mockNotification := services.NewMockNotificationService()
	handler := NewConnectionHandler(mock, mockNotification)

	requesterUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	receiverEmail := "user2@example.com"
	receiverUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	connID := "00000000-0000-0000-0000-000000000003"

	// Mock receiver UUID lookup by email
	mock.ExpectQuery("SELECT id FROM users WHERE email =").
		WithArgs(receiverEmail).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(receiverUUID.String()))

	// Mock insert connection
	mock.ExpectQuery("INSERT INTO connections").
		WithArgs(requesterUUID, receiverUUID.String()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(connID))

	// Expect synchronous requester name lookup
	mock.ExpectQuery("SELECT username FROM users WHERE id").
		WithArgs(requesterUUID).
		WillReturnRows(pgxmock.NewRows([]string{"username"}).AddRow("RequesterUser"))

	reqBody := ConnectionRequest{ReceiverEmail: receiverEmail}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Use TestUserKey to bypass global DB lookup for user ID
	ctx := context.WithValue(req.Context(), TestUserKey, requesterUUID) // Context value is UUID
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.SendConnectionRequest(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSearchUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mockNotification := services.NewMockNotificationService()
	handler := NewConnectionHandler(mock, mockNotification)

	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Expect search query
	username := "johndoe"
	mock.ExpectQuery(`SELECT u.id, u.email, u.username.*`).
		WithArgs("johndoe", userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "email", "username", "is_connected"}).
			AddRow("00000000-0000-0000-0000-000000000002", "john@example.com", &username, false))

	req := httptest.NewRequest("GET", "/api/users/search?q=johndoe", nil)
	// Use TestUserKey to bypass global DB lookup for user ID
	ctx := context.WithValue(req.Context(), TestUserKey, userUUID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.SearchUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var resp []map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	// Email should be masked (nil/missing) because q != email and not connected
	assert.Nil(t, resp[0]["email"])
	assert.Equal(t, "johndoe", resp[0]["username"])

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListConnections(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mockNotification := services.NewMockNotificationService()
	handler := NewConnectionHandler(mock, mockNotification)

	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	connID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	requesterUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	// Mock list connections
	requesterUsername := "user2"
	receiverUsername := "user1"
	mock.ExpectQuery("SELECT c.id, c.requester_id, c.receiver_id, c.status").
		WithArgs(userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "requester_id", "receiver_id", "status", "requester_email", "receiver_email", "requester_username", "receiver_username"}).
			AddRow(connID.String(), requesterUUID.String(), userUUID.String(), "pending", "user2@example.com", "user1@example.com", &requesterUsername, &receiverUsername))

	req := httptest.NewRequest("GET", "/api/connections", nil)
	// Use TestUserKey to bypass global DB lookup for user ID
	ctx := context.WithValue(req.Context(), TestUserKey, userUUID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ListConnections(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var conns []ConnectionResponse
	err = json.Unmarshal(w.Body.Bytes(), &conns)
	assert.NoError(t, err)
	assert.Len(t, conns, 1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAcceptConnectionRequest(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mockNotification := services.NewMockNotificationService()
	handler := NewConnectionHandler(mock, mockNotification)

	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	connID := "00000000-0000-0000-0000-000000000003"

	// Mock update connection
	mock.ExpectExec("UPDATE connections SET status = 'accepted'").
		WithArgs(connID, userUUID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	req := httptest.NewRequest("PUT", "/api/connections/"+connID, nil)

	// Setup Chi context for URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", connID)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)

	// Use TestUserKey to bypass global DB lookup for user ID
	ctx = context.WithValue(ctx, TestUserKey, userUUID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.AcceptConnectionRequest(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]bool
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp["success"])

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteConnectionRequest(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mockNotification := services.NewMockNotificationService()
	handler := NewConnectionHandler(mock, mockNotification)

	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	connID := "00000000-0000-0000-0000-000000000003"

	// Mock delete connection
	mock.ExpectExec("DELETE FROM connections").
		WithArgs(connID, userUUID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	req := httptest.NewRequest("DELETE", "/api/connections/"+connID, nil)

	// Setup Chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", connID)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)

	// Use TestUserKey to bypass global DB lookup for user ID
	ctx = context.WithValue(ctx, TestUserKey, userUUID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.DeleteConnectionRequest(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]bool
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp["success"])

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
