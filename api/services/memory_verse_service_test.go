package services

import (
	"context"
	"regexp"
	"testing"
	"time"

	"discipleship_journal_api/models"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetPacks(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	service := NewMemoryVerseService(mock)
	userID := uuid.New()

	sysID := uuid.New()
	sysTitle := "System Pack"
	sysIdentifier := "SYS"
	sysDesc := "Desc"
	sysPublic := true
	sysCreated := time.Now()
	sysUpdated := time.Now()
	sysCount := 5

	// Test System Packs
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT vp.id, vp.user_id, vp.title, vp.identifier, vp.description, vp.is_public, vp.created_at, vp.updated_at,
		       (SELECT count(*) FROM memory_verses mv WHERE mv.verse_pack_id = vp.id) as verse_count
		FROM verse_packs vp WHERE vp.user_id IS NULL ORDER BY vp.title`)).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "identifier", "description", "is_public", "created_at", "updated_at", "verse_count"}).
			AddRow(sysID, nil, sysTitle, &sysIdentifier, &sysDesc, sysPublic, sysCreated, sysUpdated, sysCount))

	packs, err := service.GetPacks(context.Background(), userID, "system")
	assert.NoError(t, err)
	assert.Len(t, packs, 1)
	assert.Equal(t, "System Pack", packs[0].Title)

	// Test User Packs
	userPackID := uuid.New()
	userPackTitle := "My Pack"
	var userPackIdent *string = nil
	var userPackDesc *string = nil
	userPackPublic := false
	userPackCreated := time.Now()
	userPackUpdated := time.Now()
	userPackCount := 2

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT vp.id, vp.user_id, vp.title, vp.identifier, vp.description, vp.is_public, vp.created_at, vp.updated_at,
		       (SELECT count(*) FROM memory_verses mv WHERE mv.verse_pack_id = vp.id) as verse_count
		FROM verse_packs vp WHERE vp.user_id = $1 ORDER BY vp.created_at DESC`)).
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "identifier", "description", "is_public", "created_at", "updated_at", "verse_count"}).
			AddRow(userPackID, &userID, userPackTitle, userPackIdent, userPackDesc, userPackPublic, userPackCreated, userPackUpdated, userPackCount))

	packs, err = service.GetPacks(context.Background(), userID, "user")
	assert.NoError(t, err)
	assert.Len(t, packs, 1)
	assert.Equal(t, "My Pack", packs[0].Title)
}

func TestClonePack(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	service := NewMemoryVerseService(mock)
	userID := uuid.New()
	packID := uuid.New()
	identifier := "SRC"
	sourcePack := &models.VersePack{
		ID: packID, Title: "Source Pack", Identifier: identifier, IsPublic: true,
	}

	// 1. GetPack
	mock.ExpectQuery(`SELECT .* FROM verse_packs vp WHERE vp.id = \$1 .*`).
		WithArgs(packID, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "identifier", "description", "is_public", "created_at", "updated_at", "verse_count"}).
			AddRow(sourcePack.ID, nil, sourcePack.Title, &sourcePack.Identifier, nil, sourcePack.IsPublic, time.Now(), time.Now(), 10))

	// 2. CreatePack
	// Expect empty string for description (not nil) because struct field is string.
	mock.ExpectExec(`INSERT INTO verse_packs`).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), "New Title", "SRC", "", false, pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// 3. GetVerses (Source) with new query
	var verseTitle *string = nil // No title

	// We use a loose regex match for the complex query to avoid whitespace issues
	// Match: SELECT ... FROM memory_verses mv LEFT JOIN user_verse_preferences ... LEFT JOIN users ... WHERE mv.verse_pack_id = $1 ...
	mock.ExpectQuery(`SELECT .* FROM memory_verses mv LEFT JOIN user_verse_preferences uvp .* LEFT JOIN users u .* WHERE mv.verse_pack_id = \$1 ORDER BY mv.created_at ASC`).
		WithArgs(packID, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "verse_pack_id", "reference", "title", "effective_version", "tags", "created_at", "updated_at"}).
			AddRow(uuid.New(), packID, "John 3:16", verseTitle, "ESV", []byte(`["Love"]`), time.Now(), time.Now()))

	// 4. CreateVerse (Clone)
	mock.ExpectExec(`INSERT INTO memory_verses`).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), "John 3:16", "", "ESV", []byte(`["Love"]`), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	newPack, err := service.ClonePack(context.Background(), packID, userID, "New Title")
	assert.NoError(t, err)
	assert.Equal(t, "New Title", newPack.Title)
	assert.Equal(t, 1, newPack.VerseCount)
}

func TestSearchVerses(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	service := NewMemoryVerseService(mock)
	userID := uuid.New()
	query := "John"

	verseID := uuid.New()
	packID := uuid.New()
	packTitle := "My Pack"
	var verseTitle *string = nil

	// Updated regex for SearchVerses with joins
	// SELECT ... FROM memory_verses mv JOIN verse_packs vp ... LEFT JOIN user_verse_preferences ... LEFT JOIN users ... WHERE ...
	mock.ExpectQuery(`SELECT .* FROM memory_verses mv JOIN verse_packs vp .* LEFT JOIN user_verse_preferences uvp .* LEFT JOIN users u .* WHERE .*`).
		WithArgs(userID, "%"+query+"%").
		WillReturnRows(pgxmock.NewRows([]string{"id", "verse_pack_id", "reference", "title", "effective_version", "tags", "pack_title", "created_at", "updated_at"}).
			AddRow(verseID, packID, "John 3:16", verseTitle, "ESV", []byte(`["Love"]`), packTitle, time.Now(), time.Now()))

	verses, err := service.SearchVerses(context.Background(), userID, query)
	assert.NoError(t, err)
	assert.Len(t, verses, 1)
	assert.Equal(t, "John 3:16", verses[0].Reference)
	assert.Equal(t, packTitle, verses[0].PackTitle)
}
