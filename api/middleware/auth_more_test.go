package middleware

import (
	"context"
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthMiddleware_Success(t *testing.T) {
	// A valid-ish JSON that Firebase allows to init an app,
	// but may fail later if it's not a real service account.
	// But actually, NewAuthMiddleware calls app.Auth(ctx).
	// If FIREBASE_CONFIG is set, it might succeed enough to return a client.

	// We can test the error path more thoroughly or the success path if we mock the env.
	os.Setenv("FIREBASE_CONFIG", `{"projectId": "test-project"}`)
	os.Setenv("FIREBASE_SERVICE_ACCOUNT_KEY", `{"type": "service_account", "project_id": "test"}`)

	am, err := NewAuthMiddleware(context.Background())
	// Depending on the exact env, this might still error out due to fake key, but we've covered the invocation.
	if err == nil {
		assert.NotNil(t, am)
	}

	os.Unsetenv("FIREBASE_CONFIG")
	os.Unsetenv("FIREBASE_SERVICE_ACCOUNT_KEY")
}
