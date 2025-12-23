package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"discipleship_journal_api/models"

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
		mock.ExpectExec("UPDATE notes SET deleted_at=.*").
			WithArgs(noteID, userID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := service.DeleteNote(ctx, userID, noteID)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE notes SET deleted_at=.*").
			WithArgs(noteID, userID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := service.DeleteNote(ctx, userID, noteID)

		assert.Error(t, err)
		assert.Equal(t, models.ErrNotFound, err)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("UPDATE notes SET deleted_at=.*").
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

func TestUpdateNote(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)

	ctx := context.Background()
	userID := "user-123"
	noteID := "note-123"
	title := "Updated Title"
	content := json.RawMessage(`{"text": "updated"}`)

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE notes").
			WithArgs(title, content, noteID, userID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := service.UpdateNote(ctx, userID, noteID, title, content)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE notes").
			WithArgs(title, content, noteID, userID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := service.UpdateNote(ctx, userID, noteID, title, content)

		assert.Error(t, err)
		assert.Equal(t, models.ErrNotFound, err)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec("UPDATE notes").
			WithArgs(title, content, noteID, userID).
			WillReturnError(errors.New("db error"))

		err := service.UpdateNote(ctx, userID, noteID, title, content)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetNote(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)

	ctx := context.Background()
	userID := "user-123"
	noteID := "note-123"
	title := "Test Note"
	content := json.RawMessage(`{"text": "hello"}`)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, title, content, created_at, updated_at, deleted_at FROM notes").
			WithArgs(noteID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "updated_at", "deleted_at"}).
				AddRow(noteID, userID, title, content, now, now, nil))

		note, err := service.GetNote(ctx, userID, noteID)

		assert.NoError(t, err)
		assert.NotNil(t, note)
		assert.Equal(t, noteID, note.ID)
		assert.Equal(t, userID, note.UserID)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, title, content, created_at, updated_at, deleted_at FROM notes").
			WithArgs(noteID, userID).
			WillReturnError(pgx.ErrNoRows)

		note, err := service.GetNote(ctx, userID, noteID)

		assert.Error(t, err)
		assert.Nil(t, note)
		assert.Equal(t, models.ErrNotFound, err)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, title, content, created_at, updated_at, deleted_at FROM notes").
			WithArgs(noteID, userID).
			WillReturnError(errors.New("db error"))

		note, err := service.GetNote(ctx, userID, noteID)

		assert.Error(t, err)
		assert.Nil(t, note)
		assert.Equal(t, "db error", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetNotes(t *testing.T) {
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
	page := 1
	limit := 10

	t.Run("success no search", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery("SELECT id, user_id, title, content, created_at, updated_at, deleted_at FROM notes").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "updated_at", "deleted_at"}).
				AddRow("note-123", userID, title, content, now, now, nil))

		notes, total, err := service.GetNotes(ctx, userID, page, limit, "")

		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, 1, total)
		assert.Equal(t, "note-123", notes[0].ID)
	})

	t.Run("success with search", func(t *testing.T) {
		searchQuery := "test"
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes").
			WithArgs(userID, "%"+searchQuery+"%").
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery("SELECT id, user_id, title, content, created_at, updated_at, deleted_at FROM notes").
			WithArgs(userID, "%"+searchQuery+"%").
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "updated_at", "deleted_at"}).
				AddRow("note-123", userID, title, content, now, now, nil))

		notes, total, err := service.GetNotes(ctx, userID, page, limit, searchQuery)

		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, 1, total)
	})

	t.Run("count error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes").
			WithArgs(userID).
			WillReturnError(errors.New("count error"))

		notes, total, err := service.GetNotes(ctx, userID, page, limit, "")

		assert.Error(t, err)
		assert.Nil(t, notes)
		assert.Equal(t, 0, total)
		assert.Equal(t, "count error", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
