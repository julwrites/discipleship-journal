package database

import (
	"errors"
	"testing"
	"fmt"
	"os"

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

func TestRunMigrations_ConfigError(t *testing.T) {
	os.Unsetenv("DB_USERNAME")
	err := RunMigrations()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to build connection string")
}

func TestRunMigrations_PingError(t *testing.T) {
	os.Setenv("DB_USERNAME", "user")
	os.Setenv("DB_NAME", "db")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "54322") // closed

	err := RunMigrations()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping database")

	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
}

func TestIsMigrationConflictError(t *testing.T) {
	cases := []struct {
		name string
		msg  string
		want bool
	}{
		{"duplicate column", "migration failed: Duplicate column name 'username'", true},
		{"duplicate key name", "migration failed: Duplicate key name 'idx_users_username'", true},
		{"duplicate entry insert", "migration failed: Duplicate entry '02fc167c' for key 'PRIMARY'", true},
		{"already exists", "Table 'notes' already exists", true},
		{"doesn't exist rename source", "Table 'journal_entries' doesn't exist", true},
		{"unknown table drop", "Unknown table 'old_table'", true},
		{"real connection error", "connection refused", false},
		{"syntax error", "You have an error in your SQL syntax", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := fmt.Errorf("%s", tc.msg)
			assert.Equal(t, tc.want, isMigrationConflictError(err))
		})
	}
}
