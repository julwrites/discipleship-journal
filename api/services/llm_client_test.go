package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLLMClient is a mock implementation of LLMClient
type MockLLMClient struct {
	mock.Mock
	name string
}

func (m *MockLLMClient) Query(ctx context.Context, prompt string, schema string) (string, string, error) {
	args := m.Called(ctx, prompt, schema)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockLLMClient) Stream(ctx context.Context, prompt string) (<-chan string, string, error) {
	args := m.Called(ctx, prompt)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).(<-chan string), args.String(1), args.Error(2)
}

func (m *MockLLMClient) Name() string {
	return m.name
}

func TestFallbackClient_Query(t *testing.T) {
	ctx := context.Background()
	prompt := "hello"
	schema := "{}"

	t.Run("first provider succeeds", func(t *testing.T) {
		m1 := &MockLLMClient{name: "p1"}
		m2 := &MockLLMClient{name: "p2"}

		m1.On("Query", ctx, prompt, schema).Return("response1", "p1", nil)

		fc := NewFallbackClient([]LLMClient{m1, m2}, []string{"p1", "p2"})

		resp, provider, err := fc.Query(ctx, prompt, schema)

		assert.NoError(t, err)
		assert.Equal(t, "response1", resp)
		assert.Equal(t, "p1", provider)
		m1.AssertExpectations(t)
		m2.AssertNotCalled(t, "Query")
	})

	t.Run("first fails, second succeeds", func(t *testing.T) {
		m1 := &MockLLMClient{name: "p1"}
		m2 := &MockLLMClient{name: "p2"}

		m1.On("Query", ctx, prompt, schema).Return("", "", errors.New("fail"))
		m2.On("Query", ctx, prompt, schema).Return("response2", "p2", nil)

		fc := NewFallbackClient([]LLMClient{m1, m2}, []string{"p1", "p2"})

		resp, provider, err := fc.Query(ctx, prompt, schema)

		assert.NoError(t, err)
		assert.Equal(t, "response2", resp)
		assert.Equal(t, "p2", provider)
		m1.AssertExpectations(t)
		m2.AssertExpectations(t)
	})

	t.Run("all fail", func(t *testing.T) {
		m1 := &MockLLMClient{name: "p1"}
		m1.On("Query", ctx, prompt, schema).Return("", "", errors.New("fail1"))

		fc := NewFallbackClient([]LLMClient{m1}, []string{"p1"})

		resp, provider, err := fc.Query(ctx, prompt, schema)

		assert.Error(t, err)
		assert.Equal(t, "", resp)
		assert.Equal(t, "", provider)
		m1.AssertExpectations(t)
	})

	t.Run("preferred provider override via SetPreferredProvider", func(t *testing.T) {
		m1 := &MockLLMClient{name: "p1"}
		m2 := &MockLLMClient{name: "p2"}

		// p2 is preferred, so it should be called first
		m2.On("Query", ctx, prompt, schema).Return("response2", "p2", nil)

		fc := NewFallbackClient([]LLMClient{m1, m2}, []string{"p1", "p2"})
		fc.SetPreferredProvider("p2")

		resp, provider, err := fc.Query(ctx, prompt, schema)

		assert.NoError(t, err)
		assert.Equal(t, "response2", resp)
		assert.Equal(t, "p2", provider)
		m2.AssertExpectations(t)
		m1.AssertNotCalled(t, "Query")
	})

	t.Run("preferred provider override via Context", func(t *testing.T) {
		m1 := &MockLLMClient{name: "p1"}
		m2 := &MockLLMClient{name: "p2"}

		ctxWithPref := context.WithValue(ctx, AIProviderKey, "p2")

		// p2 is preferred via context
		m2.On("Query", ctxWithPref, prompt, schema).Return("response2", "p2", nil)

		fc := NewFallbackClient([]LLMClient{m1, m2}, []string{"p1", "p2"})

		resp, provider, err := fc.Query(ctxWithPref, prompt, schema)

		assert.NoError(t, err)
		assert.Equal(t, "response2", resp)
		assert.Equal(t, "p2", provider)
		m2.AssertExpectations(t)
		m1.AssertNotCalled(t, "Query")
	})

	t.Run("provider not found", func(t *testing.T) {
		m1 := &MockLLMClient{name: "p1"}
		fc := NewFallbackClient([]LLMClient{m1}, []string{"nonexistent"})

		// Should look for nonexistent, skip it, then try p1 if it was in the map?
		// Wait, priorityOrder only dictates order. If a name is in priorityOrder but not in clients map, it is skipped.
		// If p1 is not in priorityOrder, is it ever called?
		// Logic: order = [preferred, ...priorityOrder]
		// For loop iterates order.
		// If p1 is NOT in priorityOrder and NOT preferred, it is NEVER called.
		// Let's verify this behavior.

		resp, provider, err := fc.Query(ctx, prompt, schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "all providers failed")
		assert.Equal(t, "", resp)
		assert.Equal(t, "", provider)
	})
}

func TestFallbackClient_Stream(t *testing.T) {
	ctx := context.Background()
	prompt := "hello"

	t.Run("success", func(t *testing.T) {
		m1 := &MockLLMClient{name: "p1"}
		ch := make(chan string)

		m1.On("Stream", ctx, prompt).Return((<-chan string)(ch), "p1", nil)

		fc := NewFallbackClient([]LLMClient{m1}, []string{"p1"})

		out, prov, err := fc.Stream(ctx, prompt)

		assert.NoError(t, err)
		assert.Equal(t, "p1", prov)
		assert.Equal(t, (<-chan string)(ch), out)
	})

	t.Run("all fail", func(t *testing.T) {
		m1 := &MockLLMClient{name: "p1"}
		m1.On("Stream", ctx, prompt).Return(nil, "", errors.New("fail1"))

		fc := NewFallbackClient([]LLMClient{m1}, []string{"p1"})

		out, prov, err := fc.Stream(ctx, prompt)

		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "", prov)
	})

	t.Run("fallback logic", func(t *testing.T) {
		m1 := &MockLLMClient{name: "p1"}
		m2 := &MockLLMClient{name: "p2"}
		ch := make(chan string)

		m1.On("Stream", ctx, prompt).Return(nil, "", errors.New("fail1"))
		m2.On("Stream", ctx, prompt).Return((<-chan string)(ch), "p2", nil)

		fc := NewFallbackClient([]LLMClient{m1, m2}, []string{"p1", "p2"})

		out, prov, err := fc.Stream(ctx, prompt)

		assert.NoError(t, err)
		assert.Equal(t, "p2", prov)
		assert.Equal(t, (<-chan string)(ch), out)
	})
}

func TestFallbackClient_Name(t *testing.T) {
	fc := NewFallbackClient(nil, nil)
	assert.Equal(t, "fallback", fc.Name())
}
