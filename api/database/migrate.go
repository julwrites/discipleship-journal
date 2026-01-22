package database

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"discipleship_journal_api/migrations"
	"github.com/golang-migrate/migrate/v4"
	pgx_migrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// RunMigrations applies all pending database migrations
func RunMigrations() error {
	// 1. Build connection string using shared logic
	dbURL, err := BuildConnectionString()
	if err != nil {
		return fmt.Errorf("failed to build connection string for migrations: %w", err)
	}

	// 2. Parse config for pgx
	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("failed to parse config from connection string: %w", err)
	}

	// 3. Configure Cloud SQL Connector if needed
	cloudSQLInstance := os.Getenv("CLOUD_SQL_INSTANCE")
	if cloudSQLInstance != "" {
		// Initialize the connector dialer with IAM AuthN
		d, err := cloudsqlconn.NewDialer(context.Background(), cloudsqlconn.WithIAMAuthN())
		if err != nil {
			return fmt.Errorf("failed to initialize Cloud SQL dialer for migrations: %w", err)
		}
		defer func() { _ = d.Close() }()

		config.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return d.Dial(ctx, cloudSQLInstance)
		}
	}

	// 4. Open database connection using pgx stdlib
	db := stdlib.OpenDB(*config)
	defer func() { _ = db.Close() }()

	// 5. Verify connection before starting migrate
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database before migrations: %w", err)
	}

	slog.Info("Starting database migrations...", "instance", cloudSQLInstance)

	// Use embedded migrations
	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs source driver: %w", err)
	}

	// Create migrate driver from existing DB connection
	driver, err := pgx_migrate.WithInstance(db, &pgx_migrate.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// Use NewWithInstance instead of NewWithDatabaseInstance to support source driver
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx", driver)
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
