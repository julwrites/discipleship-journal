package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBibleAIClient is a mock implementation of services.BibleAIClient
type MockBibleAIClient struct {
	mock.Mock
}

func (m *MockBibleAIClient) Query(ctx context.Context, prompt string, schema string) (string, string, error) {
	args := m.Called(ctx, prompt, schema)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockBibleAIClient) Stream(ctx context.Context, prompt string) (<-chan string, string, error) {
	args := m.Called(ctx, prompt)
	return args.Get(0).(<-chan string), args.String(1), args.Error(2)
}

func (m *MockBibleAIClient) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockBibleAIClient) GetPassage(ctx context.Context, reference string, version string) (map[string]interface{}, error) {
	args := m.Called(ctx, reference, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockBibleAIClient) ChatCompletion(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockBibleAIClient) StreamChatCompletion(ctx context.Context, payload map[string]interface{}) (<-chan string, <-chan error, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(<-chan string), args.Get(1).(<-chan error), args.Error(2)
}

func (m *MockBibleAIClient) GetVersions(ctx context.Context, params map[string]string) (map[string]interface{}, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockBibleAIClient) GetSystemPrompt(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func TestBibleHandler_GetBiblePassage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		handler := NewBibleHandler(mockClient)

		mockClient.On("GetPassage", mock.Anything, "John 3:16", "ESV").Return(map[string]interface{}{
			"text": "For God so loved the world...",
		}, nil)

		req, _ := http.NewRequest("GET", "/api/bible/passage?ref=John%203:16&version=ESV", nil)
		rr := httptest.NewRecorder()

		handler.GetBiblePassage(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, "For God so loved the world...", resp["text"])
	})

	t.Run("MissingRef", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		handler := NewBibleHandler(mockClient)

		req, _ := http.NewRequest("GET", "/api/bible/passage", nil)
		rr := httptest.NewRecorder()

		handler.GetBiblePassage(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("ClientError", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		handler := NewBibleHandler(mockClient)

		mockClient.On("GetPassage", mock.Anything, "John 3:16", "").Return(nil, errors.New("api error"))

		req, _ := http.NewRequest("GET", "/api/bible/passage?ref=John%203:16", nil)
		rr := httptest.NewRecorder()

		handler.GetBiblePassage(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestBibleHandler_GetBibleVersions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		handler := NewBibleHandler(mockClient)

		mockClient.On("GetVersions", mock.Anything, mock.MatchedBy(func(params map[string]string) bool {
			return params["language"] == "en"
		})).Return(map[string]interface{}{
			"data": []interface{}{"ESV", "KJV"},
		}, nil)

		req, _ := http.NewRequest("GET", "/api/bible/versions?language=en", nil)
		rr := httptest.NewRecorder()

		handler.GetBibleVersions(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("ClientError", func(t *testing.T) {
		mockClient := new(MockBibleAIClient)
		handler := NewBibleHandler(mockClient)

		mockClient.On("GetVersions", mock.Anything, mock.Anything).Return(nil, errors.New("failed"))

		req, _ := http.NewRequest("GET", "/api/bible/versions", nil)
		rr := httptest.NewRecorder()

		handler.GetBibleVersions(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
