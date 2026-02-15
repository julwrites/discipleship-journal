package services

import (
	"context"
	"testing"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"firebase.google.com/go/v4/messaging"
	"github.com/stretchr/testify/mock"
)

// SimpleMockDB implements DBInterface for benchmarks
type SimpleMockDB struct{}

func (m *SimpleMockDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return &SimpleMockRows{tokens: []string{"token-1"}}, nil
}
func (m *SimpleMockDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row { return nil }
func (m *SimpleMockDB) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}
func (m *SimpleMockDB) Begin(ctx context.Context) (pgx.Tx, error) { return nil, nil }

type SimpleMockRows struct {
	pgx.Rows
	tokens []string
	index  int
}

func (m *SimpleMockRows) Next() bool {
	return m.index < len(m.tokens)
}
func (m *SimpleMockRows) Scan(dest ...any) error {
	*dest[0].(*string) = m.tokens[m.index]
	m.index++
	return nil
}
func (m *SimpleMockRows) Close()     {}
func (m *SimpleMockRows) Err() error { return nil }

func BenchmarkSendNotificationN(b *testing.B) {
	mockDB := &SimpleMockDB{}
	mockMsgClient := new(MockMessagingClient)
	service := NewNotificationService(mockDB, mockMsgClient)

	ctx := context.Background()
	userIDs := make([]string, 10)
	for i := range userIDs {
		userIDs[i] = "user-1"
	}

	mockMsgClient.On("SendEachForMulticast", mock.Anything, mock.Anything).
		Return(&messaging.BatchResponse{SuccessCount: 1, FailureCount: 0}, nil).
		Maybe()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, userID := range userIDs {
			_ = service.SendNotification(ctx, userID, "Title", "Body", nil)
		}
	}
}

func BenchmarkSendMulticastNotification(b *testing.B) {
	mockDB := &SimpleMockDB{}
	mockMsgClient := new(MockMessagingClient)
	service := NewNotificationService(mockDB, mockMsgClient)

	ctx := context.Background()
	userIDs := make([]string, 10)
	for i := range userIDs {
		userIDs[i] = "user-1"
	}

	mockMsgClient.On("SendEachForMulticast", mock.Anything, mock.Anything).
		Return(&messaging.BatchResponse{SuccessCount: 1, FailureCount: 0}, nil).
		Maybe()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.SendMulticastNotification(ctx, userIDs, "Title", "Body", nil)
	}
}
