package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"firebase.google.com/go/v4/auth"
	"discipleship_journal_api/middleware"
	"github.com/google/uuid"
)

func TestNotificationHandler_RegisterDevice(t *testing.T) {
	// Mock NotificationService
	mockNotif := new(MockNotificationServiceWithMock)
	h := NewNotificationHandler(mockNotif)

	testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	testCases := []struct {
		name           string
		requestBody    string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:        "Success",
			requestBody: `{"token": "fcm-token-123", "device_type": "android"}`,
			setupMock: func() {
				// Expects UUID string, not Firebase UID, because Handler calls GetUserUUID then calls service with UUID string.
				mockNotif.On("RegisterDevice", mock.Anything, testUUID.String(), "fcm-token-123", "android").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Invalid Request",
			requestBody: `{"token": ""}`,
			setupMock:   func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			req := httptest.NewRequest("POST", "/notifications/device", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Inject Auth Token
			dummyToken := &auth.Token{UID: "firebase-uid-123"}
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, dummyToken)

			// Inject User UUID (since handler calls GetUserUUID)
			ctx = context.WithValue(ctx, TestUserKey, testUUID)

			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.RegisterDevice(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			mockNotif.AssertExpectations(t)
		})
	}
}
