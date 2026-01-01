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

func TestShareNoteToGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		mockNotify := new(MockNotificationServiceWithMock)
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
		mockDB.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM group_members WHERE group_id = \$1 AND user_id = \$2\)`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(true))

		// 2. Verify ownership
		mockDB.ExpectQuery(`SELECT user_id FROM notes WHERE id = \$1 AND deleted_at IS NULL`).
			WithArgs(noteID.String()).
			WillReturnRows(mockDB.NewRows([]string{"user_id"}).AddRow(userID.String()))

		// 3. Create share
		mockDB.ExpectExec(`INSERT INTO group_shares`).
			WithArgs(groupID.String(), noteID.String(), userID, "Check this out").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// 4. Fetch details for notification
		mockDB.ExpectQuery(`SELECT name FROM groups WHERE id = \$1`).
			WithArgs(groupID.String()).
			WillReturnRows(mockDB.NewRows([]string{"name"}).AddRow("Test Group"))

		mockDB.ExpectQuery(`SELECT display_name FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnRows(mockDB.NewRows([]string{"display_name"}).AddRow("Test User"))

		// 5. Get members for notification (async in handler, but pgxmock requires strict ordering/expectations)
		// Since the goroutine runs asynchronously, we might have a race condition in test expectations.
		// However, in typical unit tests, the goroutine might not run fast enough or we can't easily sync.
		// For this test, we can verify the HTTP response first.
		// Ideally, we should sync or structure code to allow testing.
		// Given the handler structure:
		// go func() { ... }()
		// The expectations for the async part might fail or not be met before the test ends.
		// We'll proceed without expecting the async queries for now, or we can use a sleep.
		// A better approach for testability is to mock the `go func` or make it synchronous in tests,
		// but since we can't change the handler code right now easily without refactoring,
		// we will focus on the main flow.
		// Update: pgxmock will complain if expectations are not met. If we add expectations for the async part,
		// we must ensure they run.
		// Let's TRY to add them and a small sleep, or just Mock the DB interactions for the async part too.
		// BUT: pgxmock isn't thread-safe for the SAME connection if used concurrently, but here it's likely sharing the pool.

		// Actually, since it spawns a goroutine using `h.db`, and `h.db` is our mock,
		// we can expect the queries if we wait.

		mockDB.ExpectQuery(`SELECT user_id FROM group_members WHERE group_id = \$1 AND user_id != \$2`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"user_id"}).AddRow(memberID.String()))

		mockNotify.On("SendNotification", mock.Anything, memberID.String(), "New Shared Note", "Test User shared a note in Test Group", map[string]string{
			"type":     "note_share",
			"group_id": groupID.String(),
			"note_id":  noteID.String(),
		}).Return(nil)

		handler.ShareNoteToGroup(w, req)

		// Give time for goroutine
		time.Sleep(50 * time.Millisecond)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
		mockNotify.AssertExpectations(t)
	})

	t.Run("NotMember", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)

		userID := uuid.New()
		groupID := uuid.New()
		noteID := uuid.New()

		reqBody := `{"note_id": "` + noteID.String() + `", "comment": "test"}`
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

		handler.ShareNoteToGroup(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("NoteNotFound", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		handler := NewGroupShareHandler(mockDB, nil)

		userID := uuid.New()
		groupID := uuid.New()
		noteID := uuid.New()

		reqBody := `{"note_id": "` + noteID.String() + `", "comment": "test"}`
		req := httptest.NewRequest("POST", "/api/groups/"+groupID.String()+"/shares", strings.NewReader(reqBody))

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

		mockDB.ExpectQuery(`SELECT user_id FROM notes`).
			WithArgs(noteID.String()).
			WillReturnError(pgx.ErrNoRows)

		handler.ShareNoteToGroup(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
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
		rows := mockDB.NewRows([]string{"id", "group_id", "note_id", "title", "display_name", "shared_at", "comment"}).
			AddRow(shareID.String(), groupID.String(), noteID.String(), "My Note", "User Name", time.Now(), "Comment")

		mockDB.ExpectQuery(`SELECT gs.id, gs.group_id, gs.note_id`).
			WithArgs(groupID.String()).
			WillReturnRows(rows)

		handler.ListGroupShares(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp []SharedNoteResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "My Note", resp[0].Title)

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("AccessDenied", func(t *testing.T) {
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

		mockDB.ExpectQuery(`SELECT EXISTS`).
			WithArgs(groupID.String(), userID).
			WillReturnRows(mockDB.NewRows([]string{"exists"}).AddRow(false))

		handler.ListGroupShares(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGetSharedNoteDetails(t *testing.T) {
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
		// Use map for content (jsonb) to avoid type mismatch in pgxmock
		content := map[string]interface{}{"text": "hello"}

		mockDB.ExpectQuery(`SELECT gs.id, gs.group_id, gs.note_id`).
			WithArgs(shareID.String(), groupID.String()).
			WillReturnRows(mockDB.NewRows([]string{"id", "group_id", "note_id", "title", "content", "display_name", "shared_at", "comment"}).
				AddRow(shareID.String(), groupID.String(), noteID.String(), "Title", content, "User", time.Now(), "Comment"))

		handler.GetSharedNoteDetails(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp SharedNoteResponse
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

		mockDB.ExpectQuery(`SELECT gs.id`).
			WithArgs(shareID.String(), groupID.String()).
			WillReturnError(pgx.ErrNoRows)

		handler.GetSharedNoteDetails(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}
