package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"firebase.google.com/go/v4/auth"
	"discipleship_journal_api/middleware"
	"github.com/google/uuid"
)

func TestGroupHandler_CreateGroup(t *testing.T) {
	mockDB, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close(context.Background())

	mockNotif := new(MockNotificationServiceWithMock)
	h := NewGroupHandler(mockDB, mockNotif)

	testCases := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:        "Success",
			requestBody: `{"name": "Bible Study", "description": "Weekly study"}`,
			setupMock: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("INSERT INTO groups").
					WithArgs("Bible Study", "Weekly study", uuid.MustParse("00000000-0000-0000-0000-000000000001")).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("group-123"))
				mockDB.ExpectExec("INSERT INTO group_members").
					WithArgs("group-123", uuid.MustParse("00000000-0000-0000-0000-000000000001")).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mockDB.ExpectCommit()
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "Invalid Request",
			requestBody: `{"name": ""}`, // too short
			setupMock:   func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "DB Error",
			requestBody: `{"name": "Bible Study", "description": "Weekly study"}`,
			setupMock: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("INSERT INTO groups").
					WithArgs("Bible Study", "Weekly study", uuid.MustParse("00000000-0000-0000-0000-000000000001")).
					WillReturnError(pgx.ErrTxClosed) // Simulate error
				mockDB.ExpectRollback()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			req := httptest.NewRequest("POST", "/groups", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Inject test user via TestUserKey
			testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
			ctx := context.WithValue(req.Context(), TestUserKey, testUUID)

			dummyToken := &auth.Token{UID: "firebase-uid-123"}
			ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)

			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.CreateGroup(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestGroupHandler_ListMyGroups(t *testing.T) {
	mockDB, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close(context.Background())

	mockNotif := new(MockNotificationServiceWithMock)
	h := NewGroupHandler(mockDB, mockNotif)

	t.Run("Success", func(t *testing.T) {
		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		mockDB.ExpectQuery(`SELECT g.id, g.name, g.description, g.created_by, gm.role`).
			WithArgs(testUUID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_by", "role"}).
				AddRow("g1", "Group 1", "Desc 1", "creator1", "admin"))

		req := httptest.NewRequest("GET", "/groups", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.ListMyGroups(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var groups []GroupResponse
		err := json.Unmarshal(w.Body.Bytes(), &groups)
		assert.NoError(t, err)
		assert.Len(t, groups, 1)
		assert.Equal(t, "Group 1", groups[0].Name)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_SearchGroups(t *testing.T) {
	mockDB, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close(context.Background())

	mockNotif := new(MockNotificationServiceWithMock)
	h := NewGroupHandler(mockDB, mockNotif)

	testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	t.Run("Success", func(t *testing.T) {
		mockDB.ExpectQuery(`SELECT g.id, g.name, g.description, g.created_by`).
			WithArgs("%Bible%", testUUID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_by", "role"}).
				AddRow("g1", "Bible Study", "Desc", "creator", "member"))

		req := httptest.NewRequest("GET", "/groups/search?q=Bible", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.SearchGroups(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Query Too Short", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/groups/search?q=Bi", nil)
		w := httptest.NewRecorder()
		h.SearchGroups(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_JoinGroup(t *testing.T) {
	mockDB, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close(context.Background())

	mockNotif := new(MockNotificationServiceWithMock)
	h := NewGroupHandler(mockDB, mockNotif)

	testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	t.Run("Success", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs("g1", testUUID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		mockDB.ExpectExec("INSERT INTO group_members").
			WithArgs("g1", testUUID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		req := httptest.NewRequest("POST", "/groups/g1/join", nil)

		// Setup chi context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "g1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.JoinGroup(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Already Member", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs("g1", testUUID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		req := httptest.NewRequest("POST", "/groups/g1/join", nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "g1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.JoinGroup(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_GetGroupMembers(t *testing.T) {
	mockDB, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close(context.Background())

	mockNotif := new(MockNotificationServiceWithMock)
	h := NewGroupHandler(mockDB, mockNotif)

	testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	t.Run("Success", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs("g1", testUUID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		mockDB.ExpectQuery("SELECT gm.user_id").
			WithArgs("g1").
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "display_name", "email", "role", "joined_at"}).
				AddRow("u1", "User 1", "u1@example.com", "admin", time.Now()))

		req := httptest.NewRequest("GET", "/groups/g1/members", nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "g1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetGroupMembers(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_AddGroupMember(t *testing.T) {
	mockDB, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close(context.Background())

	mockNotif := new(MockNotificationServiceWithMock)
	h := NewGroupHandler(mockDB, mockNotif)

	testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	t.Run("Success - Admin Adds Member", func(t *testing.T) {
		// 1. Verify admin
		mockDB.ExpectQuery("SELECT role FROM group_members").
			WithArgs("g1", testUUID).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("admin"))

		// 2. Insert member
		mockDB.ExpectExec("INSERT INTO group_members").
			WithArgs("g1", "target-user").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// 3. Get group name
		mockDB.ExpectQuery("SELECT name FROM groups").
			WithArgs("g1").
			WillReturnRows(pgxmock.NewRows([]string{"name"}).AddRow("Test Group"))

		// 4. Notification
		mockNotif.On("SendNotification", mock.Anything, "target-user", "Group Invitation", mock.Anything, mock.Anything).Return(nil).Maybe()

		reqBody := `{"user_id": "target-user"}`
		req := httptest.NewRequest("POST", "/groups/g1/members", strings.NewReader(reqBody))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "g1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.AddGroupMember(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Give goroutine a moment to run
		time.Sleep(10 * time.Millisecond)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Forbidden - Not Admin", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT role FROM group_members").
			WithArgs("g1", testUUID).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("member"))

		reqBody := `{"user_id": "target-user"}`
		req := httptest.NewRequest("POST", "/groups/g1/members", strings.NewReader(reqBody))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "g1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.AddGroupMember(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_RemoveGroupMember(t *testing.T) {
	mockDB, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close(context.Background())

	mockNotif := new(MockNotificationServiceWithMock)
	h := NewGroupHandler(mockDB, mockNotif)

	testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	t.Run("Success", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT role FROM group_members").
			WithArgs("g1", testUUID).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("admin"))

		mockDB.ExpectExec("DELETE FROM group_members").
			WithArgs("g1", "target-user").
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		req := httptest.NewRequest("DELETE", "/groups/g1/members/target-user", nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "g1")
		rctx.URLParams.Add("userId", "target-user")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.RemoveGroupMember(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}
