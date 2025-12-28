package handlers

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockNotificationServiceWithMock is a struct that embeds mock.Mock and implements services.NotificationService
// This allows us to use testify/mock assertions
type MockNotificationServiceWithMock struct {
	mock.Mock
}

func (m *MockNotificationServiceWithMock) RegisterDevice(ctx context.Context, userID, token, deviceType string) error {
	args := m.Called(ctx, userID, token, deviceType)
	return args.Error(0)
}

func (m *MockNotificationServiceWithMock) SendNotification(ctx context.Context, userID, title, body string, data map[string]string) error {
	args := m.Called(ctx, userID, title, body, data)
	return args.Error(0)
}
