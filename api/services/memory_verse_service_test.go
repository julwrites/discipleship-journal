package services

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"testing"
	"time"

	"database/sql"
	"discipleship_journal_api/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetPacks(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := NewMemoryVerseService(db)
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
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "identifier", "description", "is_public", "created_at", "updated_at", "verse_count"}).
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
		FROM verse_packs vp WHERE vp.user_id = ? ORDER BY vp.created_at DESC`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "identifier", "description", "is_public", "created_at", "updated_at", "verse_count"}).
			AddRow(userPackID, &userID, userPackTitle, userPackIdent, userPackDesc, userPackPublic, userPackCreated, userPackUpdated, userPackCount))

	packs, err = service.GetPacks(context.Background(), userID, "user")
	assert.NoError(t, err)
	assert.Len(t, packs, 1)
	assert.Equal(t, "My Pack", packs[0].Title)
}

func TestClonePack(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := NewMemoryVerseService(db)
	userID := uuid.New()
	packID := uuid.New()
	identifier := "SRC"
	sourcePack := &models.VersePack{
		ID: packID, Title: "Source Pack", Identifier: identifier, IsPublic: true,
	}

	// 1. GetPack
	mock.ExpectQuery(`SELECT .* FROM verse_packs vp WHERE vp.id = \? .*`).
		WithArgs(packID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "identifier", "description", "is_public", "created_at", "updated_at", "verse_count"}).
			AddRow(sourcePack.ID, nil, sourcePack.Title, &sourcePack.Identifier, nil, sourcePack.IsPublic, time.Now(), time.Now(), 10))

	// 2. CreatePack
	// Expect empty string for description (not nil) because struct field is string.
	mock.ExpectExec(`INSERT INTO verse_packs`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "New Title", "SRC", "", false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 3. GetVerses (Source) with new query
	var verseTitle *string = nil // No title

	// We use a loose regex match for the complex query to avoid whitespace issues
	// Match: SELECT ... FROM memory_verses mv LEFT JOIN user_verse_preferences ... LEFT JOIN users ... WHERE mv.verse_pack_id = ? ...
	mock.ExpectQuery(`SELECT .* FROM memory_verses mv LEFT JOIN user_verse_preferences uvp .* LEFT JOIN users u .* WHERE mv.verse_pack_id = \? ORDER BY mv.created_at ASC`).
		WithArgs(packID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "verse_pack_id", "reference", "title", "effective_version", "version_source", "tags", "created_at", "updated_at"}).
			AddRow(uuid.New(), packID, "John 3:16", verseTitle, "ESV", "original", []byte(`["Love"]`), time.Now(), time.Now()))

	// 4. CreateVerse (Clone)
	mock.ExpectExec(`INSERT INTO memory_verses`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "John 3:16", "", "ESV", []byte(`["Love"]`), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	newPack, err := service.ClonePack(context.Background(), packID, userID, "New Title", true)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", newPack.Title)
	assert.Equal(t, 1, newPack.VerseCount)
}

func TestClonePack_Original(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := NewMemoryVerseService(db)
	userID := uuid.New()
	packID := uuid.New()
	identifier := "SRC"
	sourcePack := &models.VersePack{
		ID: packID, Title: "Source Pack", Identifier: identifier, IsPublic: true,
	}

	// 1. GetPack
	mock.ExpectQuery(`SELECT .* FROM verse_packs vp WHERE vp.id = \? .*`).
		WithArgs(packID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "identifier", "description", "is_public", "created_at", "updated_at", "verse_count"}).
			AddRow(sourcePack.ID, nil, sourcePack.Title, &sourcePack.Identifier, nil, sourcePack.IsPublic, time.Now(), time.Now(), 10))

	// 2. CreatePack
	mock.ExpectExec(`INSERT INTO verse_packs`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "New Title", "SRC", "", false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 3. GetOriginalVerses (Source)
	var verseTitle *string = nil

	// Match simple query: SELECT ... FROM memory_verses mv WHERE mv.verse_pack_id = ? ORDER BY mv.created_at ASC
	mock.ExpectQuery(`SELECT .* FROM memory_verses mv WHERE mv.verse_pack_id = \? ORDER BY mv.created_at ASC`).
		WithArgs(packID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "verse_pack_id", "reference", "title", "version", "version_source", "tags", "created_at", "updated_at"}).
			AddRow(uuid.New(), packID, "John 3:16", verseTitle, "KJV", "original", []byte(`["Love"]`), time.Now(), time.Now()))

	// 4. CreateVerse (Clone) - Should use KJV from original
	mock.ExpectExec(`INSERT INTO memory_verses`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "John 3:16", "", "KJV", []byte(`["Love"]`), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	newPack, err := service.ClonePack(context.Background(), packID, userID, "New Title", false)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", newPack.Title)
	assert.Equal(t, 1, newPack.VerseCount)
}

func TestSearchVerses(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := NewMemoryVerseService(db)
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
		WillReturnRows(sqlmock.NewRows([]string{"id", "verse_pack_id", "reference", "title", "effective_version", "version_source", "tags", "pack_title", "created_at", "updated_at"}).
			AddRow(verseID, packID, "John 3:16", verseTitle, "ESV", "original", []byte(`["Love"]`), packTitle, time.Now(), time.Now()))

	verses, err := service.SearchVerses(context.Background(), userID, query)
	assert.NoError(t, err)
	assert.Len(t, verses, 1)
	assert.Equal(t, "John 3:16", verses[0].Reference)
	assert.Equal(t, packTitle, verses[0].PackTitle)
}

func TestDeletePack(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := NewMemoryVerseService(db)
	userID := uuid.New()
	packID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM verse_packs").
			WithArgs(packID, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := service.DeletePack(context.Background(), packID, userID)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM verse_packs").
			WithArgs(packID, userID).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := service.DeletePack(context.Background(), packID, userID)
		assert.Error(t, err)
		assert.Equal(t, models.ErrNotFound, err)
	})
}

func TestDeleteVerse(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := NewMemoryVerseService(db)
	userID := uuid.New()
	verseID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM memory_verses mv USING verse_packs vp").
			WithArgs(verseID, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := service.DeleteVerse(context.Background(), verseID, userID)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM memory_verses mv USING verse_packs vp").
			WithArgs(verseID, userID).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := service.DeleteVerse(context.Background(), verseID, userID)
		assert.Error(t, err)
		assert.Equal(t, models.ErrNotFound, err)
	})
}

func TestUpdateVerse(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := NewMemoryVerseService(db)
	userID := uuid.New()
	verseID := uuid.New()
	verse := &models.MemoryVerse{
		ID:        verseID,
		Reference: "John 3:16",
		Title:     "The Gospel",
		Version:   "ESV",
		Tags:      []string{"Love"},
	}
	tagsJSON, _ := json.Marshal(verse.Tags)

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE memory_verses mv SET").
			WithArgs(verseID, verse.Reference, verse.Title, verse.Version, tagsJSON, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := service.UpdateVerse(context.Background(), verse, userID)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE memory_verses mv SET").
			WithArgs(verseID, verse.Reference, verse.Title, verse.Version, tagsJSON, userID).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := service.UpdateVerse(context.Background(), verse, userID)
		assert.Error(t, err)
		assert.Equal(t, models.ErrNotFound, err)
	})
}

func TestSetVersePreference(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewMemoryVerseService(db)
	userID := uuid.New()
	verseID := uuid.New()
	version := "NIV"

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO user_verse_preferences (user_id, verse_id, version_override)`)).
		WithArgs(userID, verseID, version).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = service.SetVersePreference(context.Background(), userID, verseID, version)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetVersePreferencesBatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewMemoryVerseService(db)
	userID := uuid.New()
	verseID1 := uuid.New()
	verseID2 := uuid.New()
	verseIDs := []uuid.UUID{verseID1, verseID2}
	version := "ESV"

	mock.ExpectExec(`INSERT INTO user_verse_preferences .*`).
		WithArgs(userID, verseIDs[0], version, userID, verseIDs[1], version).
		WillReturnResult(sqlmock.NewResult(1, 2))

	err = service.SetVersePreferencesBatch(context.Background(), userID, verseIDs, version)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClonePack_SourceNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := NewMemoryVerseService(db)
	userID := uuid.New()
	packID := uuid.New()

	// 1. GetPack - Error
	mock.ExpectQuery(`SELECT .* FROM verse_packs vp WHERE vp.id = \? .*`).
		WithArgs(packID, userID).
		WillReturnError(sql.ErrNoRows)

	newPack, err := service.ClonePack(context.Background(), packID, userID, "New Title", true)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.Nil(t, newPack)
}

func TestClonePack_DBCreationError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := NewMemoryVerseService(db)
	userID := uuid.New()
	packID := uuid.New()
	sourcePack := &models.VersePack{
		ID: packID, Title: "Source Pack", Identifier: "SRC", IsPublic: true,
	}

	// 1. GetPack - Success
	mock.ExpectQuery(`SELECT .* FROM verse_packs vp WHERE vp.id = \? .*`).
		WithArgs(packID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "identifier", "description", "is_public", "created_at", "updated_at", "verse_count"}).
			AddRow(sourcePack.ID, nil, sourcePack.Title, &sourcePack.Identifier, nil, sourcePack.IsPublic, time.Now(), time.Now(), 10))

	// 2. CreatePack - Error
	mock.ExpectExec(`INSERT INTO verse_packs`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "New Title", "SRC", "", false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("db insert error"))

	newPack, err := service.ClonePack(context.Background(), packID, userID, "New Title", true)
	assert.Error(t, err)
	assert.EqualError(t, err, "db insert error")
	assert.Nil(t, newPack)
}

func TestRemoveVersePreference(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewMemoryVerseService(db)
	userID := uuid.New()
	verseID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM user_verse_preferences WHERE user_id = ? AND verse_id = ?`)).
		WithArgs(userID, verseID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = service.RemoveVersePreference(context.Background(), userID, verseID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
