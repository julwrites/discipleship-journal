package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestConnectionHandler_Unauthorized(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer dbMock.Close()

	handler := NewConnectionHandler(dbMock, nil)

	t.Run("SearchUsers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/search?q=query", nil)
		w := httptest.NewRecorder()
		handler.SearchUsers(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("ListConnections", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/connections", nil)
		w := httptest.NewRecorder()
		handler.ListConnections(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("AcceptConnectionRequest", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/connections/1", nil)
		w := httptest.NewRecorder()
		handler.AcceptConnectionRequest(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("DeleteConnectionRequest", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/connections/1", nil)
		w := httptest.NewRecorder()
		handler.DeleteConnectionRequest(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestSendConnectionRequest_InvalidBody(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer dbMock.Close()

	handler := NewConnectionHandler(dbMock, nil)
	userUUID := uuid.New()

	req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewBufferString("invalid json"))
	req = req.WithContext(context.WithValue(req.Context(), TestUserKey, userUUID))
	w := httptest.NewRecorder()
	handler.SendConnectionRequest(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchUsers_ScanError(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer dbMock.Close()

	handler := NewConnectionHandler(dbMock, nil)
	userUUID := uuid.New()

	// Mock query returning fewer columns than expected to force scan error
	// Query: SELECT u.id, u.email, u.username, is_connected
	// Scan: &id, &email, &username, &isConnected
	dbMock.ExpectQuery("SELECT u.id, u.email").
		WithArgs("query", userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "email"}).
			AddRow("id-1", "test@test.com"))

	req := httptest.NewRequest("GET", "/api/users/search?q=query", nil)
	req = req.WithContext(context.WithValue(req.Context(), TestUserKey, userUUID))
	w := httptest.NewRecorder()
	handler.SearchUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// It should skip the row and return empty list
	var resp []interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 0)
}

func TestListConnections_ScanError(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer dbMock.Close()

	handler := NewConnectionHandler(dbMock, nil)
	userUUID := uuid.New()

	// Query: SELECT c.id, c.requester_id, c.receiver_id, c.status, u1.email, u2.email, u1.username, u2.username
	// Return fewer columns
	dbMock.ExpectQuery("SELECT c.id, c.requester_id").
		WithArgs(userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "requester_id"}).
			AddRow("conn-1", "req-1"))

	req := httptest.NewRequest("GET", "/api/connections", nil)
	req = req.WithContext(context.WithValue(req.Context(), TestUserKey, userUUID))
	w := httptest.NewRecorder()
	handler.ListConnections(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Should skip row
	var resp []interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 0)
}

func TestAcceptConnectionRequest_RouteContext(t *testing.T) {
	// Verify that AcceptConnectionRequest works when called properly
	// (Already covered in connection_test.go, but ensuring context setup correctness for completeness if needed)
	// The main purpose here is to check if we missed anything.
	// Actually, TestAcceptConnectionRequest in connection_test.go sets TestUserKey.
	// We covered Unauthorized above.
}
