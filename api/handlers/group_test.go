package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNotificationService for testing
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) RegisterDevice(ctx context.Context, userID, token, deviceType string) error {
	args := m.Called(ctx, userID, token, deviceType)
	return args.Error(0)
}

func (m *MockNotificationService) SendNotification(ctx context.Context, userID, title, body string, data map[string]string) error {
	args := m.Called(ctx, userID, title, body, data)
	return args.Error(0)
}

func TestGroupHandler_CreateGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		mockNotif := new(MockNotificationService)
		h := NewGroupHandler(mockDB, mockNotif)

		userID := uuid.New()
		groupID := uuid.New()

		mockDB.ExpectBegin()
		mockDB.ExpectQuery("INSERT INTO groups").
			WithArgs("New Group", "Description", userID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(groupID.String()))
		mockDB.ExpectExec("INSERT INTO group_members").
			WithArgs(groupID.String(), userID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
		mockDB.ExpectCommit()

		reqBody := `{"name": "New Group", "description": "Description"}`
		req := httptest.NewRequest("POST", "/groups", createBody(reqBody))

		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &auth.Token{UID: "test-uid"})
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.CreateGroup(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_ListMyGroups(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		h := NewGroupHandler(mockDB, nil)
		userID := uuid.New()

		mockDB.ExpectQuery("SELECT g.id, g.name").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_by", "role"}).
				AddRow("group-1", "Group 1", "Desc 1", "creator-1", "admin").
				AddRow("group-2", "Group 2", "Desc 2", "creator-2", "member"))

		req := httptest.NewRequest("GET", "/groups", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &auth.Token{UID: "test-uid"})
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.ListMyGroups(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var groups []GroupResponse
		err = json.Unmarshal(w.Body.Bytes(), &groups)
		assert.NoError(t, err)
		assert.Len(t, groups, 2)
		assert.Equal(t, "Group 1", groups[0].Name)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_SearchGroups(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		h := NewGroupHandler(mockDB, nil)
		userID := uuid.New()

		mockDB.ExpectQuery("SELECT g.id, g.name").
			WithArgs("%query%", userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_by", "role"}).
				AddRow("group-1", "Group Query", "Desc", "creator", ""))

		req := httptest.NewRequest("GET", "/groups/search?q=query", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &auth.Token{UID: "test-uid"})
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		h.SearchGroups(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var groups []GroupResponse
		err = json.Unmarshal(w.Body.Bytes(), &groups)
		assert.NoError(t, err)
		assert.Len(t, groups, 1)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_JoinGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		h := NewGroupHandler(mockDB, nil)
		userID := uuid.New()
		groupID := "group-id"

		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		mockDB.ExpectExec("INSERT INTO group_members").
			WithArgs(groupID, userID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		req := httptest.NewRequest("POST", "/groups/"+groupID+"/join", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &auth.Token{UID: "test-uid"})
		req = req.WithContext(ctx)

		// Add chi param
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		h.JoinGroup(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_LeaveGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		h := NewGroupHandler(mockDB, nil)
		userID := uuid.New()
		groupID := "group-id"

		mockDB.ExpectExec("DELETE FROM group_members").
			WithArgs(groupID, userID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		req := httptest.NewRequest("POST", "/groups/"+groupID+"/leave", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &auth.Token{UID: "test-uid"})
		req = req.WithContext(ctx)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		h.LeaveGroup(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_GetGroupMembers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		h := NewGroupHandler(mockDB, nil)
		userID := uuid.New()
		groupID := "group-id"

		mockDB.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		mockDB.ExpectQuery("SELECT gm.user_id").
			WithArgs(groupID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "display_name", "email", "role", "joined_at"}).
				AddRow("u1", "User 1", "u1@example.com", "admin", time.Now()))

		req := httptest.NewRequest("GET", "/groups/"+groupID+"/members", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &auth.Token{UID: "test-uid"})
		req = req.WithContext(ctx)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		h.GetGroupMembers(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_AddGroupMember(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		mockNotif := new(MockNotificationService)
		h := NewGroupHandler(mockDB, mockNotif)
		userID := uuid.New()
		groupID := "group-id"
		targetID := "target-uuid"

		mockDB.ExpectQuery("SELECT role").
			WithArgs(groupID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("admin"))

		mockDB.ExpectExec("INSERT INTO group_members").
			WithArgs(groupID, targetID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mockDB.ExpectQuery("SELECT name FROM groups").
			WithArgs(groupID).
			WillReturnRows(pgxmock.NewRows([]string{"name"}).AddRow("Test Group"))

		mockNotif.On("SendNotification", mock.Anything, targetID, "Group Invitation", mock.Anything, mock.Anything).Return(nil)

		reqBody := `{"user_id": "target-uuid"}`
		req := httptest.NewRequest("POST", "/groups/"+groupID+"/members", createBody(reqBody))
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &auth.Token{UID: "test-uid"})
		req = req.WithContext(ctx)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		h.AddGroupMember(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGroupHandler_RemoveGroupMember(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		h := NewGroupHandler(mockDB, nil)
		userID := uuid.New()
		groupID := "group-id"
		targetID := "target-uuid"

		mockDB.ExpectQuery("SELECT role").
			WithArgs(groupID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("admin"))

		mockDB.ExpectExec("DELETE FROM group_members").
			WithArgs(groupID, targetID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		req := httptest.NewRequest("DELETE", "/groups/"+groupID+"/members/"+targetID, nil)
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &auth.Token{UID: "test-uid"})
		req = req.WithContext(ctx)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		rctx.URLParams.Add("userId", targetID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		h.RemoveGroupMember(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}
