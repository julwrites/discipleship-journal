package database

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestApplyMigrations(t *testing.T) {
	err := errors.New("Duplicate column name")
	assert.True(t, isMigrationConflictError(err))

	err = errors.New("Duplicate key name")
	assert.True(t, isMigrationConflictError(err))

	err = errors.New("Duplicate entry")
	assert.True(t, isMigrationConflictError(err))

	err = errors.New("already exists")
	assert.True(t, isMigrationConflictError(err))

	err = errors.New("doesn't exist")
	assert.True(t, isMigrationConflictError(err))

	err = errors.New("Unknown table")
	assert.True(t, isMigrationConflictError(err))

	err = errors.New("Other error")
	assert.False(t, isMigrationConflictError(err))
}
