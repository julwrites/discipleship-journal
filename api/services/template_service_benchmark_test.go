package services

import (
	"context"
	"testing"
	"time"

	"discipleship_journal_api/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// SlowMockBibleAIClient adds delay to GetPassage
type SlowMockBibleAIClient struct {
	*MockBibleAIClient
	Delay time.Duration
}

func (m *SlowMockBibleAIClient) GetPassage(ctx context.Context, reference string, version string) (map[string]interface{}, error) {
	time.Sleep(m.Delay)
	return m.MockBibleAIClient.GetPassage(ctx, reference, version)
}

// SimpleBenchmarkMockDB implements DBInterfaceWithQuery for benchmarking
type SimpleBenchmarkMockDB struct {}

func (m *SimpleBenchmarkMockDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *SimpleBenchmarkMockDB) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (m *SimpleBenchmarkMockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}
func (m *SimpleBenchmarkMockDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return &SimpleBenchmarkRow{}
}

type SimpleBenchmarkRow struct{}

func (r *SimpleBenchmarkRow) Scan(dest ...any) error {
	// Index 9 is AllowUserPassages (*bool)
	if len(dest) > 9 {
		if ptr, ok := dest[9].(*bool); ok {
			*ptr = true
		}
	}
	// Index 5 is Prompts (*[]byte) - required to avoid unmarshal error
	if len(dest) > 5 {
		if ptr, ok := dest[5].(*[]byte); ok {
			*ptr = []byte("{}")
		}
	}
	return nil
}

func BenchmarkGenerateContent_MultipleReferences(b *testing.B) {
	// Simulate a slow network call (e.g., 20ms)
	mockAI := &SlowMockBibleAIClient{
		MockBibleAIClient: NewMockBibleAIClient(),
		Delay:             20 * time.Millisecond,
	}
	mockDB := &SimpleBenchmarkMockDB{}
	service := NewTemplateService(mockDB, mockAI)

	tmplID := uuid.New()
	req := models.GenerateRequest{
		Inputs: map[string]string{"key": "value"},
		UserPassages: []string{"Gen 1:1", "Ex 20:1", "Lev 19:18", "Num 6:24", "Deut 6:5"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GenerateContent(context.Background(), tmplID, req)
	}
}
