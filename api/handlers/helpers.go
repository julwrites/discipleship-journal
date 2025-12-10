package handlers

import (
	"context"
	"errors"
	"log/slog"

	"discipleship_journal_api/database"
	"github.com/google/uuid"
)

// GetUserUUID is a helper to get the UUID of the user from the database given the Firebase UID.
func GetUserUUID(ctx context.Context, firebaseUID string) (uuid.UUID, error) {
	if firebaseUID == "" {
		return uuid.Nil, errors.New("invalid firebase UID")
	}

	var id string
	err := database.DB.QueryRow(ctx, "SELECT id FROM users WHERE firebase_uid = $1", firebaseUID).Scan(&id)
	if err != nil {
		slog.Error("Failed to get user UUID", "firebaseUID", firebaseUID, "error", err)
		return uuid.Nil, err
	}

	return uuid.Parse(id)
}
