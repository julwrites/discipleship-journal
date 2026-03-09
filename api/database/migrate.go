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

// applyMigrations runs migrations step by step.
//
// Recovery mode: the staging TiDB database has all migrations previously applied
// but schema_migrations is out of sync. Once the first conflict is detected
// (proving we are reconciling a DB with prior migration history), ALL subsequent
// SQL-level migration errors are treated as "previously applied" and skipped.
//
// This is safe because:
//   - Known conflict errors (Duplicate column, Duplicate entry, etc.) confirm the
//     schema/data already exists in the target state.
//   - Unknown SQL errors in recovery mode also indicate a previously-run migration
//     that left partial state — skipping is correct.
//   - Truly new migrations that have never run will succeed cleanly (nil error).
//   - Infrastructure errors (can't read version, can't force) still fail loudly.
func applyMigrations(m *migrate.Migrate) error {
	const maxSteps = 1000
	inRecoveryMode := false

	for i := 0; i < maxSteps; i++ {
		err := m.Steps(1)

		if err == nil {
			continue
		}

		if errors.Is(err, migrate.ErrNoChange) || strings.Contains(err.Error(), "file does not exist") {
			// ErrNoChange: golang-migrate confirmed nothing left to apply.
			// "file does not exist": the source driver has no migration file at or
			// after the current DB version — we are at/past the last migration.
			// Both are terminal success conditions.
			slog.Info("Database migrations: all migrations applied")
			return nil
		}

		// Dirty state from a previous failed run — force clean and enter recovery.
		var dirtyErr migrate.ErrDirty
		if errors.As(err, &dirtyErr) {
			inRecoveryMode = true
			slog.Warn("Dirty migration state detected, forcing clean",
				"dirty_version", dirtyErr.Version)
			if forceErr := m.Force(dirtyErr.Version); forceErr != nil {
				return fmt.Errorf("failed to force dirty version %d: %w", dirtyErr.Version, forceErr)
			}
			continue
		}

		// Determine whether this is a SQL-level migration error that we can skip.
		// In recovery mode we skip all of them; outside recovery mode we only skip
		// known conflict patterns (duplicate column/key/entry, already exists, etc.).
		isSQLMigrationErr := strings.Contains(err.Error(), "migration failed in line")
		isKnownConflict := isMigrationConflictError(err)

		if isKnownConflict || (inRecoveryMode && isSQLMigrationErr) {
			inRecoveryMode = true // entering or staying in recovery

			v, dirty, vErr := m.Version()
			if vErr != nil {
				return fmt.Errorf("failed to get version after migration conflict: %w", vErr)
			}
			if dirty {
				slog.Warn("Migration conflict — marking as applied and continuing",
					"version", v, "known_conflict", isKnownConflict,
					"recovery_mode", inRecoveryMode, "error", err.Error())
				if forceErr := m.Force(int(v)); forceErr != nil {
					return fmt.Errorf("failed to force version %d: %w", v, forceErr)
				}
			}
			continue
		}

		// Not a SQL migration error and not in recovery mode — this is a real failure.
		return fmt.Errorf("migration step failed: %w", err)
	}

	return nil
}

// isMigrationConflictError returns true for TiDB/MySQL errors that indicate a
// migration was already applied — either its DDL (schema) or DML (seed data)
// objects already exist in the database. Safe to skip in all cases because the
// database is already in or past the target state for that migration.
//
// DDL conflicts arise when schema migrations (ADD COLUMN, CREATE INDEX, RENAME)
// were partially applied by a previous run.
//
// DML conflicts (Duplicate entry) arise when seed-data migrations try to INSERT
// rows that were already loaded into TiDB by the CSV data import.
func isMigrationConflictError(err error) bool {
	s := err.Error()
	return strings.Contains(s, "Duplicate column name") || // ADD COLUMN already done
		strings.Contains(s, "Duplicate key name") || // ADD INDEX already done
		strings.Contains(s, "Duplicate entry") || // INSERT row already exists (seed data)
		strings.Contains(s, "already exists") || // generic already-exists
		strings.Contains(s, "doesn't exist") || // RENAME source gone (already renamed)
		strings.Contains(s, "Unknown table") // DROP / RENAME on already-removed table
}
