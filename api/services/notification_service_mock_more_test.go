package services

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMockNotificationService_Multicast(t *testing.T) {
	mockService := NewMockNotificationService()
	err := mockService.SendMulticastNotification(context.Background(), []string{"user1"}, "title", "body", nil)
	assert.NoError(t, err)

	mockService.SendMulticastNotificationFunc = nil
	err = mockService.SendMulticastNotification(context.Background(), []string{"user1"}, "title", "body", nil)
	assert.NoError(t, err)
}

func TestMockNotificationService_New(t *testing.T) {
    mockService := NewMockNotificationService()
    assert.NotNil(t, mockService)
}
