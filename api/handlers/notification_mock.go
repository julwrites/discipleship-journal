package handlers

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// MockNotificationService is a mock implementation of services.NotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) RegisterDevice(ctx context.Context, userID, token, deviceType string) error {
	args := m.Called(ctx, userID, token, deviceType)
	return args.Error(0)
}

func (m *MockNotificationService) SendNotification(ctx context.Context, userID, title, body string, data map[string]string) error {
	args := m.Called(ctx, userID, title, body, data)
	return args.Error(0)
}
