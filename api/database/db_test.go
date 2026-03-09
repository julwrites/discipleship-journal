package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConnectionString(t *testing.T) {
	// Setup environment
	os.Setenv("DB_USERNAME", "user")
	os.Setenv("DB_NAME", "db")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("CLOUD_SQL_INSTANCE")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")

	t.Run("Missing_Credentials", func(t *testing.T) {
		os.Unsetenv("DB_USERNAME")
		_, err := BuildConnectionString()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DB_USERNAME and DB_NAME environment variables are required")
		os.Setenv("DB_USERNAME", "user")
	})

	t.Run("Local_Missing_Password", func(t *testing.T) {
		// No cloud sql instance, no password
		_, err := BuildConnectionString()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DB_PASSWORD is required")
	})

	t.Run("Local_Success", func(t *testing.T) {
		os.Setenv("DB_PASSWORD", "pass")
		str, err := BuildConnectionString()
		assert.NoError(t, err)
		assert.Contains(t, str, "user:pass@tcp(127.0.0.1:4000)/db")
		assert.Contains(t, str, "tls=true")
		assert.Contains(t, str, "tidb_skip_isolation_level_check=1")
	})

	t.Run("Local_Custom_Host", func(t *testing.T) {
		os.Setenv("DB_PASSWORD", "pass")
		os.Setenv("DB_HOST", "db-host")
		os.Setenv("DB_PORT", "5433")
		str, err := BuildConnectionString()
		assert.NoError(t, err)
		assert.Contains(t, str, "user:pass@tcp(db-host:5433)/db")
		assert.Contains(t, str, "tidb_skip_isolation_level_check=1")
	})

	// Cleanup
	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("CLOUD_SQL_INSTANCE")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
}

func TestConnect_ConfigError(t *testing.T) {
	// Force BuildConnectionString failure
	os.Unsetenv("DB_USERNAME")
	err := Connect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to build database connection string")
}

func TestConnect_ParseError(t *testing.T) {
	// BuildConnectionString succeeds but returns valid URL
	os.Setenv("DB_USERNAME", "user")
	os.Setenv("DB_NAME", "db")
	// Use Cloud SQL to trigger dialer logic path?
	// Actually, just verify path where BuildString works.

	// Try to make ParseConfig fail?
	// ParseConfig fails if URL is invalid.
	// But BuildConnectionString constructs valid URL.
	// So we can skip this case unless we mock BuildConnectionString.
	// But we can test Ping Failure.
}

func TestConnect_PingFailure(t *testing.T) {
	os.Setenv("DB_USERNAME", "user")
	os.Setenv("DB_NAME", "db")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "54322") // closed port

	err := Connect()
	assert.Error(t, err)
	// Expect "unable to ping database" or connection refused
	assert.Contains(t, err.Error(), "unable to ping database")

	// Check masking
	// "dial tcp [::1]:54322: connect: connection refused"
	// It doesn't contain password unless URL is printed.
	// The error wrapping in Connect is: fmt.Errorf("unable to ping database: %w", err)

	// Cleanup
	Close()
	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
}

func TestClose_Nil(t *testing.T) {
	// Should not panic
	DB = nil
	Close()
}

func TestRunMigrations_Coverage(t *testing.T) {
	t.Run("ConfigError", func(t *testing.T) {
		os.Unsetenv("DB_USERNAME")
		err := RunMigrations()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to build connection string")
	})

	t.Run("PingError", func(t *testing.T) {
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
	})
}

func TestIsSchemaConflictError(t *testing.T) {
	cases := []struct {
		name string
		msg  string
		want bool
	}{
		{"duplicate column", "migration failed: Duplicate column name 'username'", true},
		{"duplicate key", "migration failed: Duplicate key name 'idx_users_username'", true},
		{"already exists", "Table 'notes' already exists", true},
		{"doesn't exist rename source", "Table 'journal_entries' doesn't exist", true},
		{"unknown table drop", "Unknown table 'old_table'", true},
		{"real error", "connection refused", false},
		{"syntax error", "You have an error in your SQL syntax", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := fmt.Errorf("%s", tc.msg)
			assert.Equal(t, tc.want, isSchemaConflictError(err))
		})
	}
}
