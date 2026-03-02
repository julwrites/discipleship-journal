package services

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestNoteService_NilDB checks all methods for nil database connection error
func TestNoteService_NilDB(t *testing.T) {
	service := NewNoteService(nil)
	ctx := context.Background()

	t.Run("CreateNote", func(t *testing.T) {
		_, err := service.CreateNote(ctx, "uid", "title", nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection is nil")
	})

	t.Run("DeleteNote", func(t *testing.T) {
		err := service.DeleteNote(ctx, "uid", "nid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection is nil")
	})

	t.Run("UpdateNote", func(t *testing.T) {
		err := service.UpdateNote(ctx, "uid", "nid", "title", nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection is nil")
	})

	t.Run("GetNote", func(t *testing.T) {
		_, err := service.GetNote(ctx, "uid", "nid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection is nil")
	})

	t.Run("GetNotes", func(t *testing.T) {
		_, _, err := service.GetNotes(ctx, "uid", 1, 10, NoteFilter{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection is nil")
	})

	t.Run("CreateTag", func(t *testing.T) {
		_, err := service.CreateTag(ctx, "uid", "name")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection is nil")
	})

	t.Run("GetUserTags", func(t *testing.T) {
		_, err := service.GetUserTags(ctx, "uid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection is nil")
	})

	t.Run("DeleteTag", func(t *testing.T) {
		err := service.DeleteTag(ctx, "uid", "tid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection is nil")
	})
}

// TestNoteService_CreateNote_DBErrors checks transaction errors
func TestNoteService_CreateNote_DBErrors(t *testing.T) {
	ctx := context.Background()
	userID := "uid"
	title := "title"

	t.Run("Begin Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		mock.ExpectBegin().WillReturnError(errors.New("begin error"))
		_, err = service.CreateNote(ctx, userID, title, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, "begin error", err.Error())
	})

	t.Run("Insert Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO notes").
			WithArgs(sqlmock.AnyArg(), userID, title, json.RawMessage(nil), "active"). // content is nil, status defaults to "active"
			WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()
		_, err = service.CreateNote(ctx, userID, title, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, "insert error", err.Error())
	})

	t.Run("Commit Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO notes").
			WithArgs(sqlmock.AnyArg(), userID, title, json.RawMessage(nil), "active").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(regexp.QuoteMeta("SELECT created_at, updated_at FROM notes WHERE id = ?")).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).AddRow(time.Now(), time.Now()))

		mock.ExpectCommit().WillReturnError(errors.New("commit error"))
		mock.ExpectRollback()

		_, err = service.CreateNote(ctx, userID, title, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, "commit error", err.Error())
	})
}

// TestNoteService_UpdateNote_DBErrors checks transaction errors
func TestNoteService_UpdateNote_DBErrors(t *testing.T) {
	ctx := context.Background()
	userID := "uid"
	noteID := "nid"
	title := "title"

	t.Run("Begin Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		mock.ExpectBegin().WillReturnError(errors.New("begin error"))
		err = service.UpdateNote(ctx, userID, noteID, title, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, "begin error", err.Error())
	})

	t.Run("Delete Tags Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE notes").
			WithArgs(title, json.RawMessage(nil), noteID, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("DELETE FROM note_tags").
			WithArgs(noteID).
			WillReturnError(errors.New("delete tags error"))
		mock.ExpectRollback()

		err = service.UpdateNote(ctx, userID, noteID, title, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, "delete tags error", err.Error())
	})
}

// TestNoteService_GetNote_TagErrors checks tag fetching errors
func TestNoteService_GetNote_TagErrors(t *testing.T) {
	ctx := context.Background()
	userID := "uid"
	noteID := "nid"

	t.Run("Tag Query Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		mock.ExpectQuery("SELECT id, user_id, title").
			WithArgs(noteID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow(noteID, userID, "title", json.RawMessage("{}"), "active", time.Now(), time.Now(), nil))

		mock.ExpectQuery("SELECT .* FROM tags").
			WithArgs(noteID).
			WillReturnError(errors.New("tag query error"))

		_, err = service.GetNote(ctx, userID, noteID)
		assert.Error(t, err)
		assert.Equal(t, "tag query error", err.Error())
	})

	t.Run("Tag Scan Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		mock.ExpectQuery("SELECT id, user_id, title").
			WithArgs(noteID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow(noteID, userID, "title", json.RawMessage("{}"), "active", time.Now(), time.Now(), nil))

		// Return rows with wrong type/column count to force scan error
		mock.ExpectQuery("SELECT .* FROM tags").
			WithArgs(noteID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
				AddRow("tid", "uid", "name", "invalid-date")) // invalid date string scan error

		_, err = service.GetNote(ctx, userID, noteID)
		assert.Error(t, err)
		// Check that it is an error, exact message varies by driver/version
		assert.NotNil(t, err)
	})
}

// TestNoteService_GetNotes_SortAndFilter checks all filter paths
func TestNoteService_GetNotes_SortAndFilter(t *testing.T) {
	ctx := context.Background()
	userID := "uid"
	now := time.Now()

	t.Run("Full Filter", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		startDate := now.Add(-time.Hour)
		endDate := now.Add(time.Hour)
		filter := NoteFilter{
			SearchQuery: "search",
			Tag:         "tag",
			StartDate:   &startDate,
			EndDate:     &endDate,
			SortBy:      "created_at",
			SortOrder:   "asc",
		}

		// Count Query
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(userID, "%search%", "%search%", "tag", startDate, endDate).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Select Query
		// Expects correct ORDER BY clause constructed
		mock.ExpectQuery("SELECT .* FROM notes .* ORDER BY created_at ASC").
			WithArgs(userID, "%search%", "%search%", "tag", startDate, endDate).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow("nid", userID, "title", json.RawMessage("{}"), "active", now, now, nil))

		// Tags Query
		mock.ExpectQuery("SELECT .* FROM tags").
			WithArgs("nid").
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "note_id"}))

		var getErr error
		_, _, getErr = service.GetNotes(ctx, userID, 1, 10, filter)
		assert.NoError(t, getErr)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Scan Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		mock.ExpectQuery("SELECT COUNT").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery("SELECT .* FROM notes").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow("nid", userID, "title", json.RawMessage("{}"), "active", "invalid-date", now, nil)) // Invalid date

		var getErr error
		notes, _, getErr := service.GetNotes(ctx, userID, 1, 10, NoteFilter{})
		assert.NoError(t, getErr)
		assert.Empty(t, notes)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

// TestNoteService_GetUserTags_Errors
func TestNoteService_GetUserTags_Errors(t *testing.T) {
	ctx := context.Background()
	userID := "uid"

	t.Run("Query Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		mock.ExpectQuery("SELECT .* FROM tags").
			WithArgs(userID).
			WillReturnError(errors.New("query error"))

		_, err = service.GetUserTags(ctx, userID)
		assert.Error(t, err)
		assert.Equal(t, "query error", err.Error())
	})

	t.Run("Scan Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		mock.ExpectQuery("SELECT .* FROM tags").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
				AddRow("tid", "uid", "name", "invalid-date"))

		_, err = service.GetUserTags(ctx, userID)
		assert.Error(t, err)
		assert.NotNil(t, err)
	})
}
