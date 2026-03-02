package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"discipleship_journal_api/middleware"

	"firebase.google.com/go/v4/auth"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestShareItemToGroup_Coverage(t *testing.T) {
	t.Run("Missing_IDs", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		_ = db

		_ = dbMock
		assert.NoError(t, err)
		defer db.Close()

		handler := NewGroupShareHandler(db, nil)
		userID := uuid.New().String()
		groupID := uuid.New().String()

		reqBody := `{"comment": "test"}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.ShareItemToGroup(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Either note_id or verse_pack_id is required")
	})

	t.Run("Membership_Check_DB_Error", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		_ = db

		_ = dbMock
		assert.NoError(t, err)
		defer db.Close()

		handler := NewGroupShareHandler(db, nil)
		userID := uuid.New().String()
		groupID := uuid.New().String()
		noteID := uuid.New().String()

		dbMock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnError(errors.New("db error"))

		reqBody := `{"note_id": "` + noteID + `"}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.ShareItemToGroup(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("VersePack_Info_DB_Error", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		_ = db

		_ = dbMock
		assert.NoError(t, err)
		defer db.Close()

		handler := NewGroupShareHandler(db, nil)
		userID := uuid.New().String()
		groupID := uuid.New().String()
		packID := uuid.New().String()

		dbMock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		dbMock.ExpectQuery("SELECT user_id, title, is_public FROM verse_packs").
			WithArgs(packID).
			WillReturnError(errors.New("db error"))

		reqBody := `{"verse_pack_id": "` + packID + `"}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.ShareItemToGroup(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Existing_Share_Check_DB_Error", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		_ = db

		_ = dbMock
		assert.NoError(t, err)
		defer db.Close()

		handler := NewGroupShareHandler(db, nil)
		userID := uuid.New().String()
		groupID := uuid.New().String()
		packID := uuid.New().String()

		dbMock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		dbMock.ExpectQuery("SELECT user_id, title, is_public FROM verse_packs").
			WithArgs(packID).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "title", "is_public"}).AddRow(userID, "Title", false))

		dbMock.ExpectQuery("SELECT id FROM group_shares").
			WithArgs(groupID, packID).
			WillReturnError(errors.New("db error"))

		// Note: The logic in handler falls back to INSERT if err != nil?
		// Actually: `if err == nil { UPDATE } else { INSERT }`
		// So if DB error occurs, it tries INSERT.
		// If INSERT fails (due to same DB error potentially or not), then it returns error.
		// Let's expect INSERT and make it fail.

		dbMock.ExpectExec("INSERT INTO group_shares").
			WithArgs(groupID, packID, userID, "").
			WillReturnError(errors.New("insert error"))

		reqBody := `{"verse_pack_id": "` + packID + `", "comment": ""}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.ShareItemToGroup(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Update_Share_DB_Error", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		_ = db

		_ = dbMock
		assert.NoError(t, err)
		defer db.Close()

		handler := NewGroupShareHandler(db, nil)
		userID := uuid.New().String()
		groupID := uuid.New().String()
		packID := uuid.New().String()
		shareID := uuid.New().String()

		dbMock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		dbMock.ExpectQuery("SELECT user_id, title, is_public FROM verse_packs").
			WithArgs(packID).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "title", "is_public"}).AddRow(userID, "Title", false))

		dbMock.ExpectQuery("SELECT id FROM group_shares").
			WithArgs(groupID, packID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(shareID))

		dbMock.ExpectExec("UPDATE group_shares").
			WithArgs(shareID, groupID, "comment").
			WillReturnError(errors.New("update error"))

		reqBody := `{"verse_pack_id": "` + packID + `", "comment": "comment"}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID+"/shares", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.ShareItemToGroup(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestListGroupShares_Coverage(t *testing.T) {
	t.Run("Membership_Check_DB_Error", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		_ = db

		_ = dbMock
		assert.NoError(t, err)
		defer db.Close()

		handler := NewGroupShareHandler(db, nil)
		userID := uuid.New().String()
		groupID := uuid.New().String()

		dbMock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnError(errors.New("db error"))

		req := httptest.NewRequest("GET", "/api/groups/"+groupID+"/shares", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.ListGroupShares(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetSharedItemDetails_Coverage(t *testing.T) {
	t.Run("Membership_Check_DB_Error", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		_ = db

		_ = dbMock
		assert.NoError(t, err)
		defer db.Close()

		handler := NewGroupShareHandler(db, nil)
		userID := uuid.New().String()
		groupID := uuid.New().String()
		shareID := uuid.New().String()

		dbMock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnError(errors.New("db error"))

		req := httptest.NewRequest("GET", "/api/groups/"+groupID+"/shares/"+shareID, nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
		ctx = context.WithValue(ctx, TestUserKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", groupID)
		rctx.URLParams.Add("shareId", shareID)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.GetSharedItemDetails(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
