package services

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBibleVersionService_GetVersions(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	_ = mockDB
	require.NoError(t, err)
	defer db.Close()

	service := NewBibleVersionService(db)

	rows := sqlmock.NewRows([]string{
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

func TestBibleVersionService_ScrapeVersions(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `
			<html>
			<body>
				<select name="version" class="search-dropdown">
					<option value="ESV">English Standard Version</option>
					<option value="NIV">New International Version</option>
					<option value="">Select Version</option>
				</select>
			</body>
			</html>
		`)
	}))
	defer ts.Close()

	// Instantiate service manually to inject dependencies
	service := &bibleVersionService{
		db:         nil, // Not needed for scraping
		httpClient: ts.Client(),
		scrapeURL:  ts.URL,
	}

	versions, err := service.ScrapeVersions(context.Background())
	assert.NoError(t, err)
	assert.Len(t, versions, 2)
	assert.Equal(t, "English Standard Version", versions["ESV"])
	assert.Equal(t, "New International Version", versions["NIV"])
}

func TestBibleVersionService_SyncVersions(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `
			<select name="version" class="search-dropdown">
				<option value="ESV">English Standard Version</option>
			</select>
		`)
	}))
	defer ts.Close()

	db, mockDB, err := sqlmock.New()
	_ = mockDB
	require.NoError(t, err)
	defer db.Close()

	service := &bibleVersionService{
		db:         db,
		httpClient: ts.Client(),
		scrapeURL:  ts.URL,
	}

	// Expect transaction
	mockDB.ExpectBegin()
	// Expect insert
	mockDB.ExpectExec("INSERT INTO bible_versions").
		WithArgs("English Standard Version", "ESV").
		WillReturnResult(sqlmock.NewResult(1, 1))
	// Expect commit
	mockDB.ExpectCommit()

	err = service.SyncVersions(context.Background())
	assert.NoError(t, err)

	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestBibleVersionService_ScrapeVersions_Error(t *testing.T) {
	// Mock server returning 500
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	service := &bibleVersionService{
		httpClient: ts.Client(),
		scrapeURL:  ts.URL,
	}

	versions, err := service.ScrapeVersions(context.Background())
	assert.Error(t, err)
	assert.Nil(t, versions)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}

func TestBibleVersionService_SyncVersions_ScrapeError(t *testing.T) {
	// Mock server returning 500
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	db, mockDB, err := sqlmock.New()
	_ = mockDB
	require.NoError(t, err)
	defer db.Close()

	service := &bibleVersionService{
		db:         db,
		httpClient: ts.Client(),
		scrapeURL:  ts.URL,
	}

	err = service.SyncVersions(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}

func TestBibleVersionService_SyncVersions_DBError(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `
			<select name="version" class="search-dropdown">
				<option value="ESV">English Standard Version</option>
			</select>
		`)
	}))
	defer ts.Close()

	t.Run("begin_error", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		_ = mockDB
		require.NoError(t, err)
		defer db.Close()

		service := &bibleVersionService{
			db:         db,
			httpClient: ts.Client(),
			scrapeURL:  ts.URL,
		}

		mockDB.ExpectBegin().WillReturnError(fmt.Errorf("begin error"))

		err = service.SyncVersions(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "begin error")
	})

	t.Run("exec_error", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		_ = mockDB
		require.NoError(t, err)
		defer db.Close()

		service := &bibleVersionService{
			db:         db,
			httpClient: ts.Client(),
			scrapeURL:  ts.URL,
		}

		mockDB.ExpectBegin()
		mockDB.ExpectExec("INSERT INTO bible_versions").
			WithArgs("English Standard Version", "ESV").
			WillReturnError(fmt.Errorf("exec error"))
		mockDB.ExpectRollback()

		err = service.SyncVersions(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exec error")
	})

	t.Run("commit_error", func(t *testing.T) {
		db, mockDB, err := sqlmock.New()
		_ = mockDB
		require.NoError(t, err)
		defer db.Close()

		service := &bibleVersionService{
			db:         db,
			httpClient: ts.Client(),
			scrapeURL:  ts.URL,
		}

		mockDB.ExpectBegin()
		mockDB.ExpectExec("INSERT INTO bible_versions").
			WithArgs("English Standard Version", "ESV").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.ExpectCommit().WillReturnError(fmt.Errorf("commit error"))

		err = service.SyncVersions(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "commit error")
	})
}
