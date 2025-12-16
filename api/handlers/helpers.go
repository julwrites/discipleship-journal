package handlers

import (
	"context"
	"errors"

	"discipleship_journal_api/database"
	"github.com/google/uuid"
)

// GetUserUUID is a helper to get the UUID of the user from the database given the Firebase UID.
func GetUserUUID(ctx context.Context, firebaseUID string) (uuid.UUID, error) {
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
