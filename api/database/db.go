package database

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// BuildConnectionString builds a MySQL connection string (DSN) from individual components
func BuildConnectionString() (string, error) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	if username == "" || dbName == "" {
		return "", fmt.Errorf("DB_USERNAME and DB_NAME environment variables are required")
	}

	if password == "" {
		return "", fmt.Errorf("DB_PASSWORD is required")
	}

	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "4000" // TiDB default
	}

	// Format: user:password@tcp(host:port)/dbname?tls=true&parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=true&parseTime=true", username, password, host, port, dbName)
	return dsn, nil
}

func Connect() error {
	dsn, err := BuildConnectionString()
	if err != nil {
		return fmt.Errorf("failed to build database connection string: %w", err)
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		// Mask password in error message for security
		errMsg := err.Error()
		if password := os.Getenv("DB_PASSWORD"); password != "" {
			errMsg = strings.ReplaceAll(errMsg, password, "***")
		}
		return fmt.Errorf("unable to open database: %s", errMsg)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(2)
	DB.SetConnMaxLifetime(time.Hour)
	DB.SetConnMaxIdleTime(30 * time.Minute)

	// Create a context with a timeout for the initial ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	fmt.Println("Connected to TiDB database")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
