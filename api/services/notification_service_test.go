package services

import (
	"context"
	"errors"
	"testing"

	"firebase.google.com/go/v4/messaging"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	t.Run("Success", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening stub database: %s", err)
		}
		defer db.Close()

		mockMsgClient := new(MockMessagingClient)
		service := NewNotificationService(db, mockMsgClient)

		mockDB.ExpectExec("INSERT INTO user_devices").
			WithArgs("user-123", "token-123", "ios").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = service.RegisterDevice(context.Background(), "user-123", "token-123", "ios")
		assert.NoError(t, err)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("DB Error", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening stub database: %s", err)
		}
		defer db.Close()

		mockMsgClient := new(MockMessagingClient)
		service := NewNotificationService(db, mockMsgClient)

		mockDB.ExpectExec("INSERT INTO user_devices").
			WithArgs("user-123", "token-123", "ios").
			WillReturnError(errors.New("db error"))

		err = service.RegisterDevice(context.Background(), "user-123", "token-123", "ios")
		assert.Error(t, err)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestNotificationService_SendNotification(t *testing.T) {
	t.Run("Success - Sent to devices", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening stub database: %s", err)
		}
		defer db.Close()

		mockMsgClient := new(MockMessagingClient)
		service := NewNotificationService(db, mockMsgClient)

		// 1. Get tokens
		mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
			WithArgs("user-123").
			WillReturnRows(sqlmock.NewRows([]string{"fcm_token"}).AddRow("token-1"))

		// 2. Send message
		mockMsgClient.On("SendEachForMulticast", mock.Anything, mock.MatchedBy(func(msg *messaging.MulticastMessage) bool {
			return len(msg.Tokens) == 1 && msg.Tokens[0] == "token-1" && msg.Notification.Title == "Test Title"
		})).Return(&messaging.BatchResponse{SuccessCount: 1, FailureCount: 0}, nil)

		err = service.SendNotification(context.Background(), "user-123", "Test Title", "Test Body", nil)
		assert.NoError(t, err)
		mockMsgClient.AssertExpectations(t)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("No Devices", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening stub database: %s", err)
		}
		defer db.Close()

		mockMsgClient := new(MockMessagingClient)
		service := NewNotificationService(db, mockMsgClient)

		mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
			WithArgs("user-123").
			WillReturnRows(sqlmock.NewRows([]string{"fcm_token"}))

		err = service.SendNotification(context.Background(), "user-123", "Test Title", "Test Body", nil)
		assert.NoError(t, err)
		// Should not call SendEachForMulticast
		mockMsgClient.AssertNotCalled(t, "SendEachForMulticast")
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Cleanup Invalid Tokens", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening stub database: %s", err)
		}
		defer db.Close()

		mockMsgClient := new(MockMessagingClient)
		service := NewNotificationService(db, mockMsgClient)

		// 1. Get tokens
		mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
			WithArgs("user-123").
			WillReturnRows(sqlmock.NewRows([]string{"fcm_token"}).AddRow("valid-token").AddRow("invalid-token"))

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
		mockDB.ExpectExec("DELETE FROM user_devices").
			WithArgs("user-123", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = service.SendNotification(context.Background(), "user-123", "Test Title", "Test Body", nil)
		assert.NoError(t, err)
		mockMsgClient.AssertExpectations(t)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestNotificationService_SendNotification_GetTokensError(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening stub database: %s", err)
	}
	defer db.Close()

	mockMsgClient := new(MockMessagingClient)
	service := NewNotificationService(db, mockMsgClient)

	mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
		WithArgs("user-123").
		WillReturnError(errors.New("db error"))

	err = service.SendNotification(context.Background(), "user-123", "Title", "Body", nil)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestNotificationService_SendNotification_SendError(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening stub database: %s", err)
	}
	defer db.Close()

	mockMsgClient := new(MockMessagingClient)
	service := NewNotificationService(db, mockMsgClient)

	mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
		WithArgs("user-123").
		WillReturnRows(sqlmock.NewRows([]string{"fcm_token"}).AddRow("token"))

	mockMsgClient.On("SendEachForMulticast", mock.Anything, mock.Anything).
		Return(nil, errors.New("fcm error"))

	err = service.SendNotification(context.Background(), "user-123", "Title", "Body", nil)
	assert.Error(t, err)
	assert.Equal(t, "fcm error", err.Error())
}

func TestNotificationService_SendNotification_CleanupError(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening stub database: %s", err)
	}
	defer db.Close()

	mockMsgClient := new(MockMessagingClient)
	service := NewNotificationService(db, mockMsgClient)

	// 1. Get tokens
	mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
		WithArgs("user-123").
		WillReturnRows(sqlmock.NewRows([]string{"fcm_token"}).AddRow("invalid-token"))

	// 2. Send message failure
	mockMsgClient.On("SendEachForMulticast", mock.Anything, mock.Anything).
		Return(&messaging.BatchResponse{
			SuccessCount: 0,
			FailureCount: 1,
			Responses: []*messaging.SendResponse{
				{Success: false, Error: errors.New("error")},
			},
		}, nil)

	// 3. Remove invalid token -> DB Error
	mockDB.ExpectExec("DELETE FROM user_devices").
		WithArgs("user-123", sqlmock.AnyArg()).
		WillReturnError(errors.New("db delete error"))

	err = service.SendNotification(context.Background(), "user-123", "Title", "Body", nil)
	// Should not return error, just log it
	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestNotificationService_SendMulticastNotification(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening stub database: %s", err)
		}
		defer db.Close()

		mockMsgClient := new(MockMessagingClient)
		service := NewNotificationService(db, mockMsgClient)

		userIDs := []string{"user-1", "user-2"}

		// 1. Get tokens
		mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
			WithArgs(userIDs[0], userIDs[1]).
			WillReturnRows(sqlmock.NewRows([]string{"fcm_token"}).AddRow("token-1").AddRow("token-2"))

		// 2. Send message
		mockMsgClient.On("SendEachForMulticast", mock.Anything, mock.MatchedBy(func(msg *messaging.MulticastMessage) bool {
			return len(msg.Tokens) == 2 && msg.Tokens[0] == "token-1" && msg.Tokens[1] == "token-2"
		})).Return(&messaging.BatchResponse{SuccessCount: 2, FailureCount: 0}, nil)

		err = service.SendMulticastNotification(context.Background(), userIDs, "Title", "Body", nil)
		assert.NoError(t, err)
		mockMsgClient.AssertExpectations(t)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("No Users", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		defer db.Close()
		service := NewNotificationService(db, nil)
		err := service.SendMulticastNotification(context.Background(), []string{}, "Title", "Body", nil)
		assert.NoError(t, err)
	})

	t.Run("Cleanup Failed Tokens", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error opening stub database: %s", err)
		}
		defer db.Close()

		mockMsgClient := new(MockMessagingClient)
		service := NewNotificationService(db, mockMsgClient)

		userIDs := []string{"user-1", "user-2"}

		// 1. Get tokens
		mockDB.ExpectQuery("SELECT fcm_token FROM user_devices").
			WithArgs(userIDs[0], userIDs[1]).
			WillReturnRows(sqlmock.NewRows([]string{"fcm_token"}).AddRow("token-fail"))

		// 2. Send message
		mockMsgClient.On("SendEachForMulticast", mock.Anything, mock.Anything).
			Return(&messaging.BatchResponse{
				SuccessCount: 0,
				FailureCount: 1,
				Responses: []*messaging.SendResponse{
					{Success: false, Error: errors.New("invalid")},
				},
			}, nil)

		// 3. Remove tokens
		mockDB.ExpectExec("DELETE FROM user_devices").
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = service.SendMulticastNotification(context.Background(), userIDs, "Title", "Body", nil)
		assert.NoError(t, err)
		mockMsgClient.AssertExpectations(t)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}
