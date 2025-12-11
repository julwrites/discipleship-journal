package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"discipleship_journal_api/database"
	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestShareNoteToGroup(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	database.DB = mock

	userID := "uid-1"
	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	groupID := "group-1"
	noteID := "note-1"

	reqBody := ShareNoteRequest{
		NoteID:  noteID,
		Comment: "Check this out",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/groups/"+groupID+"/shares", bytes.NewBuffer(body))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", groupID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	token := &auth.Token{UID: userID}
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, token))

	// 1. GetUserUUID
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// 2. Check Member
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(groupID, userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	// 3. Verify Ownership
	mock.ExpectQuery("SELECT user_id FROM notes").
		WithArgs(noteID).
		WillReturnRows(pgxmock.NewRows([]string{"user_id"}).AddRow(userUUID.String()))

	// 4. Insert Share
	mock.ExpectExec("INSERT INTO group_shares").
		WithArgs(groupID, noteID, userUUID, reqBody.Comment).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	w := httptest.NewRecorder()
	ShareNoteToGroup(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListGroupShares(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	database.DB = mock

	userID := "uid-1"
	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	groupID := "group-1"

	req := httptest.NewRequest("GET", "/groups/"+groupID+"/shares", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", groupID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	token := &auth.Token{UID: userID}
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, token))

	// 1. GetUserUUID
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(userUUID))

	// 2. Check Member
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(groupID, userUUID).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	// 3. List Shares
	// now := time.Now()

	mock.ExpectQuery("SELECT gs.id, gs.group_id").
		WithArgs(groupID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "group_id", "note_id", "title", "display_name", "shared_at", "comment"}).
			AddRow("share-1", groupID, "note-1", "My Note", "User 1", time.Now(), "Comment"))

	w := httptest.NewRecorder()
	ListGroupShares(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
