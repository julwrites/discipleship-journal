package database

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/stub"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"discipleship_journal_api/migrations"
)

func TestApplyMigrations_CoverageMore(t *testing.T) {
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

func TestApplyMigrations_Stub(t *testing.T) {
	sourceDriver, err := iofs.New(migrations.FS, ".")
	assert.NoError(t, err)

	d := &stub.Stub{}
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "stub", d)
	if err == nil && m != nil {
		_ = applyMigrations(m)
	}
}
