package handlers

import (
	"context"
)

// MockNotificationService is a mock implementation of services.NotificationService
type MockNotificationService struct {}

func (m *MockNotificationService) RegisterDevice(ctx context.Context, userID, token, deviceType string) error {
	return nil
}

func (m *MockNotificationService) SendNotification(ctx context.Context, userID, title, body string, data map[string]string) error {
	return nil
}
