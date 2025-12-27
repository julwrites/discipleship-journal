package services

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"firebase.google.com/go/v4/messaging"
)

// MockMessagingClient is a mock for MessagingClientInterface
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
	mockDB, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close(context.Background())

	mockMsgClient := new(MockMessagingClient)
	service := NewNotificationService(mockDB, mockMsgClient)

	t.Run("Success", func(t *testing.T) {
		mockDB.ExpectExec("INSERT INTO user_devices").
			WithArgs("user-123", "token-123", "ios").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := service.RegisterDevice(context.Background(), "user-123", "token-123", "ios")
		assert.NoError(t, err)
	})

	t.Run("DB Error", func(t *testing.T) {
		mockDB.ExpectExec("INSERT INTO user_devices").
			WithArgs("user-123", "token-123", "ios").
			WillReturnError(errors.New("db error"))

		err := service.RegisterDevice(context.Background(), "user-123", "token-123", "ios")
		assert.Error(t, err)
	})
}

func TestNotificationService_SendNotification(t *testing.T) {
	mockDB, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close(context.Background())

	mockMsgClient := new(MockMessagingClient)
	service := NewNotificationService(mockDB, mockMsgClient)

	t.Run("Success - Sent to devices", func(t *testing.T) {
		// 1. Get tokens
		mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
			WithArgs("user-123").
			WillReturnRows(pgxmock.NewRows([]string{"fcm_token"}).AddRow("token-1"))

		// 2. Send message
		mockMsgClient.On("SendEachForMulticast", mock.Anything, mock.MatchedBy(func(msg *messaging.MulticastMessage) bool {
			return len(msg.Tokens) == 1 && msg.Tokens[0] == "token-1" && msg.Notification.Title == "Test Title"
		})).Return(&messaging.BatchResponse{SuccessCount: 1, FailureCount: 0}, nil)

		err := service.SendNotification(context.Background(), "user-123", "Test Title", "Test Body", nil)
		assert.NoError(t, err)
		mockMsgClient.AssertExpectations(t)
	})

	t.Run("No Devices", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
			WithArgs("user-123").
			WillReturnRows(pgxmock.NewRows([]string{"fcm_token"}))

		err := service.SendNotification(context.Background(), "user-123", "Test Title", "Test Body", nil)
		assert.NoError(t, err)
		// Should not call SendEachForMulticast
		mockMsgClient.AssertNotCalled(t, "SendEachForMulticast")
	})

	t.Run("Cleanup Invalid Tokens", func(t *testing.T) {
		// 1. Get tokens
		mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
			WithArgs("user-123").
			WillReturnRows(pgxmock.NewRows([]string{"fcm_token"}).AddRow("valid-token").AddRow("invalid-token"))

		// 2. Send message (returns 1 success, 1 failure)
		mockMsgClient.On("SendEachForMulticast", mock.Anything, mock.Anything).
			Return(&messaging.BatchResponse{
				SuccessCount: 1,
				FailureCount: 1,
				Responses: []*messaging.SendResponse{
					{Success: true},
					{Success: false, Error: errors.New("registration-token-not-registered")},
				},
			}, nil)

		// 3. Remove invalid token
		// Expects ANY($2) for the token list
		// pgxmock with ANY matching is tricky. We can match arguments loosely.
		mockDB.ExpectExec("DELETE FROM user_devices").
			WithArgs("user-123", pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := service.SendNotification(context.Background(), "user-123", "Test Title", "Test Body", nil)
		assert.NoError(t, err)
		mockMsgClient.AssertExpectations(t)
	})
}
