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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotificationHandler_RegisterDevice(t *testing.T) {
	testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	testCases := []struct {
		name           string
		requestBody    string
		token          *auth.Token
		skipUserKey    bool
		setupMock      func(mockNotif *MockNotificationServiceWithMock)
		expectedStatus int
	}{
		{
			name:        "Success",
			requestBody: `{"token": "fcm-token-123", "device_type": "android"}`,
			token:       &auth.Token{UID: "firebase-uid-123"},
			setupMock: func(mockNotif *MockNotificationServiceWithMock) {
				// Expects UUID string, not Firebase UID, because Handler calls GetUserUUID then calls service with UUID string.
				mockNotif.On("RegisterDevice", mock.Anything, testUUID.String(), "fcm-token-123", "android").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Request",
			requestBody:    `{"token": ""}`,
			token:          &auth.Token{UID: "firebase-uid-123"},
			setupMock:      func(mockNotif *MockNotificationServiceWithMock) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unauthorized - No Token",
			requestBody:    `{"token": "fcm-token-123"}`,
			token:          nil,
			setupMock:      func(mockNotif *MockNotificationServiceWithMock) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Unauthorized - GetUserUUID Error",
			requestBody:    `{"token": "fcm-token-123"}`,
			token:          &auth.Token{UID: "firebase-uid-123"},
			skipUserKey:    true, // This will cause GetUserUUID to fail because DB is not initialized
			setupMock:      func(mockNotif *MockNotificationServiceWithMock) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Service Error",
			requestBody:    `{"token": "fcm-token-123"}`,
			token:          &auth.Token{UID: "firebase-uid-123"},
			setupMock:      func(mockNotif *MockNotificationServiceWithMock) {
				mockNotif.On("RegisterDevice", mock.Anything, testUUID.String(), "fcm-token-123", "").
					Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock NotificationService
			mockNotif := new(MockNotificationServiceWithMock)
			h := NewNotificationHandler(mockNotif)

			tc.setupMock(mockNotif)

			req := httptest.NewRequest("POST", "/notifications/device", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			ctx := req.Context()

			// Inject Auth Token if present
			if tc.token != nil {
				ctx = context.WithValue(ctx, middleware.UserContextKey, tc.token)
			}

			// Inject User UUID (unless skipped to test GetUserUUID failure)
			if !tc.skipUserKey {
				ctx = context.WithValue(ctx, TestUserKey, testUUID)
			}

			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.RegisterDevice(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			mockNotif.AssertExpectations(t)
		})
	}
}
