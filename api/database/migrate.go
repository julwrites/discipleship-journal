package database

import (
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

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("Database migrations: No changes required")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	slog.Info("Database migrations applied successfully")
	return nil
}
