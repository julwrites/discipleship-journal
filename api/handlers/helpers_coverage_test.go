package handlers

import (
	"context"
	"errors"
	"testing"

	"discipleship_journal_api/middleware"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetUserUUIDFromContext_Coverage(t *testing.T) {
	t.Run("TestUserKey_UUID", func(t *testing.T) {
		id := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, id)
		got, err := GetUserUUIDFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, id, got)
	})

	t.Run("TestUserKey_String", func(t *testing.T) {
		id := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, id.String())
		got, err := GetUserUUIDFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, id, got)
	})

	t.Run("TestUserKey_InvalidType", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TestUserKey, 12345)
		_, err := GetUserUUIDFromContext(ctx)
		// Should fall through to other checks, then fail
		assert.Error(t, err)
	})

	t.Run("Production_Token_With_DB", func(t *testing.T) {
		// Mock DB via provider
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mock.Close()

		oldProvider := dbProvider
		dbProvider = func() DBInterface { return mock }
		defer func() { dbProvider = oldProvider }()

		uid := "firebase-uid"
		dbUUID := uuid.New()

		mock.ExpectQuery("SELECT id FROM users").
			WithArgs(uid).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(dbUUID))

		ctx := context.WithValue(context.Background(), middleware.UserContextKey, &auth.Token{UID: uid})
		got, err := GetUserUUIDFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, dbUUID, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Missing_Context", func(t *testing.T) {
		_, err := GetUserUUIDFromContext(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "user not found in context", err.Error())
	})
}

func TestGetUserUUID_Coverage(t *testing.T) {
	t.Run("Empty_UID", func(t *testing.T) {
		_, err := GetUserUUID(context.Background(), "")
		assert.Error(t, err)
		assert.Equal(t, "invalid firebase UID", err.Error())
	})

	t.Run("DB_Not_Initialized", func(t *testing.T) {
		oldProvider := dbProvider
		dbProvider = func() DBInterface { return nil }
		defer func() { dbProvider = oldProvider }()

		_, err := GetUserUUID(context.Background(), "uid")
		assert.Error(t, err)
		assert.Equal(t, "database not initialized", err.Error())
	})

	t.Run("DB_Error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mock.Close()

		oldProvider := dbProvider
		dbProvider = func() DBInterface { return mock }
		defer func() { dbProvider = oldProvider }()

		mock.ExpectQuery("SELECT id FROM users").
			WithArgs("uid").
			WillReturnError(errors.New("db error"))

		_, err = GetUserUUID(context.Background(), "uid")
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("TestUserKey_Override_Direct", func(t *testing.T) {
		id := uuid.New()
		ctx := context.WithValue(context.Background(), TestUserKey, id)
		got, err := GetUserUUID(ctx, "any-uid")
		assert.NoError(t, err)
		assert.Equal(t, id, got)
	})
}
