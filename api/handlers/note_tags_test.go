package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/models"
	"discipleship_journal_api/services"

	"firebase.google.com/go/v4/auth"
	"github.com/DATA-DOG/go-sqlmock"
	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTagsHandler(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"

	t.Run("success", func(t *testing.T) {
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		tags := []services.Tag{
			{ID: "tag-1", Name: "Bible"},
			{ID: "tag-2", Name: "Prayer"},
		}

		noteServiceMock.On("GetUserTags", mock.Anything, userUUID).Return(tags, nil)

		req := httptest.NewRequest("GET", "/api/tags", nil)
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.GetTags(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []services.Tag
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, "Bible", resp[0].Name)

		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("GetUserTags", mock.Anything, userUUID).Return(nil, assert.AnError)

		req := httptest.NewRequest("GET", "/api/tags", nil)
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.GetTags(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		noteServiceMock.AssertExpectations(t)
	})
}

func TestCreateTagHandler_ServiceError(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	tagName := "New Tag"

	validReq := CreateTagRequest{
		Name: tagName,
	}

	db, dbMock, noteServiceMock, handler := setupTest(t)

	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

	noteServiceMock.On("CreateTag", mock.Anything, userUUID, tagName).Return(nil, errors.New("service error"))

	body, _ := json.Marshal(validReq)
	req := httptest.NewRequest("POST", "/api/tags", bytes.NewBuffer(body))
	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTag(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteTagHandler_ServiceError(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	tagID := "tag-123"

	db, dbMock, noteServiceMock, handler := setupTest(t)

	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

	noteServiceMock.On("DeleteTag", mock.Anything, userUUID, tagID).Return(errors.New("service error"))

	req := httptest.NewRequest("DELETE", "/api/tags/"+tagID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", tagID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTag(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetTags_UserNotFound(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	db, dbMock, noteServiceMock, handler := setupTest(t)
	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnError(errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/tags", nil)
	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetTags(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	noteServiceMock.AssertNotCalled(t, "GetUserTags")
}

func TestCreateTag_UserNotFound(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	db, dbMock, noteServiceMock, handler := setupTest(t)
	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnError(errors.New("db error"))

	req := httptest.NewRequest("POST", "/api/tags", bytes.NewBufferString("{}"))
	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTag(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	noteServiceMock.AssertNotCalled(t, "CreateTag")
}

func TestDeleteTag_UserNotFound(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	db, dbMock, noteServiceMock, handler := setupTest(t)
	_ = db
	_ = db
	_ = dbMock
	defer db.Close()

	dbMock.ExpectQuery("SELECT id FROM users").
		WithArgs(firebaseUID).
		WillReturnError(errors.New("db error"))

	req := httptest.NewRequest("DELETE", "/api/tags/1", nil)
	token := &auth.Token{UID: firebaseUID}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTag(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	noteServiceMock.AssertNotCalled(t, "DeleteTag")
}

func TestCreateTagHandler(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	tagName := "New Tag"

	validReq := CreateTagRequest{
		Name: tagName,
	}

	t.Run("success", func(t *testing.T) {
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		createdTag := &services.Tag{
			ID:   "tag-new",
			Name: tagName,
		}

		noteServiceMock.On("CreateTag", mock.Anything, userUUID, tagName).Return(createdTag, nil)

		body, _ := json.Marshal(validReq)
		req := httptest.NewRequest("POST", "/api/tags", bytes.NewBuffer(body))
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.CreateTag(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp services.Tag
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "tag-new", resp.ID)
		assert.Equal(t, tagName, resp.Name)

		noteServiceMock.AssertExpectations(t)
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled db expectations: %s", err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		db, dbMock, _, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		// Missing name
		invalidReq := CreateTagRequest{}
		body, _ := json.Marshal(invalidReq)
		req := httptest.NewRequest("POST", "/api/tags", bytes.NewBuffer(body))
		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		// Mock user lookup is still called before validation in the handler
		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		w := httptest.NewRecorder()
		handler.CreateTag(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteTagHandler(t *testing.T) {
	firebaseUID := "firebase-uid-123"
	userUUID := "user-uuid-123"
	tagID := "tag-123"

	t.Run("success", func(t *testing.T) {
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("DeleteTag", mock.Anything, userUUID, tagID).Return(nil)

		req := httptest.NewRequest("DELETE", "/api/tags/"+tagID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", tagID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.DeleteTag(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"success": true}`, w.Body.String())
		noteServiceMock.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		db, dbMock, noteServiceMock, handler := setupTest(t)
		_ = db
		_ = db
		_ = dbMock
		defer db.Close()

		dbMock.ExpectQuery("SELECT id FROM users").
			WithArgs(firebaseUID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userUUID))

		noteServiceMock.On("DeleteTag", mock.Anything, userUUID, tagID).Return(models.ErrNotFound)

		req := httptest.NewRequest("DELETE", "/api/tags/"+tagID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", tagID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		token := &auth.Token{UID: firebaseUID}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.DeleteTag(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		noteServiceMock.AssertExpectations(t)
	})
}
