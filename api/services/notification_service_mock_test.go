package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockNotificationService(t *testing.T) {
	service := NewMockNotificationService()

	t.Run("RegisterDevice", func(t *testing.T) {
		err := service.RegisterDevice(context.Background(), "user-id", "token", "android")
		assert.NoError(t, err)
	})

	t.Run("SendNotification", func(t *testing.T) {
		err := service.SendNotification(context.Background(), "user-id", "Title", "Body", nil)
		assert.NoError(t, err)
	})
}

func TestMockNotificationService_Error(t *testing.T) {
	service := NewMockNotificationService()
	service.RegisterDeviceFunc = func(ctx context.Context, userID, token, deviceType string) error {
		return errors.New("mock error")
	}
	service.SendNotificationFunc = func(ctx context.Context, userID, title, body string, data map[string]string) error {
		return errors.New("mock error")
	}

	t.Run("RegisterDevice", func(t *testing.T) {
		err := service.RegisterDevice(context.Background(), "user-id", "token", "android")
		assert.Error(t, err)
		assert.Equal(t, "mock error", err.Error())
	})

	t.Run("SendNotification", func(t *testing.T) {
		err := service.SendNotification(context.Background(), "user-id", "Title", "Body", nil)
		assert.Error(t, err)
		assert.Equal(t, "mock error", err.Error())
	})
}

func TestMockNotificationService_NilFuncs(t *testing.T) {
	service := &MockNotificationService{} // nil funcs

	t.Run("RegisterDevice", func(t *testing.T) {
		err := service.RegisterDevice(context.Background(), "user-id", "token", "android")
		assert.NoError(t, err)
	})

	t.Run("SendNotification", func(t *testing.T) {
		err := service.SendNotification(context.Background(), "user-id", "Title", "Body", nil)
		assert.NoError(t, err)
	})
}
