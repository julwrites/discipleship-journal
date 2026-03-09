package database

import (
	"errors"
	"fmt"
	"log/slog"

	"database/sql"
	"discipleship_journal_api/migrations"

	"github.com/golang-migrate/migrate/v4"
	mysql_migrate "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// RunMigrations applies all pending database migrations
func RunMigrations() error {
	dsn, err := BuildConnectionString()
	if err != nil {
		return fmt.Errorf("failed to build connection string for migrations: %w", err)
	}

	// multiStatements=true is required because each migration file contains
	// multiple SQL statements. This is intentionally NOT added to the main app
	// connection string (db.go) since it is a security risk for general queries.
	migrationDSN := dsn + "&multiStatements=true"

	db, err := sql.Open("mysql", migrationDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to mysql for migrations: %w", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database before migrations: %w", err)
	}

	slog.Info("Starting database migrations to TiDB...")

	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs source driver: %w", err)
	}

	driver, err := mysql_migrate.WithInstance(db, &mysql_migrate.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}

	if err := applyMigrations(m); err != nil {
		return err
	}

	slog.Info("Database migrations applied successfully")
	return nil
}

// applyMigrations runs m.Up() and handles the ErrDirty case that arises when
// a previous run failed mid-migration. All migrations use CREATE TABLE IF NOT
// EXISTS, making them idempotent, so it is safe to force-clear the dirty
// version and retry rather than requiring a manual database intervention.
func applyMigrations(m *migrate.Migrate) error {
	err := m.Up()
	if err == nil || errors.Is(err, migrate.ErrNoChange) {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("Database migrations: No changes required")
		}
		return nil
	}

	// Dirty database: a prior run failed mid-migration and left the version
	// flagged as dirty. Force-clear it so we can retry cleanly.
	var dirtyErr migrate.ErrDirty
	if errors.As(err, &dirtyErr) {
		slog.Warn("Dirty migration state detected, forcing clean and retrying",
			"dirty_version", dirtyErr.Version)

		if forceErr := m.Force(dirtyErr.Version); forceErr != nil {
			return fmt.Errorf("failed to force dirty version %d: %w", dirtyErr.Version, forceErr)
		}

		// Retry after clearing dirty flag
		if retryErr := m.Up(); retryErr != nil && !errors.Is(retryErr, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations after dirty state fix: %w", retryErr)
		}
		return nil
	}

	return fmt.Errorf("failed to apply migrations: %w", err)
}
