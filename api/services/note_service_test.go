package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateNote(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)

	ctx := context.Background()
	userID := "user-123"
	title := "Test Note"
	content := json.RawMessage(`{"text": "hello"}`)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO notes").
			WithArgs(userID, title, content).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "updated_at"}).
				AddRow("note-123", userID, title, content, now, now))

		note, err := service.CreateNote(ctx, userID, title, content)

		assert.NoError(t, err)
		assert.NotNil(t, note)
		assert.Equal(t, "note-123", note.ID)
		assert.Equal(t, userID, note.UserID)
		assert.Equal(t, title, note.Title)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO notes").
			WithArgs(userID, title, content).
			WillReturnError(errors.New("db error"))

		note, err := service.CreateNote(ctx, userID, title, content)

		assert.Error(t, err)
		assert.Nil(t, note)
		assert.Equal(t, "db error", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteNote(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)

	ctx := context.Background()
	userID := "user-123"
	noteID := "note-123"

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM notes").
			WithArgs(noteID, userID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := service.DeleteNote(ctx, userID, noteID)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM notes").
			WithArgs(noteID, userID).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := service.DeleteNote(ctx, userID, noteID)

		assert.Error(t, err)
		assert.Equal(t, pgx.ErrNoRows, err)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM notes").
			WithArgs(noteID, userID).
			WillReturnError(errors.New("db error"))

		err := service.DeleteNote(ctx, userID, noteID)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
