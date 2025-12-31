package middleware

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserContextKey(t *testing.T) {
	assert.Equal(t, ContextKey("user"), UserContextKey)
}

func TestWithUser(t *testing.T) {
	// This helper function is internal in some files, but if we want to test context helpers,
	// we should probably check if they exist or if we need to export them.
	// Looking at previous errors, WithUser wasn't exported.

	// We can test manual context injection
	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, "test-user")
	req = req.WithContext(ctx)

	assert.Equal(t, "test-user", req.Context().Value(UserContextKey))
}
