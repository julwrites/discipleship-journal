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

	"discipleship_journal_api/models"
	"discipleship_journal_api/services"
	"github.com/go-chi/chi/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetNotes_TestUserKey(t *testing.T) {
	testUserID := "test-user-id"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		_, noteServiceMock, handler := setupTest(t)

		serviceNotes := []services.Note{
			{
				ID:        "note-1",
				UserID:    testUserID,
				Title:     "Title 1",
				Content:   json.RawMessage(`{"text": "content"}`),
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		noteServiceMock.On("GetNotes", mock.Anything, testUserID, 1, 20, services.NoteFilter{}).Return(serviceNotes, 1, nil)

		req := httptest.NewRequest("GET", "/api/notes", nil)
		// Use TestUserKey instead of Firebase Token
		ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.GetNotes(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		noteServiceMock.AssertExpectations(t)
	})
}

func TestGetNotes_NilService(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer dbMock.Close()

	// Handler with nil service
	handler := NewNoteHandler(dbMock, nil)

	req := httptest.NewRequest("GET", "/api/notes", nil)
	w := httptest.NewRecorder()
	handler.GetNotes(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteNote_NilService(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer dbMock.Close()
	handler := NewNoteHandler(dbMock, nil)

	req := httptest.NewRequest("DELETE", "/api/notes/1", nil)
	w := httptest.NewRecorder()
	handler.DeleteNote(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateNote_NilService(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer dbMock.Close()
	handler := NewNoteHandler(dbMock, nil)

	req := httptest.NewRequest("POST", "/api/notes", nil)
	w := httptest.NewRecorder()
	handler.CreateNote(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateNote_NilService(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer dbMock.Close()
	handler := NewNoteHandler(dbMock, nil)

	req := httptest.NewRequest("PUT", "/api/notes/1", nil)
	w := httptest.NewRecorder()
	handler.UpdateNote(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetNote_NilService(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer dbMock.Close()
	handler := NewNoteHandler(dbMock, nil)

	req := httptest.NewRequest("GET", "/api/notes/1", nil)
	w := httptest.NewRecorder()
	handler.GetNote(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetTags_NilService(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer dbMock.Close()
	handler := NewNoteHandler(dbMock, nil)

	req := httptest.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()
	handler.GetTags(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateTag_NilService(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer dbMock.Close()
	handler := NewNoteHandler(dbMock, nil)

	req := httptest.NewRequest("POST", "/api/tags", nil)
	w := httptest.NewRecorder()
	handler.CreateTag(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteTag_NilService(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer dbMock.Close()
	handler := NewNoteHandler(dbMock, nil)

	req := httptest.NewRequest("DELETE", "/api/tags/1", nil)
	w := httptest.NewRecorder()
	handler.DeleteTag(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetTags_Success(t *testing.T) {
	testUserID := "test-user-id"
	_, noteServiceMock, handler := setupTest(t)

	tags := []services.Tag{
		{ID: "tag-1", Name: "Tag 1", UserID: testUserID},
	}

	noteServiceMock.On("GetUserTags", mock.Anything, testUserID).Return(tags, nil)

	req := httptest.NewRequest("GET", "/api/tags", nil)
	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetTags(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []services.Tag
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Tag 1", resp[0].Name)
}

func TestGetTags_Error(t *testing.T) {
	testUserID := "test-user-id"
	_, noteServiceMock, handler := setupTest(t)

	noteServiceMock.On("GetUserTags", mock.Anything, testUserID).Return(nil, errors.New("service error"))

	req := httptest.NewRequest("GET", "/api/tags", nil)
	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetTags(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateTag_Success(t *testing.T) {
	testUserID := "test-user-id"
	_, noteServiceMock, handler := setupTest(t)

	tagName := "New Tag"
	tag := &services.Tag{ID: "tag-1", Name: tagName, UserID: testUserID}

	noteServiceMock.On("CreateTag", mock.Anything, testUserID, tagName).Return(tag, nil)

	body := `{"name": "New Tag"}`
	req := httptest.NewRequest("POST", "/api/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTag(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp services.Tag
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, tagName, resp.Name)
}

func TestCreateTag_ValidationFail(t *testing.T) {
	testUserID := "test-user-id"
	_, _, handler := setupTest(t)

	// Empty name should fail required validation
	body := `{"name": ""}`
	req := httptest.NewRequest("POST", "/api/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTag(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTag_ServiceError(t *testing.T) {
	testUserID := "test-user-id"
	_, noteServiceMock, handler := setupTest(t)

	tagName := "New Tag"
	noteServiceMock.On("CreateTag", mock.Anything, testUserID, tagName).Return(nil, errors.New("service error"))

	body := `{"name": "New Tag"}`
	req := httptest.NewRequest("POST", "/api/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTag(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteTag_Success(t *testing.T) {
	testUserID := "test-user-id"
	tagID := "tag-1"
	_, noteServiceMock, handler := setupTest(t)

	noteServiceMock.On("DeleteTag", mock.Anything, testUserID, tagID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/tags/"+tagID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", tagID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTag(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTag_NotFound(t *testing.T) {
	testUserID := "test-user-id"
	tagID := "tag-1"
	_, noteServiceMock, handler := setupTest(t)

	noteServiceMock.On("DeleteTag", mock.Anything, testUserID, tagID).Return(models.ErrNotFound)

	req := httptest.NewRequest("DELETE", "/api/tags/"+tagID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", tagID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTag(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteTag_ServiceError(t *testing.T) {
	testUserID := "test-user-id"
	tagID := "tag-1"
	_, noteServiceMock, handler := setupTest(t)

	noteServiceMock.On("DeleteTag", mock.Anything, testUserID, tagID).Return(errors.New("service error"))

	req := httptest.NewRequest("DELETE", "/api/tags/"+tagID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", tagID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTag(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteNote_TestUserKey(t *testing.T) {
	testUserID := "test-user-id"
	noteID := "note-1"
	_, noteServiceMock, handler := setupTest(t)

	noteServiceMock.On("DeleteNote", mock.Anything, testUserID, noteID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/notes/"+noteID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", noteID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteNote(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateNote_TestUserKey(t *testing.T) {
	testUserID := "test-user-id"
	_, noteServiceMock, handler := setupTest(t)

	noteServiceMock.On("CreateNote", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&services.Note{ID: "new"}, nil)

	body := `{"title": "Test", "content": {}}`
	req := httptest.NewRequest("POST", "/api/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateNote(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateNote_TestUserKey(t *testing.T) {
	testUserID := "test-user-id"
	noteID := "note-1"
	_, noteServiceMock, handler := setupTest(t)

	noteServiceMock.On("UpdateNote", mock.Anything, testUserID, noteID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	body := `{"title": "Test", "content": {}}`
	req := httptest.NewRequest("PUT", "/api/notes/"+noteID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", noteID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.UpdateNote(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetNote_TestUserKey(t *testing.T) {
	testUserID := "test-user-id"
	noteID := "note-1"
	_, noteServiceMock, handler := setupTest(t)

	noteServiceMock.On("GetNote", mock.Anything, testUserID, noteID).Return(&services.Note{ID: noteID}, nil)

	req := httptest.NewRequest("GET", "/api/notes/"+noteID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", noteID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetNote(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserUUIDFromContext_NilDB(t *testing.T) {
	handler := &NoteHandler{db: nil}
	_, err := handler.getUserUUID(context.Background(), "firebase-uid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection is nil")
}

func TestGetNotes_ServiceError(t *testing.T) {
	testUserID := "test-user-id"
	_, noteServiceMock, handler := setupTest(t)

	noteServiceMock.On("GetNotes", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything).Return(nil, 0, errors.New("service error"))

	req := httptest.NewRequest("GET", "/api/notes", nil)
	ctx := context.WithValue(req.Context(), TestUserKey, testUserID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetNotes(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
