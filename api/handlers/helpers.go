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

	// We scan into uuid.UUID directly to avoid "destination kind 'string' not supported for value kind 'array' of column 'id'"
	// when using pgx with pgxmock or real DB if id is UUID type.
	// Actually, if the column is UUID, pgx might want uuid.UUID or [16]byte.
	// Since I'm using google/uuid, I should scan into that if pgx supports it (with pgx-uuid adapter)
	// or scan into string if pgx maps UUID to string.
	// The error "destination kind 'string' not supported for value kind 'array'" usually implies the driver sees it as an array of bytes (UUID)
	// but we try to scan into string.
	// Let's try scanning into uuid.UUID directly.

	var id uuid.UUID
	err := database.DB.QueryRow(ctx, "SELECT id FROM users WHERE firebase_uid = $1", firebaseUID).Scan(&id)
	if err != nil {
		slog.Error("Failed to get user UUID", "firebaseUID", firebaseUID, "error", err)
		return uuid.Nil, err
	}

	return id, nil
}
