package handlers

import (
	"context"
	"errors"

	"discipleship_journal_api/database"
	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
)

type contextKey string

const TestUserKey contextKey = "test_user_id"

// GetUserUUIDFromContext retrieves the user ID from the context, supporting both
// test overrides (TestUserKey) and production auth tokens (middleware.UserContextKey).
func GetUserUUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	// 1. Check for test override first
	if val := ctx.Value(TestUserKey); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			return id, nil
		}
		if idStr, ok := val.(string); ok {
			return uuid.Parse(idStr)
		}
	}

	// 2. Check for production auth token
	if val := ctx.Value(middleware.UserContextKey); val != nil {
		if token, ok := val.(*auth.Token); ok && token != nil {
			return GetUserUUID(ctx, token.UID)
		}
	}

	return uuid.Nil, errors.New("user not found in context")
}

// GetUserUUID is a helper to get the UUID of the user from the database given the Firebase UID.
func GetUserUUID(ctx context.Context, firebaseUID string) (uuid.UUID, error) {
	// Re-check test key just in case called directly, though GetUserUUIDFromContext handles it.
	if val := ctx.Value(TestUserKey); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			return id, nil
		}
		if idStr, ok := val.(string); ok {
			return uuid.Parse(idStr)
		}
	}

	if firebaseUID == "" {
		return uuid.Nil, errors.New("invalid firebase UID")
	}

	// Use database.DB global pool
	// Note: In tests where DB is not initialized, this will panic or fail if called.
	// Handlers using this should ideally be refactored to use dependency injection.
	if database.DB == nil {
		return uuid.Nil, errors.New("database not initialized")
	}

	var id uuid.UUID
	err := database.DB.QueryRow(ctx, "SELECT id FROM users WHERE firebase_uid=$1", firebaseUID).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
