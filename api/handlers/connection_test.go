package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestSendConnectionRequest(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	handler := NewConnectionHandler(mock)

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

func TestListConnections(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	handler := NewConnectionHandler(mock)

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
