package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Register pgx driver
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetupIntegrationDB spins up a Postgres container, applies migrations, and returns a connection pool.
// It skips the test if running in -short mode.
// The caller is responsible for closing the pool and terminating the container (via the returned cleanup func).
func SetupIntegrationDB(t *testing.T) (*pgxpool.Pool, func()) {
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

	// Use standard Docker Hub image. ECR public mirror was flaky in CI.
	postgresContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// Run Migrations
	// We need to find the absolute path to the migrations directory
	// Assuming this test is running from within api/tests/integration or similar
	// We might need to walk up to find 'api/migrations'
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %s", err)
	}

	// Helper to find migrations folder.
	// If running from api/tests/integration, it is ../../../migrations
	// If running from api/tests/contract, it is ../../../migrations
	// We'll search up to 5 levels up.
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
		// Fallback: assume we are in api root context if the above failed (unlikely)
		migrationsPath = "migrations"
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		connStr,
	)
	if err != nil {
		// Try absolute path if relative failed
		absPath, _ := filepath.Abs("../../../migrations")
		t.Logf("Migration path lookup failed, trying hardcoded relative path: %s", absPath)
		m, err = migrate.New(
			"file://"+absPath,
			connStr,
		)
		if err != nil {
			t.Fatalf("failed to create migrate instance (path: %s): %s", migrationsPath, err)
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to run migrations: %s", err)
	}

	// Create pgxpool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create pgxpool: %s", err)
	}

	cleanup := func() {
		pool.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return pool, cleanup
}
