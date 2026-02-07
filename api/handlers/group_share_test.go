package handlers

import (
	"context"
	"encoding/json"
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
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShareItemToGroup(t *testing.T) {
	t.Run("Success_Note", func(t *testing.T) {
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
		mockDB.ExpectQuery(`SELECT user_id, title FROM notes WHERE id = \$1 AND deleted_at IS NULL`).
			WithArgs(noteID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id", "title"}).AddRow(userID.String(), "My Note"))

		// 3. Create share
		mockDB.ExpectExec(`INSERT INTO group_shares`).
			WithArgs(groupID.String(), noteID.String(), userID, "Check this out").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// 4. Fetch details for notification
		mockDB.ExpectQuery(`SELECT name FROM groups WHERE id = \$1`).
			WithArgs(groupID.String()).
			WillReturnRows(mockDB.NewRows([]string{"name"}).AddRow("Test Group"))

		mockDB.ExpectQuery(`SELECT COALESCE\(username, email\) FROM users WHERE id = \$1`).
			WithArgs(userID.String()).
			WillReturnRows(mockDB.NewRows([]string{"display_name"}).AddRow("Test User"))

		mockDB.ExpectQuery(`SELECT user_id FROM group_members WHERE group_id = \$1 AND user_id != \$2`).
			WithArgs(groupID.String(), userID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id"}).AddRow(memberID.String()))

		mockNotify.On("SendNotification", mock.Anything, memberID.String(), "New Shared Item", "Test User shared \"My Note\" in Test Group", map[string]string{
			"type":        "note_share",
			"group_id":    groupID.String(),
			"resource_id": noteID.String(),
		}).Return(nil)

		handler.ShareItemToGroup(w, req)

		time.Sleep(50 * time.Millisecond) // Allow goroutine

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, true, resp["success"])

		assert.NoError(t, mockDB.ExpectationsWereMet())
		mockNotify.AssertExpectations(t)
	})

	t.Run("Success_VersePack_Public", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		mockNotify := new(MockNotificationService)
		handler := NewGroupShareHandler(mockDB, mockNotify)

		userID := uuid.New()
		groupID := uuid.New()
		packID := uuid.New()
		ownerID := uuid.New() // Different owner

		// Setup Request
		reqBody := `{"verse_pack_id": "` + packID.String() + `", "comment": "Great pack"}`
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

		// 2. Verify pack (Public)
		ownerIDStr := ownerID.String()
		mockDB.ExpectQuery(`SELECT user_id, title, is_public FROM verse_packs WHERE id = \$1`).
			WithArgs(packID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id", "title", "is_public"}).AddRow(&ownerIDStr, "Public Pack", true))

		// 3. Check existing share (assume not exists for insert path)
		mockDB.ExpectQuery(`SELECT id FROM group_shares WHERE group_id = \$1 AND verse_pack_id = \$2`).
			WithArgs(groupID.String(), packID.String()).
			WillReturnError(pgx.ErrNoRows)

		// 4. Create share
		mockDB.ExpectExec(`INSERT INTO group_shares`).
			WithArgs(groupID.String(), packID.String(), userID, "Great pack").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// 5. Notifications (simplified expectations)
		mockDB.ExpectQuery(`SELECT name FROM groups`).
			WithArgs(groupID.String()).
			WillReturnRows(mockDB.NewRows([]string{"name"}).AddRow("Group"))
		mockDB.ExpectQuery(`SELECT COALESCE`).
			WithArgs(userID.String()).
			WillReturnRows(mockDB.NewRows([]string{"display_name"}).AddRow("User"))
		mockDB.ExpectQuery(`SELECT user_id FROM group_members`).
			WithArgs(groupID.String(), userID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id"})) // No other members

		handler.ShareItemToGroup(w, req)
		time.Sleep(50 * time.Millisecond)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Success_VersePack_Owner", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, new(MockNotificationService))

		userID := uuid.New()
		groupID := uuid.New()
		packID := uuid.New()

		reqBody := `{"verse_pack_id": "` + packID.String() + `", "comment": "My pack"}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID.String()+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		ownerIDStr := userID.String()
		mockDB.ExpectQuery(`SELECT user_id, title, is_public FROM verse_packs WHERE id = \$1`).
			WithArgs(packID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id", "title", "is_public"}).AddRow(&ownerIDStr, "My Pack", false))

		// Check existing (exists, update)
		existingID := uuid.New()
		mockDB.ExpectQuery(`SELECT id FROM group_shares WHERE group_id = \$1 AND verse_pack_id = \$2`).
			WithArgs(groupID.String(), packID.String()).
			WillReturnRows(mockDB.NewRows([]string{"id"}).AddRow(existingID.String()))

		mockDB.ExpectExec(`UPDATE group_shares`).
			WithArgs(existingID.String(), groupID.String(), "My pack").
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		// Notifications (fail fast or empty)
		mockDB.ExpectQuery(`SELECT name FROM groups`).
			WithArgs(groupID.String()).
			WillReturnError(errors.New("db err"))

		mockDB.ExpectQuery(`SELECT COALESCE`).
			WithArgs(userID.String()).
			WillReturnError(errors.New("db err"))

		mockDB.ExpectQuery(`SELECT user_id FROM group_members`).
			WithArgs(groupID.String(), userID.String()).
			WillReturnError(errors.New("db err"))

		handler.ShareItemToGroup(w, req)
		time.Sleep(50 * time.Millisecond)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Error_NotMember", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)
		userID := uuid.New()
		groupID := uuid.New()
		noteID := uuid.New()

		reqBody := `{"note_id": "` + noteID.String() + `", "comment": ""}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID.String()+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(false))

		handler.ShareItemToGroup(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Error_NoteNotOwned", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)
		userID := uuid.New()
		groupID := uuid.New()
		noteID := uuid.New()
		otherUser := uuid.New()

		reqBody := `{"note_id": "` + noteID.String() + `", "comment": ""}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID.String()+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		mockDB.ExpectQuery(`SELECT user_id, title FROM notes`).
			WithArgs(noteID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id", "title"}).AddRow(otherUser.String(), "Not Mine"))

		handler.ShareItemToGroup(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Error_VersePackNotFound", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)
		userID := uuid.New()
		groupID := uuid.New()
		packID := uuid.New()

		reqBody := `{"verse_pack_id": "` + packID.String() + `", "comment": ""}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID.String()+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		mockDB.ExpectQuery(`SELECT user_id, title, is_public FROM verse_packs`).
			WithArgs(packID.String()).
			WillReturnError(pgx.ErrNoRows)

		handler.ShareItemToGroup(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Error_VersePackForbidden", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)
		userID := uuid.New()
		groupID := uuid.New()
		packID := uuid.New()
		otherUser := uuid.New()

		reqBody := `{"verse_pack_id": "` + packID.String() + `", "comment": ""}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID.String()+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		ownerIDStr := otherUser.String()
		mockDB.ExpectQuery(`SELECT user_id, title, is_public FROM verse_packs`).
			WithArgs(packID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id", "title", "is_public"}).AddRow(&ownerIDStr, "Private Pack", false))

		handler.ShareItemToGroup(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestListGroupShares(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)

		userID := uuid.New()
		groupID := uuid.New()
		shareID := uuid.New()
		noteID := uuid.New()

		req := httptest.NewRequest("GET", "/api/groups/"+groupID.String()+"/shares", nil)

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

		// 2. Query shares
		noteIDStr := noteID.String()
		var versePackIDStr *string = nil
		var subtitle *string = nil

		rows := mockDB.NewRows([]string{"id", "group_id", "note_id", "verse_pack_id", "title", "subtitle", "display_name", "shared_at", "comment", "type"}).
			AddRow(shareID.String(), groupID.String(), &noteIDStr, versePackIDStr, "My Note", subtitle, "User Name", time.Now(), "Comment", "note")

		mockDB.ExpectQuery(`SELECT gs.id, gs.group_id`).
			WithArgs(groupID.String()).
			WillReturnRows(rows)

		handler.ListGroupShares(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp []SharedItemResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "My Note", resp[0].Title)

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Empty", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)
		userID := uuid.New()
		groupID := uuid.New()

		req := httptest.NewRequest("GET", "/api/groups/"+groupID.String()+"/shares", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(`SELECT EXISTS`).WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		mockDB.ExpectQuery(`SELECT gs.id, gs.group_id`).WithArgs(groupID.String()).
			WillReturnRows(mockDB.NewRows([]string{"id", "group_id", "note_id", "verse_pack_id", "title", "subtitle", "display_name", "shared_at", "comment", "type"}))

		handler.ListGroupShares(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []SharedItemResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 0) // Should be empty list or null, assuming implementation handles it
	})

	t.Run("Error_DB", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)
		userID := uuid.New()
		groupID := uuid.New()

		req := httptest.NewRequest("GET", "/api/groups/"+groupID.String()+"/shares", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(`SELECT EXISTS`).WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		mockDB.ExpectQuery(`SELECT gs.id, gs.group_id`).WithArgs(groupID.String()).
			WillReturnError(errors.New("db error"))

		handler.ListGroupShares(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Error_NotMember", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)
		userID := uuid.New()
		groupID := uuid.New()

		req := httptest.NewRequest("GET", "/api/groups/"+groupID.String()+"/shares", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(`SELECT EXISTS`).WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(false))

		handler.ListGroupShares(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGetSharedItemDetails(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)

		userID := uuid.New()
		groupID := uuid.New()
		shareID := uuid.New()
		noteID := uuid.New()

		req := httptest.NewRequest("GET", "/api/groups/"+groupID.String()+"/shares/"+shareID.String(), nil)

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		rctx.URLParams.Add("shareId", shareID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		// 1. Verify membership
		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		// 2. Query detail
		noteIDStr := noteID.String()
		var versePackIDStr *string = nil
		content := map[string]interface{}{"text": "hello"}
		var subtitle *string = nil

		mockDB.ExpectQuery(`SELECT gs.id, gs.group_id`).
			WithArgs(shareID.String(), groupID.String()).
			WillReturnRows(mockDB.NewRows([]string{"id", "group_id", "note_id", "verse_pack_id", "title", "subtitle", "content", "display_name", "shared_at", "comment", "type"}).
				AddRow(shareID.String(), groupID.String(), &noteIDStr, versePackIDStr, "Title", subtitle, content, "User", time.Now(), "Comment", "note"))

		handler.GetSharedItemDetails(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp SharedItemResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Title", resp.Title)
		assert.Equal(t, map[string]interface{}{"text": "hello"}, resp.Content)

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)
		userID := uuid.New()
		groupID := uuid.New()
		shareID := uuid.New()

		req := httptest.NewRequest("GET", "/api/groups/"+groupID.String()+"/shares/"+shareID.String(), nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		rctx.URLParams.Add("shareId", shareID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		mockDB.ExpectQuery(`SELECT gs.id, gs.group_id`).
			WithArgs(shareID.String(), groupID.String()).
			WillReturnError(pgx.ErrNoRows)

		handler.GetSharedItemDetails(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)
		userID := uuid.New()
		groupID := uuid.New()
		shareID := uuid.New()

		req := httptest.NewRequest("GET", "/api/groups/"+groupID.String()+"/shares/"+shareID.String(), nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID.String())
		rctx.URLParams.Add("shareId", shareID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		mockDB.ExpectQuery(`SELECT gs.id, gs.group_id`).
			WithArgs(shareID.String(), groupID.String()).
			WillReturnError(errors.New("db error"))

		handler.GetSharedItemDetails(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}
