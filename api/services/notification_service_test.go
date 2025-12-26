package services

import (
	"context"
	"testing"

	"firebase.google.com/go/v4/messaging"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMessagingClient
type MockMessagingClient struct {
	mock.Mock
}

func (m *MockMessagingClient) SendEachForMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
	args := m.Called(ctx, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*messaging.BatchResponse), args.Error(1)
}

func TestNotificationService_RegisterDevice(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	service := NewNotificationService(mockDB, nil)

	ctx := context.Background()
	userID := "user-uuid"
	token := "fcm-token"
	deviceType := "android"

	mockDB.ExpectExec("INSERT INTO user_devices").
		WithArgs(userID, token, deviceType).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = service.RegisterDevice(ctx, userID, token, deviceType)
	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestNotificationService_SendNotification(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	mockMsg := new(MockMessagingClient)
	service := NewNotificationService(mockDB, mockMsg)

	ctx := context.Background()
	userID := "user-uuid"
	title := "Title"
	body := "Body"
	data := map[string]string{"key": "value"}

	// 1. Get tokens
	mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"fcm_token"}).AddRow("token-1"))

	// 2. Send message
	mockMsg.On("SendEachForMulticast", ctx, mock.MatchedBy(func(msg *messaging.MulticastMessage) bool {
		return len(msg.Tokens) == 1 && msg.Tokens[0] == "token-1" &&
			msg.Notification.Title == title && msg.Notification.Body == body &&
			msg.Data["key"] == "value"
	})).Return(&messaging.BatchResponse{SuccessCount: 1, FailureCount: 0}, nil)

	err = service.SendNotification(ctx, userID, title, body, data)
	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockMsg.AssertExpectations(t)
}

func TestNotificationService_SendNotification_CleanupInvalidTokens(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	mockMsg := new(MockMessagingClient)
	service := NewNotificationService(mockDB, mockMsg)

	ctx := context.Background()
	userID := "user-uuid"

	// 1. Get tokens
	mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"fcm_token"}).AddRow("token-valid").AddRow("token-invalid"))

	// 2. Send message
	mockMsg.On("SendEachForMulticast", ctx, mock.Anything).Return(&messaging.BatchResponse{
		SuccessCount: 1,
		FailureCount: 1,
		Responses: []*messaging.SendResponse{
			{Success: true},
			{Success: false, Error: assert.AnError}, // Corresponds to token-invalid
		},
	}, nil)

	// 3. Cleanup
	mockDB.ExpectExec("DELETE FROM user_devices").
		WithArgs(userID, []string{"token-invalid"}).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = service.SendNotification(ctx, userID, "T", "B", nil)
	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockMsg.AssertExpectations(t)
}
