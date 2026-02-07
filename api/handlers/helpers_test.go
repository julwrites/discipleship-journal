package handlers

import (
	"context"
	"net/http/httptest"
	"testing"

	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserUUIDFromContext(t *testing.T) {
	t.Run("TestUserKey_UUID", func(t *testing.T) {
		id := uuid.New()
		req := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, id)

		resID, err := GetUserUUIDFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, id, resID)
	})

	t.Run("TestUserKey_String", func(t *testing.T) {
		id := uuid.New()
		req := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, id.String())

		resID, err := GetUserUUIDFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, id, resID)
	})

	t.Run("UserContextKey_WithTestUserKeyOverride", func(t *testing.T) {
		// Should prefer TestUserKey if present
		testID := uuid.New()
		req := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, testID)
		ctx = context.WithValue(ctx, middleware.UserContextKey, &auth.Token{UID: "uid"})

		resID, err := GetUserUUIDFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, testID, resID)
	})

	t.Run("NoUserContext", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		// No context keys

		_, err := GetUserUUIDFromContext(req.Context())
		assert.Error(t, err)
		assert.Equal(t, "user not found in context", err.Error())
	})

	t.Run("TestUserKey_InvalidType", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), TestUserKey, 12345) // int, not UUID or string

		_, err := GetUserUUIDFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "user not found in context", err.Error())
	})

	t.Run("UserContextKey_InvalidType", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, "invalid-token-type")

		_, err := GetUserUUIDFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "user not found in context", err.Error())
	})

	t.Run("UserContextKey_NilToken", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		var token *auth.Token = nil
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, token)

		_, err := GetUserUUIDFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "user not found in context", err.Error())
	})
}

func TestGetUserUUID(t *testing.T) {
	t.Run("WithTestUserKey", func(t *testing.T) {
		id := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, id)

		resID, err := GetUserUUID(ctx, "any-uid")
		assert.NoError(t, err)
		assert.Equal(t, id, resID)
	})

	t.Run("WithTestUserKeyString", func(t *testing.T) {
		id := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, id.String())

		resID, err := GetUserUUID(ctx, "any-uid")
		assert.NoError(t, err)
		assert.Equal(t, id, resID)
	})

	t.Run("EmptyUID", func(t *testing.T) {
		// If TestUserKey is nil
		_, err := GetUserUUID(context.Background(), "")
		assert.Error(t, err)
		assert.Equal(t, "invalid firebase UID", err.Error())
	})

    t.Run("DBNotInitialized", func(t *testing.T) {
        // If we call with UID but no TestUserKey, it hits DB.
        // If DB is nil, it returns error or panics?
        // Code: if database.DB == nil { return uuid.Nil, errors.New("database not initialized") }

        // This test assumes database.DB is nil (which it is in isolated tests usually).
        // If other tests initialize it, this might fail or crash if DB is closed.
        // Assuming isolated run or DB is nil.
        _, err := GetUserUUID(context.Background(), "some-uid")
        // We expect "database not initialized" OR db error if initialized.
        // It's safer to just check error.
        assert.Error(t, err)
    })
}
