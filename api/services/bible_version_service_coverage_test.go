package services

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBibleVersionService_GetVersions_Coverage(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := NewBibleVersionService(db)
	ctx := context.Background()

	t.Run("Query_Error", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT .* FROM bible_versions").
			WillReturnError(errors.New("query error"))

		_, err := service.GetVersions(ctx)
		assert.Error(t, err)
		assert.Equal(t, "failed to query versions: query error", err.Error())
	})

	t.Run("Scan_Error", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT .* FROM bible_versions").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "abbreviation", "created_at", "updated_at"}).
				AddRow("id", "Name", "Abbr", "invalid-date", time.Now()))

		_, err := service.GetVersions(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to scan version")
	})
}

// MockHTTPClient for ScrapeVersions
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestBibleVersionService_ScrapeVersions_Coverage(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Use struct directly to inject mock http client
	mockHTTP := new(MockHTTPClient)
	service := &bibleVersionService{
		db:         db,
		httpClient: mockHTTP,
		scrapeURL:  "http://test.com",
	}
	ctx := context.Background()

	t.Run("HTTP_Do_Error", func(t *testing.T) {
		mockHTTP.On("Do", mock.Anything).Return(nil, errors.New("http error")).Once()

		_, err := service.ScrapeVersions(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch versions page")
	})

	t.Run("Status_Not_OK", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       http.NoBody,
		}
		mockHTTP.On("Do", mock.Anything).Return(resp, nil).Once()

		_, err := service.ScrapeVersions(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
	})
}
