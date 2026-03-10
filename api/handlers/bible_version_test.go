package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBibleVersionService is a mock implementation of services.BibleVersionService
type MockBibleVersionService struct {
	mock.Mock
}

func (m *MockBibleVersionService) GetVersions(ctx context.Context) ([]services.BibleVersion, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]services.BibleVersion), args.Error(1)
}

func (m *MockBibleVersionService) SyncVersions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBibleVersionService) ScrapeVersions(ctx context.Context) (map[string]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func TestBibleVersionHandler_GetVersions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockBibleVersionService)
		handler := NewBibleVersionHandler(mockService)

		expectedVersions := []services.BibleVersion{
			{ID: "1", Name: "Version 1", Abbreviation: "V1"},
			{ID: "2", Name: "Version 2", Abbreviation: "V2"},
		}

		mockService.On("GetVersions", mock.Anything).Return(expectedVersions, nil)

		req, _ := http.NewRequest("GET", "/bible-versions", nil)
		rr := httptest.NewRecorder()

		handler.GetVersions(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var versions []services.BibleVersion
		err := json.Unmarshal(rr.Body.Bytes(), &versions)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedVersions), len(versions))
		assert.Equal(t, expectedVersions[0].Name, versions[0].Name)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockBibleVersionService)
		handler := NewBibleVersionHandler(mockService)

		mockService.On("GetVersions", mock.Anything).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest("GET", "/bible-versions", nil)
		rr := httptest.NewRecorder()

		handler.GetVersions(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
