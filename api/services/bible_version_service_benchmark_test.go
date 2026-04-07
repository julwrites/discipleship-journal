package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"io"

	"github.com/DATA-DOG/go-sqlmock"
)

type stringReadCloser struct {
	io.Reader
}

func (s *stringReadCloser) Close() error { return nil }

type MockHTTPClientForBenchmark struct {
	body string
}

func (m *MockHTTPClientForBenchmark) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       &stringReadCloser{Reader: strings.NewReader(m.body)},
	}, nil
}

func BenchmarkSyncVersions(b *testing.B) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	// Generate a large number of versions for benchmarking
	const numVersions = 100
	var sb strings.Builder
	sb.WriteString(`<html><body><select name="version" class="search-dropdown">`)
	for i := 0; i < numVersions; i++ {
		// Ensure zero padding for alphabetical order match
		sb.WriteString(fmt.Sprintf(`<option value="V%03d">Version %03d</option>`, i, i))
	}
	sb.WriteString(`</select></body></html>`)
	body := sb.String()

	service := &bibleVersionService{
		db:         db,
		httpClient: &MockHTTPClientForBenchmark{body: body},
		scrapeURL:  "http://localhost",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mockDB.ExpectBegin()

		// For batching, we expect ONE Exec call with ALL arguments
		args := make([]interface{}, 0, numVersions*2)
		for j := 0; j < numVersions; j++ {
			args = append(args, fmt.Sprintf("Version %03d", j), fmt.Sprintf("V%03d", j))
		}

		mockDB.ExpectExec("INSERT INTO bible_versions").
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(1, int64(numVersions)))

		mockDB.ExpectCommit()

		err := service.SyncVersions(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}
