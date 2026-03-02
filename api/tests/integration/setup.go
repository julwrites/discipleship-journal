package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// SetupIntegrationDB spins up a MySQL container, applies migrations, and returns a connection pool.
// It skips the test if running in -short mode.
// The caller is responsible for closing the pool and terminating the container (via the returned cleanup func).
func SetupIntegrationDB(t *testing.T) (*sql.DB, func()) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Pre-check for Docker availability to avoid panics from testcontainers-go
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		t.Skipf("skipping integration test: docker not available or permission denied: %v", err)
	}

	ctx := context.Background()

	dbName := "journal_test"
	dbUser := "testuser"
	dbPassword := "testpassword"

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8",
		mysql.WithDatabase(dbName),
		mysql.WithUsername(dbUser),
		mysql.WithPassword(dbPassword),
	)
	if err != nil {
		t.Fatalf("failed to start mysql container: %s", err)
	}

	connStr, err := mysqlContainer.ConnectionString(ctx, "multiStatements=true&parseTime=true")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// Run Migrations
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %s", err)
	}

	migrationsPath := ""
	checkPath := wd
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(checkPath, "migrations")
		if _, err := os.Stat(candidate); err == nil {
			migrationsPath = candidate
			break
		}
		checkPath = filepath.Dir(checkPath)
	}

	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	migrateConnStr := "mysql://" + connStr
	m, err := migrate.New(
		"file://"+migrationsPath,
		migrateConnStr,
	)
	if err != nil {
		absPath, _ := filepath.Abs("../../../migrations")
		t.Logf("Migration path lookup failed, trying hardcoded relative path: %s", absPath)
		m, err = migrate.New(
			"file://"+absPath,
			migrateConnStr,
		)
		if err != nil {
			t.Fatalf("failed to create migrate instance (path: %s): %s", migrationsPath, err)
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to run migrations: %s", err)
	}

	pool, err := sql.Open("mysql", connStr)
	if err != nil {
		t.Fatalf("failed to create pool: %s", err)
	}

	cleanup := func() {
		pool.Close()
		if err := mysqlContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return pool, cleanup
}
