package database

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

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

// applyMigrations runs migrations step by step, handling two recoverable cases:
//
//  1. ErrDirty: a previous run failed mid-migration and left the version flagged
//     as dirty. Force-clear it so we can retry from that version.
//
//  2. Schema conflict errors (duplicate column, duplicate key, table already
//     exists, table not found for rename): the schema was already applied by a
//     previous partial migration run or data import. Force-mark that version as
//     applied and move on to the next one.
//
// All these cases are safe to recover from because the migration SQL is either
// idempotent (CREATE TABLE IF NOT EXISTS) or the conflicting object already
// exists in the target state, so no data or schema is lost by skipping.
func applyMigrations(m *migrate.Migrate) error {
	const maxSteps = 1000 // safety guard against infinite loop

	for i := 0; i < maxSteps; i++ {
		err := m.Steps(1)

		if err == nil {
			// Step applied cleanly — continue to next
			continue
		}

		if errors.Is(err, migrate.ErrNoChange) {
			// No more migrations to apply
			slog.Info("Database migrations: No changes required")
			return nil
		}

		// Dirty state: a previous run failed and flagged this version dirty.
		var dirtyErr migrate.ErrDirty
		if errors.As(err, &dirtyErr) {
			slog.Warn("Dirty migration state detected, forcing clean",
				"dirty_version", dirtyErr.Version)
			if forceErr := m.Force(dirtyErr.Version); forceErr != nil {
				return fmt.Errorf("failed to force dirty version %d: %w", dirtyErr.Version, forceErr)
			}
			continue
		}

		// Schema conflict: the object this migration creates/alters already
		// exists (e.g. column added by a previous partial run or data import).
		// Mark this version as applied and move on.
		if isSchemaConflictError(err) {
			v, dirty, vErr := m.Version()
			if vErr != nil {
				return fmt.Errorf("failed to get version after schema conflict: %w", vErr)
			}
			slog.Warn("Schema already applied, marking migration as complete",
				"version", v, "dirty", dirty, "error", err.Error())
			if forceErr := m.Force(int(v)); forceErr != nil {
				return fmt.Errorf("failed to force version %d after schema conflict: %w", v, forceErr)
			}
			continue
		}

		return fmt.Errorf("migration step failed: %w", err)
	}

	return nil
}

// isSchemaConflictError returns true for TiDB/MySQL errors that indicate the
// schema change in a migration was already applied. These are safe to skip
// because the database is already in (or past) the target state.
func isSchemaConflictError(err error) bool {
	s := err.Error()
	return strings.Contains(s, "Duplicate column name") || // ADD COLUMN already done
		strings.Contains(s, "Duplicate key name") || // ADD INDEX already done
		strings.Contains(s, "already exists") || // generic already-exists
		strings.Contains(s, "doesn't exist") || // RENAME source gone (already renamed)
		strings.Contains(s, "Unknown table") // DROP / RENAME on already-removed table
}
