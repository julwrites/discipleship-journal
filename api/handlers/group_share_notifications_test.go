package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestShareItemToGroup_NotificationErrors(t *testing.T) {
	t.Run("DB_Member_Query_Error", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		mockNotify := new(MockNotificationService)
		handler := NewGroupShareHandler(mockDB, mockNotify)

		userID := uuid.New()
		groupID := uuid.New()
		noteID := uuid.New()

		// Setup Request
		reqBody := `{"note_id": "` + noteID.String() + `", "comment": "Check this out"}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID.String()+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// Mock Context with User
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)

		// Setup Chi Context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		// 1. Verify membership
		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		// 2. Verify ownership
		mockDB.ExpectQuery(`SELECT user_id, title FROM notes`).
			WithArgs(noteID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id", "title"}).AddRow(userID.String(), "My Note"))

		// 3. Create share
		mockDB.ExpectExec(`INSERT INTO group_shares`).
			WithArgs(groupID.String(), noteID.String(), userID, "Check this out").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// 4. Notification Logic - DB Error on group name
		// Implementation tries to fetch group name. If fails, defaults.
		mockDB.ExpectQuery(`SELECT name FROM groups`).
			WithArgs(groupID.String()).
			WillReturnError(errors.New("db error"))

		// Fetch sharer name
		mockDB.ExpectQuery(`SELECT COALESCE`).
			WithArgs(userID.String()).
			WillReturnError(errors.New("db error"))

		// Fetch members -> Error here should abort notification loop
		mockDB.ExpectQuery(`SELECT user_id FROM group_members`).
			WithArgs(groupID.String(), userID.String()).
			WillReturnError(errors.New("critical db error"))

		// No notification should be sent
		// No expectations on mockNotify

		handler.ShareItemToGroup(w, req)

		time.Sleep(50 * time.Millisecond) // Allow goroutine

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
		mockNotify.AssertExpectations(t)
	})

	t.Run("SendNotification_Error", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		mockNotify := new(MockNotificationService)
		handler := NewGroupShareHandler(mockDB, mockNotify)

		userID := uuid.New()
		groupID := uuid.New()
		noteID := uuid.New()
		memberID := uuid.New()

		// Setup Request
		reqBody := `{"note_id": "` + noteID.String() + `", "comment": "Check this out"}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID.String()+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		// 1. Verify membership
		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		// 2. Verify ownership
		mockDB.ExpectQuery(`SELECT user_id, title FROM notes`).
			WithArgs(noteID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id", "title"}).AddRow(userID.String(), "My Note"))

		// 3. Create share
		mockDB.ExpectExec(`INSERT INTO group_shares`).
			WithArgs(groupID.String(), noteID.String(), userID, "Check this out").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// 4. Notification Logic
		mockDB.ExpectQuery(`SELECT name FROM groups`).
			WithArgs(groupID.String()).
			WillReturnRows(mockDB.NewRows([]string{"name"}).AddRow("Group"))

		mockDB.ExpectQuery(`SELECT COALESCE`).
			WithArgs(userID.String()).
			WillReturnRows(mockDB.NewRows([]string{"display_name"}).AddRow("User"))

		mockDB.ExpectQuery(`SELECT user_id FROM group_members`).
			WithArgs(groupID.String(), userID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id"}).AddRow(memberID.String()))

		// Expect SendMulticastNotification to fail
		mockNotify.On("SendMulticastNotification", mock.Anything, []string{memberID.String()}, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("notification service error"))

		handler.ShareItemToGroup(w, req)

		time.Sleep(50 * time.Millisecond) // Allow goroutine

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
		mockNotify.AssertExpectations(t)
	})
}
