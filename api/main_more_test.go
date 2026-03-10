package main

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMockDatabaseMethods(t *testing.T) {
	mock := &mockDatabase{}
	ctx := context.Background()

	_, err := mock.QueryContext(ctx, "query")
	assert.Error(t, err)

	row := mock.QueryRowContext(ctx, "query")
	assert.Nil(t, row) // Checking if it actually returns nil

	_, err = mock.ExecContext(ctx, "query")
	assert.Error(t, err)

	_, err = mock.BeginTx(ctx, nil)
	assert.Error(t, err)
}
