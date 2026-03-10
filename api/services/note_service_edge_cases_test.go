package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNoteService_GetNotes_TagErrors(t *testing.T) {
	ctx := context.Background()
	userID := "uid"
	now := time.Now()

	t.Run("Tag_Query_Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		// Count Query
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Notes Query
		mock.ExpectQuery("SELECT .* FROM notes").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow("nid", userID, "title", json.RawMessage("{}"), "active", now, now, nil))

		// Tags Query -> Error
		mock.ExpectQuery("SELECT .* FROM tags").
			WithArgs("nid").
			WillReturnError(errors.New("tag query error"))

		_, _, err = service.GetNotes(ctx, userID, 1, 10, NoteFilter{})
		assert.Error(t, err)
		assert.Equal(t, "tag query error", err.Error())
	})

	t.Run("Tag_Scan_Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		_ = mock
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		defer db.Close()
		service := NewNoteService(db)

		// Count Query
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Notes Query
		mock.ExpectQuery("SELECT .* FROM notes").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow("nid", userID, "title", json.RawMessage("{}"), "active", now, now, nil))

		// Tags Query -> Rows with bad data
		mock.ExpectQuery("SELECT .* FROM tags").
			WithArgs("nid").
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "note_id"}).
				AddRow("tid", userID, "name", "invalid-time", "nid"))

		_, _, err = service.GetNotes(ctx, userID, 1, 10, NoteFilter{})
		assert.Error(t, err)
		// Scan error message depends on implementation/driver
	})
}
