package database

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all pending database migrations
func RunMigrations() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = os.Getenv("DB_USERNAME") // Fallback
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	cloudSQLInstance := os.Getenv("CLOUD_SQL_INSTANCE")

	// Construct connection string securely with URL encoding
	userInfo := url.UserPassword(dbUser, dbPassword)
	var dbURL string

	// Check if we are using Cloud SQL via Unix socket (common in Cloud Run)
	if cloudSQLInstance != "" {
		// Socket connection format: postgres://user:password@localhost/dbname?host=/cloudsql/instance
		socketPath := fmt.Sprintf("/cloudsql/%s", cloudSQLInstance)

		// Build URL query parameters
		query := url.Values{}
		query.Add("host", socketPath)
		query.Add("sslmode", "disable")

		u := url.URL{
			Scheme:   "postgres",
			User:     userInfo,
			Host:     "localhost", // Host is ignored when using unix socket in query, but required for URL parsing
			Path:     dbName,
			RawQuery: query.Encode(),
		}
		dbURL = u.String()
	} else {
		// Standard TCP connection
		if dbHost == "" {
			dbHost = "localhost"
		}
		if dbPort == "" {
			dbPort = "5432"
		}

		u := url.URL{
			Scheme:   "postgres",
			User:     userInfo,
			Host:     fmt.Sprintf("%s:%s", dbHost, dbPort),
			Path:     dbName,
			RawQuery: "sslmode=disable",
		}
		dbURL = u.String()
	}

	// We mask the password for logging purposes (although url.UserPassword might handle stringification safely, being explicit is safer)
	maskedUserInfo := url.UserPassword(dbUser, "***")
	var maskedURL string
	if cloudSQLInstance != "" {
		u := url.URL{Scheme: "postgres", User: maskedUserInfo, Host: "localhost", Path: dbName, RawQuery: fmt.Sprintf("host=/cloudsql/%s", cloudSQLInstance)}
		maskedURL = u.String()
	} else {
		u := url.URL{Scheme: "postgres", User: maskedUserInfo, Host: fmt.Sprintf("%s:%s", dbHost, dbPort), Path: dbName}
		maskedURL = u.String()
	}

	slog.Info("Running database migrations...", "url_masked", maskedURL)

	// Source URL: file://migrations
	// This assumes the migrations folder is in the current working directory
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}
	defer m.Close()

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
