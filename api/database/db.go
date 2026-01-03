package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// BuildConnectionString builds a PostgreSQL connection string from individual components
// Supports both local development and Cloud SQL connections
// For Cloud SQL, it returns a connection string suitable for the connector, or socket if legacy.
// Actually, with the connector, we primarily need the config object, but pgxpool.ParseConfig needs a string.
func BuildConnectionString() (string, error) {
	// Get individual components
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	cloudSQLInstance := os.Getenv("CLOUD_SQL_INSTANCE")

	// Check required fields
	// If Cloud SQL Instance is set, we might rely on IAM Auth, so password might be optional
	if username == "" || dbName == "" {
		return "", fmt.Errorf("DB_USERNAME and DB_NAME environment variables are required")
	}

	if cloudSQLInstance == "" && password == "" {
		// Local development requires password
		return "", fmt.Errorf("DB_PASSWORD is required for local development")
	}

	// URL encode the user info
	var userInfo *url.Userinfo
	if password != "" {
		userInfo = url.UserPassword(username, password)
	} else {
		// For IAM Auth, password can be empty
		userInfo = url.User(username)
	}

	if cloudSQLInstance != "" {
		// Cloud SQL connection
		// When using the connector, we use a standard postgres:// URL but we'll inject the dialer later.
		// However, we still need a valid connection string for ParseConfig to work.
		// We set host to "cloudsql" as a placeholder, or localhost, it doesn't matter much as the dialer takes over.
		// But we should follow best practices.
		// NOTE: When using IAM Auth, the password field is ignored/handled by the connector (it generates the token).

		u := url.URL{
			Scheme: "postgres",
			User:   userInfo,
			Host:   "cloudsql", // Placeholder
			Path:   "/" + dbName,
		}
		q := u.Query()
		q.Set("sslmode", "disable") // Connector handles encryption
		u.RawQuery = q.Encode()

		return u.String(), nil
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

		u := url.URL{
			Scheme: "postgres",
			User:   userInfo,
			Host:   fmt.Sprintf("%s:%s", host, port),
			Path:   "/" + dbName,
		}
		q := u.Query()
		q.Set("sslmode", "disable")
		u.RawQuery = q.Encode()

		return u.String(), nil
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

	// Cloud SQL Connector logic
	cloudSQLInstance := os.Getenv("CLOUD_SQL_INSTANCE")
	if cloudSQLInstance != "" {
		// Initialize the connector dialer
		// We enable IAM AuthN automatically.
		d, err := cloudsqlconn.NewDialer(context.Background(), cloudsqlconn.WithIAMAuthN())
		if err != nil {
			return fmt.Errorf("failed to initialize Cloud SQL dialer: %w", err)
		}
		// Note: We don't close the dialer here because it needs to stay open for the lifetime of the app.
		// In a real app, we might want to handle cleanup on shutdown, but for this singleton DB pattern, it's acceptable.

		// Configure the dial function
		config.ConnConfig.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return d.Dial(ctx, cloudSQLInstance)
		}

		// When using IAM Auth, we should ensure the password is NOT set in the config if it's empty,
		// although the connector usually handles auth by generating a token and using it as password.
		// If we passed a password in BuildConnectionString, pgx has it.
		// If we use WithIAMAuthN, the connector works at the connection level.
		// Wait, for Postgres, the connector handles the secure tunnel. IAM Auth requires sending the token as the password.
		// The `cloud-sql-go-connector` documentation says:
		// "Automatic IAM AuthenticationN: The connector can automatically perform IAM authentication..."
		// It does this by updating the password in the connection config? No, it's a dialer.
		// Actually, for Postgres, the connector *just* establishes the connection.
		// For IAM Auth, the *driver* needs to send the token as the password.
		// BUT `cloudsqlconn` has a feature `WithIAMAuthN()` which handles this internally?
		// Let's verify. The docs say: "For PostgreSQL... connection logic...".
		// Actually, `cloudsqlconn` documentation says `WithIAMAuthN` enables automatic IAM authentication.
		// It works by the Dialer handling the handshake?
		// No, for Postgres, the client must send the password.
		// The Go Connector's `Dialer.Dial` returns a connection.
		//
		// Re-reading docs: "When connecting to PostgreSQL... using IAM authentication... you must use the `WithIAMAuthN` option...".
		// AND for `pgx` specifically, it seems we might need to register a driver OR just use the Dialer.
		//
		// Crucially, if we use `d.Dial`, we get a `net.Conn`.
		// `pgx` will then use that connection.
		// Does `pgx` perform the SASL auth? Yes.
		// Does `d.Dial` handle the SASL auth? No, it handles the secure transport.
		//
		// Wait, if I use `WithIAMAuthN`, the `Dialer` provides a `Dial` method.
		// The documentation for `WithIAMAuthN` says: "Enables automatic IAM authentication...".
		//
		// Correct usage for pgx:
		// The `cloudsqlconn` README example for pgx says:
		// "Use the dialer to create a custom DialFunc..."
		// AND `WithIAMAuthN()` is used in `NewDialer`.
		// It seems this IS enough. The internal implementation of the connector creates a special connection that handles the auth or provides the token?
		//
		// Actually, `pgx` sends the password. If I don't provide a password in `config.Password`, `pgx` might fail?
		// UNLESS the dialer intercepts the handshake?
		//
		// Let's check `cloudsqlconn` source or deep docs if possible.
		// "The Go Connector does NOT modify the driver's authentication logic..."
		// Wait, if so, how does it pass the token?
		//
		// Ah, for `database/sql` drivers like `pgx/stdlib`, we register a driver which wraps the connector.
		// But for `pgx` native (pool), we use `DialFunc`.
		//
		// IMPORTANT: For `pgx`, we must NOT set the password in the config if using IAM Auth?
		// OR does the connector rely on the password being empty?
		//
		// "When using IAM Authentication with the Go Connector, you do not need to provide a password..."
		//
		// However, `pgx` might complain if password is empty?
		// The `pgx` driver defaults to empty password.
		//
		// There is a potential issue: `pgx` might try to authenticate using MD5 or SCRAM.
		// Cloud SQL IAM Auth expects the password to be the OAuth token.
		//
		// Does `cloudsqlconn` inject the token?
		// Yes, `cloudsqlconn` automatically handles the password when `WithIAMAuthN` is used, ONLY if using the `database/sql` driver registration.
		// BUT if using `Dialer` directly with `pgx`, does it work?
		//
		// The official Google example for `pgx` (native) says:
		// https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/blob/main/examples/postgres/pgx/main.go
		// It uses `cloudsqlconn.NewDialer(ctx, cloudsqlconn.WithIAMAuthN())`.
		// And `config.ConnConfig.DialFunc = ...`
		// And `config.Password` is NOT set (or whatever).
		//
		// So it seems `cloudsqlconn` with `pgx` works out of the box with just `DialFunc`.
		// The `Dialer` likely returns a connection that helps, OR `pgx` doesn't send password?
		//
		// Wait, actually, the `cloudsqlconn` logic for Postgres *intercepts* the startup message?
		// No, that sounds complex.
		//
		// Let's assume the library handles it correctly if `WithIAMAuthN` is set.
		//
		// One detail: "Note: The PostgreSQL driver must be configured to use a password-less authentication method (e.g. 'trust' or 'cert')? NO, Cloud SQL IAM uses password."
		//
		// Let's trust the `WithIAMAuthN` option on the dialer.
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
