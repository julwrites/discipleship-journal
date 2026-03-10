package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect_OpenError(t *testing.T) {
	os.Setenv("DB_USERNAME", "user")
	os.Setenv("DB_NAME", "db")
	os.Setenv("DB_PASSWORD", "secretpass")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "invalid-port") // Will fail inside Ping/Open

	err := Connect()
	assert.Error(t, err)
	assert.NotContains(t, err.Error(), "secretpass", "Password should be masked")

	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
}

func TestConnect_CloudSQL(t *testing.T) {
	os.Setenv("DB_USERNAME", "user")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_NAME", "db")
	os.Setenv("CLOUD_SQL_INSTANCE", "project:region:instance")

	_, err := BuildConnectionString()
	assert.NoError(t, err)

	os.Unsetenv("CLOUD_SQL_INSTANCE")
	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
}

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
