package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotificationHandler_RegisterDevice(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockNotificationService)
		h := NewNotificationHandler(mockService)
		userID := uuid.New()

		reqBody := `{"token": "fcm-token", "device_type": "android"}`
		req := httptest.NewRequest("POST", "/notifications/register", createBody(reqBody))

		// Inject test user key
		ctx := context.WithValue(req.Context(), TestUserKey, userID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &auth.Token{UID: "test-uid"})
		req = req.WithContext(ctx)

		mockService.On("RegisterDevice", mock.Anything, userID.String(), "fcm-token", "android").Return(nil)

		w := httptest.NewRecorder()

		h.RegisterDevice(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "registered", resp["status"])
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := new(MockNotificationService)
		h := NewNotificationHandler(mockService)

		reqBody := `{"token": "fcm-token"}`
		req := httptest.NewRequest("POST", "/notifications/register", createBody(reqBody))
		// No context injection

		w := httptest.NewRecorder()

		// We need to ensure GetUserUUID fails or the middleware context is missing.
		// If middleware key is missing, handler panics on `r.Context().Value(middleware.UserContextKey).(*auth.Token)`
		// So we must provide middleware key.
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &auth.Token{UID: "test-uid"})
		req = req.WithContext(ctx)

		// GetUserUUID will fail because DB is not initialized and no TestUserKey

		h.RegisterDevice(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
