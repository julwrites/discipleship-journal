package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGroupHandler_CreateGroup(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    string
		setupMock      func(mockDB pgxmock.PgxConnIface)
		expectedStatus int
	}{
		{
			name:        "Success",
			requestBody: `{"name": "Bible Study", "description": "Weekly study"}`,
			setupMock: func(mockDB pgxmock.PgxConnIface) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("INSERT INTO groups").
					WithArgs("Bible Study", "Weekly study", pgxmock.AnyArg(), "group").
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("00000000-0000-0000-0000-000000000123"))
				mockDB.ExpectExec("INSERT INTO group_members").
					WithArgs("00000000-0000-0000-0000-000000000123", pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mockDB.ExpectCommit()
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid Request",
			requestBody:    `{"name": ""}`, // too short
			setupMock:      func(_ pgxmock.PgxConnIface) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "DB Error",
			requestBody: `{"name": "Bible Study", "description": "Weekly study"}`,
			setupMock: func(mockDB pgxmock.PgxConnIface) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("INSERT INTO groups").
					WithArgs("Bible Study", "Weekly study", pgxmock.AnyArg(), "group").
					WillReturnError(pgx.ErrTxClosed) // Simulate error
				mockDB.ExpectRollback()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewConn()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func() { _ = mockDB.Close(context.Background()) }()

			mockNotif := new(MockNotificationServiceWithMock)
			h := NewGroupHandler(mockDB, mockNotif)

			tc.setupMock(mockDB)

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
			mockNotif.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_ListMyGroups(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		desc := "Desc 1"
		mockDB.ExpectQuery(`SELECT g.id, g.name, g.description, g.created_by, g.type, gm.role`).
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_by", "type", "role"}).
				AddRow("00000000-0000-0000-0000-000000000001", "Group 1", &desc, "00000000-0000-0000-0000-000000000002", nil, "admin"))

		req := httptest.NewRequest("GET", "/groups", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.ListMyGroups(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var groups []GroupResponse
		err = json.Unmarshal(w.Body.Bytes(), &groups)
		assert.NoError(t, err)
		assert.Len(t, groups, 1)
		assert.Equal(t, "Group 1", groups[0].Name)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_SearchGroups(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		desc := "Desc"
		// Updated to match the actual query logic including COALESCE
		mockDB.ExpectQuery(`SELECT g.id, g.name, g.description, g.created_by, g.type,
		 COALESCE\(\(SELECT role FROM group_members WHERE group_id = g.id AND user_id = \$2\), ''\) as role
		 FROM groups g
		 WHERE g.name ILIKE \$1 AND \(g.type = 'group' OR g.type IS NULL\) LIMIT 20`).
			WithArgs("%Bible%", pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_by", "type", "role"}).
				AddRow("00000000-0000-0000-0000-000000000001", "Bible Study", &desc, "00000000-0000-0000-0000-000000000002", nil, "member"))

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
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		req := httptest.NewRequest("GET", "/groups/search?q=Bi", nil)
		w := httptest.NewRecorder()
		h.SearchGroups(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_JoinGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs("00000000-0000-0000-0000-000000000001", pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		mockDB.ExpectExec("INSERT INTO group_members").
			WithArgs("00000000-0000-0000-0000-000000000001", pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		req := httptest.NewRequest("POST", "/groups/00000000-0000-0000-0000-000000000001/join", nil)

		// Setup chi context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "00000000-0000-0000-0000-000000000001")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.JoinGroup(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"success": true}`, w.Body.String())
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Already Member", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs("00000000-0000-0000-0000-000000000001", pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		req := httptest.NewRequest("POST", "/groups/00000000-0000-0000-0000-000000000001/join", nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "00000000-0000-0000-0000-000000000001")
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

func TestGroupHandler_LeaveGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		mockDB.ExpectExec("DELETE FROM group_members").
			WithArgs("00000000-0000-0000-0000-000000000001", pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		req := httptest.NewRequest("POST", "/groups/00000000-0000-0000-0000-000000000001/leave", nil)

		// Setup chi context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "00000000-0000-0000-0000-000000000001")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.LeaveGroup(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"success": true}`, w.Body.String())
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Not Member", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		mockDB.ExpectExec("DELETE FROM group_members").
			WithArgs("00000000-0000-0000-0000-000000000001", pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		req := httptest.NewRequest("POST", "/groups/00000000-0000-0000-0000-000000000001/leave", nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "00000000-0000-0000-0000-000000000001")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.LeaveGroup(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_GetGroupMembers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs("00000000-0000-0000-0000-000000000001", pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		mockDB.ExpectQuery("SELECT gm.user_id").
			WithArgs("00000000-0000-0000-0000-000000000001").
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "display_name", "email", "role", "joined_at"}).
				AddRow("00000000-0000-0000-0000-000000000002", "User 1", "u1@example.com", "admin", time.Now()))

		req := httptest.NewRequest("GET", "/groups/00000000-0000-0000-0000-000000000001/members", nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "00000000-0000-0000-0000-000000000001")
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
	t.Run("Success - Admin Adds Member", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		// 1. Verify admin
		mockDB.ExpectQuery("SELECT role FROM group_members").
			WithArgs("00000000-0000-0000-0000-000000000001", pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("admin"))

		// 2. Verify Connection
		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs(testUUID, "00000000-0000-0000-0000-000000000002").
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		// 3. Insert member
		mockDB.ExpectExec("INSERT INTO group_members").
			WithArgs("00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// 3. Get group name
		mockDB.ExpectQuery("SELECT name FROM groups").
			WithArgs("00000000-0000-0000-0000-000000000001").
			WillReturnRows(pgxmock.NewRows([]string{"name"}).AddRow("Test Group"))

		// 4. Notification
		mockNotif.On("SendNotification", mock.Anything, "00000000-0000-0000-0000-000000000002", "Group Invitation", mock.Anything, mock.Anything).Return(nil).Maybe()

		reqBody := `{"user_id": "00000000-0000-0000-0000-000000000002"}`
		req := httptest.NewRequest("POST", "/groups/00000000-0000-0000-0000-000000000001/members", strings.NewReader(reqBody))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "00000000-0000-0000-0000-000000000001")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.AddGroupMember(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"success": true}`, w.Body.String())

		// Give goroutine a moment to run
		time.Sleep(10 * time.Millisecond)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Forbidden - Not Admin", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		mockDB.ExpectQuery("SELECT role FROM group_members").
			WithArgs("00000000-0000-0000-0000-000000000001", pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("member"))

		reqBody := `{"user_id": "00000000-0000-0000-0000-000000000002"}`
		req := httptest.NewRequest("POST", "/groups/00000000-0000-0000-0000-000000000001/members", strings.NewReader(reqBody))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "00000000-0000-0000-0000-000000000001")
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

	t.Run("Forbidden - Not Connected", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		targetUUID := "00000000-0000-0000-0000-000000000002"

		// 1. Verify admin
		mockDB.ExpectQuery("SELECT role FROM group_members").
			WithArgs("00000000-0000-0000-0000-000000000001", pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("admin"))

		// 2. Verify Connection (Fail)
		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs(testUUID, targetUUID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		reqBody := `{"user_id": "` + targetUUID + `"}`
		req := httptest.NewRequest("POST", "/groups/00000000-0000-0000-0000-000000000001/members", strings.NewReader(reqBody))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "00000000-0000-0000-0000-000000000001")
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

func TestGroupHandler_GetOrCreateDirectGroup(t *testing.T) {
	t.Run("Success - New Group", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		partnerUUID := "00000000-0000-0000-0000-000000000002"

		// 1. Connection check
		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs(testUUID, partnerUUID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		// 2. Check if exists
		mockDB.ExpectQuery(`SELECT g.id FROM groups g`).
			WithArgs(testUUID, partnerUUID).
			WillReturnError(pgx.ErrNoRows)

		// 3. Get Partner Info
		mockDB.ExpectQuery("SELECT COALESCE").
			WithArgs(partnerUUID).
			WillReturnRows(pgxmock.NewRows([]string{"name", "username", "email"}).AddRow("User 2", "user2", "u2@example.com"))

		// 4. Get My Info
		mockDB.ExpectQuery("SELECT COALESCE").
			WithArgs(testUUID).
			WillReturnRows(pgxmock.NewRows([]string{"name"}).AddRow("User 1"))

		// 5. Create Group
		mockDB.ExpectBegin()
		mockDB.ExpectQuery("INSERT INTO groups").
			WithArgs("Direct: User 1 & User 2", testUUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("00000000-0000-0000-0000-000000000099"))

		// 6. Add Members
		mockDB.ExpectExec("INSERT INTO group_members").
			WithArgs("00000000-0000-0000-0000-000000000099", testUUID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mockDB.ExpectExec("INSERT INTO group_members").
			WithArgs("00000000-0000-0000-0000-000000000099", partnerUUID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mockDB.ExpectCommit()

		reqBody := `{"partner_id": "` + partnerUUID + `"}`
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"id": "00000000-0000-0000-0000-000000000099"}`, w.Body.String())
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Success - Existing Group", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		partnerUUID := "00000000-0000-0000-0000-000000000002"

		// 1. Connection check
		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs(testUUID, partnerUUID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		// 2. Check if exists
		mockDB.ExpectQuery(`SELECT g.id FROM groups g`).
			WithArgs(testUUID, partnerUUID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("00000000-0000-0000-0000-000000000088"))

		reqBody := `{"partner_id": "` + partnerUUID + `"}`
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusOK, w.Code) // Note: Handler might return 200 here? Yes, 200 in code if found.
		// Actually the handler returns 200 if found.
		assert.JSONEq(t, `{"id": "00000000-0000-0000-0000-000000000088"}`, w.Body.String())
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Failure - Not Connected", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		partnerUUID := "00000000-0000-0000-0000-000000000003"

		// 1. Connection check
		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs(testUUID, partnerUUID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		reqBody := `{"partner_id": "` + partnerUUID + `"}`
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Failure - Self Partner", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		reqBody := `{"partner_id": "00000000-0000-0000-0000-000000000001"}`
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_RemoveGroupMember(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() { _ = mockDB.Close(context.Background()) }()

		mockNotif := new(MockNotificationServiceWithMock)
		h := NewGroupHandler(mockDB, mockNotif)

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		mockDB.ExpectQuery("SELECT role FROM group_members").
			WithArgs("00000000-0000-0000-0000-000000000001", pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("admin"))

		mockDB.ExpectExec("DELETE FROM group_members").
			WithArgs("00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002").
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		req := httptest.NewRequest("DELETE", "/groups/00000000-0000-0000-0000-000000000001/members/00000000-0000-0000-0000-000000000002", nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "00000000-0000-0000-0000-000000000001")
		rctx.URLParams.Add("userId", "00000000-0000-0000-0000-000000000002")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.RemoveGroupMember(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"success": true}`, w.Body.String())
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}
