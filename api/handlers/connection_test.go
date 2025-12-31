package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/services"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
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

	uid := "firebase-uid-1"
	requesterUUID := "user-uuid-1"
	receiverEmail := "user2@example.com"
	receiverUUID := "user-uuid-2"

	// Mock requester UUID lookup
	mock.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
		WithArgs(uid).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(requesterUUID))

	// Mock receiver UUID lookup by email
	mock.ExpectQuery("SELECT id FROM users WHERE email =").
		WithArgs(receiverEmail).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(receiverUUID))

	// Mock insert connection
	mock.ExpectQuery("INSERT INTO connections").
		WithArgs(requesterUUID, receiverUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("conn-1"))

	// Expect synchronous requester name lookup
	mock.ExpectQuery("SELECT display_name FROM users WHERE id").
		WithArgs(requesterUUID).
		WillReturnRows(pgxmock.NewRows([]string{"display_name"}).AddRow("Requester Name"))

	reqBody := ConnectionRequest{ReceiverEmail: receiverEmail}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	token := &auth.Token{UID: uid}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
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

	uid := "firebase-uid-1"
	userUUID := "user-uuid-1"

	mock.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
		WithArgs(uid).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// Expect search query
	mock.ExpectQuery(`SELECT id, email, full_name, avatar_url FROM users WHERE \(email ILIKE \$1 OR full_name ILIKE \$1\) AND id != \$2 LIMIT 20`).
		WithArgs("%john%", userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "email", "full_name", "avatar_url"}).
			AddRow("u2", "john@example.com", "John Doe", ""))

	req := httptest.NewRequest("GET", "/api/users/search?q=john", nil)
	token := &auth.Token{UID: uid}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.SearchUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var resp []map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "john@example.com", resp[0]["email"])

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

	uid := "firebase-uid-1"
	userUUID := "user-uuid-1"

	// Mock user UUID lookup
	mock.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
		WithArgs(uid).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// Mock list connections
	mock.ExpectQuery("SELECT c.id, c.requester_id, c.receiver_id, c.status").
		WithArgs(userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "requester_id", "receiver_id", "status", "requester_email", "receiver_email"}).
			AddRow("conn-1", "user-uuid-2", userUUID, "pending", "user2@example.com", "user1@example.com"))

	req := httptest.NewRequest("GET", "/api/connections", nil)
	token := &auth.Token{UID: uid}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
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

	uid := "firebase-uid-1"
	userUUID := "user-uuid-1"
	connID := "conn-1"

	// Mock user UUID lookup
	mock.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
		WithArgs(uid).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// Mock update connection
	mock.ExpectExec("UPDATE connections SET status = 'accepted'").
		WithArgs(connID, userUUID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	req := httptest.NewRequest("PUT", "/api/connections/"+connID, nil)

	// Setup Chi context for URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", connID)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)

	token := &auth.Token{UID: uid}
	ctx = context.WithValue(ctx, middleware.UserContextKey, token)
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

	uid := "firebase-uid-1"
	userUUID := "user-uuid-1"
	connID := "conn-1"

	// Mock user UUID lookup
	mock.ExpectQuery("SELECT id FROM users WHERE firebase_uid=").
		WithArgs(uid).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// Mock delete connection
	mock.ExpectExec("DELETE FROM connections").
		WithArgs(connID, userUUID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	req := httptest.NewRequest("DELETE", "/api/connections/"+connID, nil)

	// Setup Chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", connID)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)

	token := &auth.Token{UID: uid}
	ctx = context.WithValue(ctx, middleware.UserContextKey, token)
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
