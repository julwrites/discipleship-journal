package database

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// BuildConnectionString builds a PostgreSQL connection string from individual components
// Supports both local development and Cloud SQL (Unix socket) connections
func BuildConnectionString() (string, error) {
	// Get individual components
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if username == "" || password == "" || dbName == "" {
		return "", fmt.Errorf("DB_USERNAME, DB_PASSWORD, and DB_NAME environment variables are required")
	}

	// Check if we're using Cloud SQL (Unix socket) or standard TCP
	cloudSQLInstance := os.Getenv("CLOUD_SQL_INSTANCE")

	if cloudSQLInstance != "" {
		// Cloud SQL Unix socket connection
		// Format: postgres://username:password@/dbname?host=/cloudsql/project:region:instance&sslmode=disable
		return fmt.Sprintf(
			"postgres://%s:%s@/%s?host=/cloudsql/%s&sslmode=disable",
			username, password, dbName, cloudSQLInstance,
		), nil
	} else {
		// Standard TCP connection (local/dev)
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")

		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "5432"
		}

		// Format: postgres://username:password@host:port/dbname?sslmode=disable
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			username, password, host, port, dbName,
		), nil
	}
}

func Connect() error {
	dbURL, err := BuildConnectionString()
	if err != nil {
		return fmt.Errorf("failed to build database connection string: %w", err)
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		// Mask password in error message for security
		errMsg := err.Error()
		if password := os.Getenv("DB_PASSWORD"); password != "" {
			errMsg = strings.ReplaceAll(errMsg, password, "***")
		}
		return fmt.Errorf("unable to create connection pool: %s", errMsg)
	}

	// Create a context with a timeout for the initial ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	fmt.Println("Connected to database")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
