package main

import (
	"context"
	"fmt"

	"database/sql"
	_ "discipleship_journal_api/docs"
)

// @title Discipleship Journal API
// @version 1.0
// @description API for the Discipleship Journal application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// mockDatabase implements DBInterface but returns errors for all operations
// Used when database connection fails
type mockDatabase struct {
	connected bool
}

func (m *mockDatabase) QueryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error) {
	return nil, fmt.Errorf("database not connected")
}

func (m *mockDatabase) QueryRowContext(ctx context.Context, sql string, args ...any) *sql.Row {
	// Not easily mockable in standard library without a real connection, return nil
	// In reality this might panic if a handler calls Scan() on a nil pointer,
	// but the application logic should check for connected state earlier.
	return nil
}

func (m *mockDatabase) ExecContext(ctx context.Context, sql string, arguments ...any) (sql.Result, error) {
	return nil, fmt.Errorf("database not connected")
}

func (m *mockDatabase) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return nil, fmt.Errorf("database not connected")
}
