package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockNotificationService_RegisterDevice(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := new(MockNotificationService)
		ctx := context.Background()
		userID := "user-123"
		token := "token-123"
		deviceType := "android"

		m.On("RegisterDevice", ctx, userID, token, deviceType).Return(nil)

		err := m.RegisterDevice(ctx, userID, token, deviceType)
		assert.NoError(t, err)
		m.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		m := new(MockNotificationService)
		ctx := context.Background()
		userID := "user-123"
		token := "token-123"
		deviceType := "android"
		expectedErr := errors.New("service error")

		m.On("RegisterDevice", ctx, userID, token, deviceType).Return(expectedErr)

		err := m.RegisterDevice(ctx, userID, token, deviceType)
		assert.Equal(t, expectedErr, err)
		m.AssertExpectations(t)
	})

	t.Run("SendNotification_Success", func(t *testing.T) {
		m := new(MockNotificationService)
		ctx := context.Background()
		userID := "user-123"
		title := "Title"
		body := "Body"
		data := map[string]string{"key": "val"}

		// Use mock.Anything to utilize the import
		m.On("SendNotification", ctx, userID, title, body, mock.Anything).Return(nil)

		err := m.SendNotification(ctx, userID, title, body, data)
		assert.NoError(t, err)
		m.AssertExpectations(t)
	})

	t.Run("SendNotification_Error", func(t *testing.T) {
		m := new(MockNotificationService)
		ctx := context.Background()
		userID := "user-123"
		title := "Title"
		body := "Body"
		data := map[string]string{"key": "val"}
		expectedErr := errors.New("service error")

		m.On("SendNotification", ctx, userID, title, body, data).Return(expectedErr)

		err := m.SendNotification(ctx, userID, title, body, data)
		assert.Equal(t, expectedErr, err)
		m.AssertExpectations(t)
	})
}
