package services

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBibleVersionService_GetVersions(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	service := NewBibleVersionService(mockDB)

	rows := pgxmock.NewRows([]string{
		"id", "name", "abbreviation", "created_at", "updated_at",
	}).AddRow(
		"uuid-1", "English Standard Version", "ESV", time.Now(), time.Now(),
	).AddRow(
		"uuid-2", "New International Version", "NIV", time.Now(), time.Now(),
	)

	mockDB.ExpectQuery("SELECT id, name, abbreviation, created_at, updated_at FROM bible_versions ORDER BY name ASC").
		WillReturnRows(rows)

	versions, err := service.GetVersions(context.Background())
	assert.NoError(t, err)
	assert.Len(t, versions, 2)
	assert.Equal(t, "ESV", versions[0].Abbreviation)
	assert.Equal(t, "NIV", versions[1].Abbreviation)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}
