package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"database/sql"
	"discipleship_journal_api/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSendConnectionRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = db

	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mockNotification := services.NewMockNotificationService()
	handler := NewConnectionHandler(db, mockNotification)

	requesterUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	receiverEmail := "user2@example.com"
	receiverUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	connID := "00000000-0000-0000-0000-000000000003"

	// Mock receiver UUID lookup by email
	mock.ExpectQuery("SELECT id FROM users WHERE email =").
		WithArgs(receiverEmail).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receiverUUID.String()))

	// Mock insert connection
	mock.ExpectQuery("INSERT INTO connections").
		WithArgs(requesterUUID, receiverUUID.String()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(connID))

	// Expect synchronous requester name lookup
	mock.ExpectQuery("SELECT username FROM users WHERE id").
		WithArgs(requesterUUID).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("RequesterUser"))

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

func TestSendConnectionRequest_Errors(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := NewConnectionHandler(db, nil)

	t.Run("MissingReceiver", func(t *testing.T) {
		reqBody := ConnectionRequest{}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewBuffer(bodyBytes))
		req = req.WithContext(context.WithValue(req.Context(), TestUserKey, uuid.New()))
		w := httptest.NewRecorder()
		handler.SendConnectionRequest(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("SelfConnection", func(t *testing.T) {
		userID := uuid.New()

		// Mock receiver lookup to return same ID
		mockDB.ExpectQuery("SELECT id FROM users WHERE email =").
			WithArgs("me@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID.String()))

		reqBody := ConnectionRequest{ReceiverEmail: "me@example.com"}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewBuffer(bodyBytes))
		req = req.WithContext(context.WithValue(req.Context(), TestUserKey, userID))
		w := httptest.NewRecorder()
		handler.SendConnectionRequest(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ReceiverNotFound", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT id FROM users WHERE email =").
			WithArgs("missing@example.com").
			WillReturnError(sql.ErrNoRows)

		reqBody := ConnectionRequest{ReceiverEmail: "missing@example.com"}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewBuffer(bodyBytes))
		req = req.WithContext(context.WithValue(req.Context(), TestUserKey, uuid.New()))
		w := httptest.NewRecorder()
		handler.SendConnectionRequest(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSearchUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = db

	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mockNotification := services.NewMockNotificationService()
	handler := NewConnectionHandler(db, mockNotification)

	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Expect search query
	username := "johndoe"
	mock.ExpectQuery(`SELECT u.id, u.email, u.username.*`).
		WithArgs("johndoe", userUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "is_connected"}).
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

func TestSearchUsers_Errors(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := NewConnectionHandler(db, nil)

	t.Run("ShortQuery", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/search?q=ab", nil)
		w := httptest.NewRecorder()
		handler.SearchUsers(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		userUUID := uuid.New()
		mockDB.ExpectQuery("SELECT u.id").
			WillReturnError(assert.AnError)

		req := httptest.NewRequest("GET", "/api/users/search?q=error", nil)
		req = req.WithContext(context.WithValue(req.Context(), TestUserKey, userUUID))
		w := httptest.NewRecorder()
		handler.SearchUsers(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestListConnections(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = db

	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mockNotification := services.NewMockNotificationService()
	handler := NewConnectionHandler(db, mockNotification)

	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	connID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	requesterUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	// Mock list connections
	requesterUsername := "user2"
	receiverUsername := "user1"
	mock.ExpectQuery("SELECT c.id, c.requester_id, c.receiver_id, c.status").
		WithArgs(userUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "requester_id", "receiver_id", "status", "requester_email", "receiver_email", "requester_username", "receiver_username"}).
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

func TestListConnections_Error(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := NewConnectionHandler(db, nil)
	userUUID := uuid.New()

	mockDB.ExpectQuery("SELECT c.id").WillReturnError(assert.AnError)

	req := httptest.NewRequest("GET", "/api/connections", nil)
	req = req.WithContext(context.WithValue(req.Context(), TestUserKey, userUUID))
	w := httptest.NewRecorder()
	handler.ListConnections(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAcceptConnectionRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = db

	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mockNotification := services.NewMockNotificationService()
	handler := NewConnectionHandler(db, mockNotification)

	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	connID := "00000000-0000-0000-0000-000000000003"

	// Mock update connection
	mock.ExpectExec("UPDATE connections SET status = 'accepted'").
		WithArgs(connID, userUUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

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

func TestAcceptConnectionRequest_Errors(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := NewConnectionHandler(db, nil)
	userUUID := uuid.New()
	connID := "conn-1"

	t.Run("DBError", func(t *testing.T) {
		mockDB.ExpectExec("UPDATE connections").WithArgs(connID, userUUID).WillReturnError(assert.AnError)
		req := httptest.NewRequest("PUT", "/api/connections/"+connID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", connID)
		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
		ctx = context.WithValue(ctx, TestUserKey, userUUID)
		w := httptest.NewRecorder()
		handler.AcceptConnectionRequest(w, req.WithContext(ctx))
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockDB.ExpectExec("UPDATE connections").WithArgs(connID, userUUID).WillReturnResult(sqlmock.NewResult(1, 0))
		req := httptest.NewRequest("PUT", "/api/connections/"+connID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", connID)
		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
		ctx = context.WithValue(ctx, TestUserKey, userUUID)
		w := httptest.NewRecorder()
		handler.AcceptConnectionRequest(w, req.WithContext(ctx))
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteConnectionRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = db

	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mockNotification := services.NewMockNotificationService()
	handler := NewConnectionHandler(db, mockNotification)

	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	connID := "00000000-0000-0000-0000-000000000003"

	// Mock delete connection
	mock.ExpectExec("DELETE FROM connections").
		WithArgs(connID, userUUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

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

func TestDeleteConnectionRequest_Errors(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := NewConnectionHandler(db, nil)
	userUUID := uuid.New()
	connID := "conn-1"

	t.Run("DBError", func(t *testing.T) {
		mockDB.ExpectExec("DELETE FROM connections").WithArgs(connID, userUUID).WillReturnError(assert.AnError)
		req := httptest.NewRequest("DELETE", "/api/connections/"+connID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", connID)
		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
		ctx = context.WithValue(ctx, TestUserKey, userUUID)
		w := httptest.NewRecorder()
		handler.DeleteConnectionRequest(w, req.WithContext(ctx))
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockDB.ExpectExec("DELETE FROM connections").WithArgs(connID, userUUID).WillReturnResult(sqlmock.NewResult(1, 0))
		req := httptest.NewRequest("DELETE", "/api/connections/"+connID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", connID)
		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
		ctx = context.WithValue(ctx, TestUserKey, userUUID)
		w := httptest.NewRecorder()
		handler.DeleteConnectionRequest(w, req.WithContext(ctx))
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSendConnectionRequest_MoreErrors(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB

	_ = db

	_ = mockDB
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mockNotif := services.NewMockNotificationService()
	handler := NewConnectionHandler(db, mockNotif)
	requesterUUID := uuid.New()

	t.Run("Unauthorized", func(t *testing.T) {
		reqBody := ConnectionRequest{ReceiverEmail: "some@email.com"}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewBuffer(bodyBytes))
		// No user in context
		w := httptest.NewRecorder()
		handler.SendConnectionRequest(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("ReceiverLookupError", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT id FROM users WHERE email =").
			WithArgs("some@email.com").
			WillReturnError(errors.New("db connection error"))

		reqBody := ConnectionRequest{ReceiverEmail: "some@email.com"}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewBuffer(bodyBytes))
		req = req.WithContext(context.WithValue(req.Context(), TestUserKey, requesterUUID))
		w := httptest.NewRecorder()
		handler.SendConnectionRequest(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("InsertConflict", func(t *testing.T) {
		receiverUUID := uuid.New()
		mockDB.ExpectQuery("SELECT id FROM users WHERE email =").
			WithArgs("some@email.com").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receiverUUID.String()))

		mockDB.ExpectQuery("INSERT INTO connections").
			WithArgs(requesterUUID, receiverUUID.String()).
			WillReturnError(errors.New("duplicate key value violates unique constraint"))

		reqBody := ConnectionRequest{ReceiverEmail: "some@email.com"}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewBuffer(bodyBytes))
		req = req.WithContext(context.WithValue(req.Context(), TestUserKey, requesterUUID))
		w := httptest.NewRecorder()
		handler.SendConnectionRequest(w, req)
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("RequesterNameLookupFailure", func(t *testing.T) {
		receiverUUID := uuid.New()
		connID := "new-conn-id"

		mockDB.ExpectQuery("SELECT id FROM users WHERE email =").
			WithArgs("some@email.com").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receiverUUID.String()))

		mockDB.ExpectQuery("INSERT INTO connections").
			WithArgs(requesterUUID, receiverUUID.String()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(connID))

		// Fail name lookup
		mockDB.ExpectQuery("SELECT username FROM users WHERE id").
			WithArgs(requesterUUID).
			WillReturnError(errors.New("db error"))

		// Notification should still be sent, but with fallback name "Someone"
		// We can't verify the arguments to mockNotif easily because it's a manual mock in services package
		// that just returns nil. But we can verify code execution path doesn't panic and returns 201.

		reqBody := ConnectionRequest{ReceiverEmail: "some@email.com"}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewBuffer(bodyBytes))
		req = req.WithContext(context.WithValue(req.Context(), TestUserKey, requesterUUID))
		w := httptest.NewRecorder()
		handler.SendConnectionRequest(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
