package services

import (
	"context"
)

// MockNotificationService is a mock implementation of NotificationService for testing.
type MockNotificationService struct {
	RegisterDeviceFunc            func(ctx context.Context, userID, token, deviceType string) error
	SendNotificationFunc          func(ctx context.Context, userID, title, body string, data map[string]string) error
	SendMulticastNotificationFunc func(ctx context.Context, userIDs []string, title, body string, data map[string]string) error
}

func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{
		RegisterDeviceFunc: func(ctx context.Context, userID, token, deviceType string) error {
			return nil
		},
		SendNotificationFunc: func(ctx context.Context, userID, title, body string, data map[string]string) error {
			return nil
		},
		SendMulticastNotificationFunc: func(ctx context.Context, userIDs []string, title, body string, data map[string]string) error {
			return nil
		},
	}
}

func (m *MockNotificationService) RegisterDevice(ctx context.Context, userID, token, deviceType string) error {
	if m.RegisterDeviceFunc != nil {
		return m.RegisterDeviceFunc(ctx, userID, token, deviceType)
	}
	return nil
}

func (m *MockNotificationService) SendMulticastNotification(ctx context.Context, userIDs []string, title, body string, data map[string]string) error {
	if m.SendMulticastNotificationFunc != nil {
		return m.SendMulticastNotificationFunc(ctx, userIDs, title, body, data)
	}
	return nil
}

func (m *MockNotificationService) SendNotification(ctx context.Context, userID, title, body string, data map[string]string) error {
	if m.SendNotificationFunc != nil {
		return m.SendNotificationFunc(ctx, userID, title, body, data)
	}
	return nil
}
