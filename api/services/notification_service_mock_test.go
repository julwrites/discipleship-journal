package services

import (
	"context"
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
