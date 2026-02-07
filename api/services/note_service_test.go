package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"discipleship_journal_api/models"

	"github.com/google/uuid"
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
	userID := uuid.New().String()
	noteID := uuid.New().String()
	tagID := uuid.New().String()
	title := "Test Note"
	content := json.RawMessage(`{"text": "hello"}`)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO notes").
			WithArgs(userID, title, content, "active").
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at"}).
				AddRow(noteID, userID, title, content, "active", now, now))
		mock.ExpectCommit()

		note, err := service.CreateNote(ctx, userID, title, content, nil)

		assert.NoError(t, err)
		assert.NotNil(t, note)
		assert.Equal(t, noteID, note.ID)
		assert.Equal(t, userID, note.UserID)
		assert.Equal(t, title, note.Title)
		assert.Equal(t, "active", note.Status)
	})

	t.Run("success with tags", func(t *testing.T) {
		tags := []string{"tag1"}
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO notes").
			WithArgs(userID, title, content, "active").
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at"}).
				AddRow(noteID, userID, title, content, "active", now, now))

		// Tag creation/lookup
		mock.ExpectQuery("INSERT INTO tags").
			WithArgs(userID, "tag1").
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
				AddRow(tagID, userID, "tag1", now))

		// Link creation
		mock.ExpectExec("INSERT INTO note_tags").
			WithArgs(noteID, tagID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mock.ExpectCommit()

		note, err := service.CreateNote(ctx, userID, title, content, tags)

		assert.NoError(t, err)
		assert.NotNil(t, note)
		assert.Len(t, note.Tags, 1)
		assert.Equal(t, "tag1", note.Tags[0].Name)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO notes").
			WithArgs(userID, title, content, "active").
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		note, err := service.CreateNote(ctx, userID, title, content, nil)

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
	userID := uuid.New().String()
	noteID := uuid.New().String()

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
	userID := uuid.New().String()
	noteID := uuid.New().String()
	title := "Updated Title"
	content := json.RawMessage(`{"text": "updated"}`)

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE notes").
			WithArgs(title, content, noteID, userID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		mock.ExpectExec("DELETE FROM note_tags").
			WithArgs(noteID).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		mock.ExpectCommit()

		err := service.UpdateNote(ctx, userID, noteID, title, content, nil)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE notes").
			WithArgs(title, content, noteID, userID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))
		// Rollback is deferred, but CreateNote returns early.
		// Wait, in my implementation UpdateNote returns error early if rows affected == 0.
		// And deferred Rollback will be called.
		mock.ExpectRollback()

		err := service.UpdateNote(ctx, userID, noteID, title, content, nil)

		assert.Error(t, err)
		assert.Equal(t, models.ErrNotFound, err)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE notes").
			WithArgs(title, content, noteID, userID).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := service.UpdateNote(ctx, userID, noteID, title, content, nil)

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
	userID := uuid.New().String()
	noteID := uuid.New().String()
	title := "Test Note"
	content := json.RawMessage(`{"text": "hello"}`)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, title, content, status, created_at, updated_at, deleted_at FROM notes").
			WithArgs(noteID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow(noteID, userID, title, content, "active", now, now, nil))

		mock.ExpectQuery("SELECT .* FROM tags .* JOIN note_tags").
			WithArgs(noteID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "created_at"}))

		note, err := service.GetNote(ctx, userID, noteID)

		assert.NoError(t, err)
		assert.NotNil(t, note)
		assert.Equal(t, noteID, note.ID)
		assert.Equal(t, userID, note.UserID)
		assert.Equal(t, "active", note.Status)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, title, content, status, created_at, updated_at, deleted_at FROM notes").
			WithArgs(noteID, userID).
			WillReturnError(pgx.ErrNoRows)

		note, err := service.GetNote(ctx, userID, noteID)

		assert.Error(t, err)
		assert.Nil(t, note)
		assert.Equal(t, models.ErrNotFound, err)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, title, content, status, created_at, updated_at, deleted_at FROM notes").
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
	userID := uuid.New().String()
	noteID := uuid.New().String()
	title := "Test Note"
	content := json.RawMessage(`{"text": "hello"}`)
	now := time.Now()
	page := 1
	limit := 10

	t.Run("success no search", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery("SELECT id, user_id, title, content, status, created_at, updated_at, deleted_at FROM notes").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow(noteID, userID, title, content, "active", now, now, nil))

		// Tags query
		mock.ExpectQuery("SELECT .* FROM tags .* JOIN note_tags").
			WithArgs([]string{noteID}).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "created_at", "note_id"}))

		notes, total, err := service.GetNotes(ctx, userID, page, limit, NoteFilter{})

		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, 1, total)
		assert.Equal(t, noteID, notes[0].ID)
		assert.Equal(t, "active", notes[0].Status)
	})

	t.Run("success with search", func(t *testing.T) {
		searchQuery := "test"
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes").
			WithArgs(userID, "%"+searchQuery+"%").
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery("SELECT id, user_id, title, content, status, created_at, updated_at, deleted_at FROM notes").
			WithArgs(userID, "%"+searchQuery+"%").
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow(noteID, userID, title, content, "active", now, now, nil))

		mock.ExpectQuery("SELECT .* FROM tags .* JOIN note_tags").
			WithArgs([]string{noteID}).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "created_at", "note_id"}))

		notes, total, err := service.GetNotes(ctx, userID, page, limit, NoteFilter{SearchQuery: searchQuery})

		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, 1, total)
	})

	t.Run("success with tag filter", func(t *testing.T) {
		tag := "tag1"
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes").
			WithArgs(userID, tag).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery("SELECT id, user_id, title, content, status, created_at, updated_at, deleted_at FROM notes").
			WithArgs(userID, tag).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
				AddRow(noteID, userID, title, content, "active", now, now, nil))

		mock.ExpectQuery("SELECT .* FROM tags .* JOIN note_tags").
			WithArgs([]string{noteID}).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "created_at", "note_id"}))

		notes, total, err := service.GetNotes(ctx, userID, page, limit, NoteFilter{Tag: tag})

		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, 1, total)
	})

	t.Run("count error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes").
			WithArgs(userID).
			WillReturnError(errors.New("count error"))

		notes, total, err := service.GetNotes(ctx, userID, page, limit, NoteFilter{})

		assert.Error(t, err)
		assert.Nil(t, notes)
		assert.Equal(t, 0, total)
		assert.Equal(t, "count error", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateTag(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	tagID := uuid.New().String()
	tagName := "tag1"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO tags").
			WithArgs(userID, tagName).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
				AddRow(tagID, userID, tagName, now))

		tag, err := service.CreateTag(ctx, userID, tagName)
		assert.NoError(t, err)
		assert.Equal(t, tagName, tag.Name)
	})
}

func TestGetUserTags(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	tagID := uuid.New().String()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, name, created_at FROM tags").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
				AddRow(tagID, userID, "tag1", now))

		tags, err := service.GetUserTags(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, tags, 1)
	})
}

func TestDeleteTag(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	tagID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM tags").
			WithArgs(tagID, userID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := service.DeleteTag(ctx, userID, tagID)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM tags").
			WithArgs(tagID, userID).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := service.DeleteTag(ctx, userID, tagID)
		assert.Error(t, err)
		assert.Equal(t, models.ErrNotFound, err)
	})
}

func TestCreateNote_TagError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	noteID := uuid.New().String()
	title := "Test Note"
	content := json.RawMessage(`{"text": "hello"}`)
	tags := []string{"tag1"}
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO notes").
		WithArgs(userID, title, content, "active").
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at"}).
			AddRow(noteID, userID, title, content, "active", now, now))

	// Fail on tag insert
	mock.ExpectQuery("INSERT INTO tags").
		WithArgs(userID, "tag1").
		WillReturnError(errors.New("tag error"))

	mock.ExpectRollback()

	note, err := service.CreateNote(ctx, userID, title, content, tags)
	assert.Error(t, err)
	assert.Nil(t, note)
	assert.Equal(t, "tag error", err.Error())
}

func TestGetNotes_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery("SELECT id, user_id").
		WithArgs(userID).
		WillReturnError(errors.New("query error"))

	notes, total, err := service.GetNotes(ctx, userID, 1, 10, NoteFilter{})
	assert.Error(t, err)
	assert.Nil(t, notes)
	assert.Equal(t, 0, total)
	assert.Equal(t, "query error", err.Error())
}

func TestGetNotes_TagsQueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	noteID := uuid.New().String()
	now := time.Now()

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery("SELECT id, user_id").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(noteID, userID, "Title", json.RawMessage("{}"), "active", now, now, nil))

	mock.ExpectQuery("SELECT .* FROM tags").
		WithArgs([]string{noteID}).
		WillReturnError(errors.New("tags error"))

	notes, total, err := service.GetNotes(ctx, userID, 1, 10, NoteFilter{})
	assert.Error(t, err)
	assert.Nil(t, notes)
	assert.Equal(t, 0, total)
	assert.Equal(t, "tags error", err.Error())
}

func TestUpdateNote_WithTags(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	noteID := uuid.New().String()
	tagID := uuid.New().String()
	title := "Updated Title"
	content := json.RawMessage(`{"text": "updated"}`)
	tags := []string{"tag1"}
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE notes").
		WithArgs(title, content, noteID, userID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Delete existing tags
	mock.ExpectExec("DELETE FROM note_tags").
		WithArgs(noteID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	// Create/Get tag
	mock.ExpectQuery("INSERT INTO tags").
		WithArgs(userID, "tag1").
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
			AddRow(tagID, userID, "tag1", now))

	// Link tag
	mock.ExpectExec("INSERT INTO note_tags").
		WithArgs(noteID, tagID).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectCommit()

	err = service.UpdateNote(ctx, userID, noteID, title, content, tags)
	assert.NoError(t, err)
}

func TestUpdateNote_WithStatus(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	noteID := uuid.New().String()
	title := "Updated Title"
	content := json.RawMessage(`{"text": "updated"}`)
	status := "archived"

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE notes").
		WithArgs(title, content, status, noteID, userID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Delete existing tags (empty tags list passed)
	mock.ExpectExec("DELETE FROM note_tags").
		WithArgs(noteID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	mock.ExpectCommit()

	err = service.UpdateNote(ctx, userID, noteID, title, content, nil, status)
	assert.NoError(t, err)
}

func TestUpdateNote_TagError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	noteID := uuid.New().String()
	title := "Updated Title"
	content := json.RawMessage(`{"text": "updated"}`)
	tags := []string{"tag1"}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE notes").
		WithArgs(title, content, noteID, userID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectExec("DELETE FROM note_tags").
		WithArgs(noteID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	// Tag error
	mock.ExpectQuery("INSERT INTO tags").
		WithArgs(userID, "tag1").
		WillReturnError(errors.New("tag error"))

	// Rollback is deferred but implicit on error return?
	// NoteService returns error immediately.
	// We expect rollback.
	// The implementation has `defer func() { _ = tx.Rollback(ctx) }()`
	mock.ExpectRollback()

	err = service.UpdateNote(ctx, userID, noteID, title, content, tags)
	assert.Error(t, err)
	assert.Equal(t, "tag error", err.Error())
}

func TestUpdateNote_LinkTagError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	noteID := uuid.New().String()
	tagID := uuid.New().String()
	title := "Updated Title"
	content := json.RawMessage(`{"text": "updated"}`)
	tags := []string{"tag1"}
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE notes").
		WithArgs(title, content, noteID, userID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectExec("DELETE FROM note_tags").
		WithArgs(noteID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	mock.ExpectQuery("INSERT INTO tags").
		WithArgs(userID, "tag1").
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
			AddRow(tagID, userID, "tag1", now))

	// Link error
	mock.ExpectExec("INSERT INTO note_tags").
		WithArgs(noteID, tagID).
		WillReturnError(errors.New("link error"))

	mock.ExpectRollback()

	err = service.UpdateNote(ctx, userID, noteID, title, content, tags)
	assert.Error(t, err)
	assert.Equal(t, "link error", err.Error())
}

func TestGetNotes_Filters(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	// Count Query - Expects StartDate and EndDate arguments
	// Args: userID($1), start($2), end($3)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes.*updated_at >= \\$2 AND updated_at <= \\$3").
		WithArgs(userID, start, end).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))

	// Data Query - Expects arguments
	mock.ExpectQuery("SELECT .* FROM notes.*updated_at >= \\$2 AND updated_at <= \\$3").
		WithArgs(userID, start, end).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}))

	notes, total, err := service.GetNotes(ctx, userID, 1, 10, NoteFilter{StartDate: &start, EndDate: &end})
	assert.NoError(t, err)
	assert.Len(t, notes, 0)
	assert.Equal(t, 0, total)
}

func TestGetNotes_Sorting(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()

	// 1. Sort by title ASC
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery("SELECT .* FROM notes.*ORDER BY title ASC").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}))

	_, _, err = service.GetNotes(ctx, userID, 1, 10, NoteFilter{SortBy: "title", SortOrder: "asc"})
	assert.NoError(t, err)

	// 2. Sort by created_at DESC (default order, explicit sort by)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM notes").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery("SELECT .* FROM notes.*ORDER BY created_at DESC").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}))

	_, _, err = service.GetNotes(ctx, userID, 1, 10, NoteFilter{SortBy: "created_at", SortOrder: "desc"})
	assert.NoError(t, err)
}

func TestCreateNote_NilDB(t *testing.T) {
	service := NewNoteService(nil)
	ctx := context.Background()

	_, err := service.CreateNote(ctx, "user", "title", nil, nil)
	assert.Error(t, err)
	assert.Equal(t, "database connection is nil", err.Error())
}

func TestDeleteNote_NilDB(t *testing.T) {
	service := NewNoteService(nil)
	ctx := context.Background()

	err := service.DeleteNote(ctx, "user", "note")
	assert.Error(t, err)
	assert.Equal(t, "database connection is nil", err.Error())
}

func TestUpdateNote_NilDB(t *testing.T) {
	service := NewNoteService(nil)
	ctx := context.Background()

	err := service.UpdateNote(ctx, "user", "note", "title", nil, nil)
	assert.Error(t, err)
	assert.Equal(t, "database connection is nil", err.Error())
}

func TestGetNote_NilDB(t *testing.T) {
	service := NewNoteService(nil)
	ctx := context.Background()

	_, err := service.GetNote(ctx, "user", "note")
	assert.Error(t, err)
	assert.Equal(t, "database connection is nil", err.Error())
}

func TestGetNote_TagsError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := NewNoteService(mock)
	ctx := context.Background()
	userID := uuid.New().String()
	noteID := uuid.New().String()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, title, content, status, created_at, updated_at, deleted_at FROM notes").
		WithArgs(noteID, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(noteID, userID, "Title", json.RawMessage("{}"), "active", now, now, nil))

	mock.ExpectQuery("SELECT .* FROM tags .* JOIN note_tags").
		WithArgs(noteID).
		WillReturnError(errors.New("tags error"))

	note, err := service.GetNote(ctx, userID, noteID)
	assert.Error(t, err)
	assert.Nil(t, note)
	assert.Equal(t, "tags error", err.Error())
}

func TestGetNotes_NilDB(t *testing.T) {
	service := NewNoteService(nil)
	ctx := context.Background()

	_, _, err := service.GetNotes(ctx, "user", 1, 10, NoteFilter{})
	assert.Error(t, err)
	assert.Equal(t, "database connection is nil", err.Error())
}

func TestCreateTag_NilDB(t *testing.T) {
	service := NewNoteService(nil)
	ctx := context.Background()

	_, err := service.CreateTag(ctx, "user", "tag")
	assert.Error(t, err)
	assert.Equal(t, "database connection is nil", err.Error())
}

func TestGetUserTags_NilDB(t *testing.T) {
	service := NewNoteService(nil)
	ctx := context.Background()

	_, err := service.GetUserTags(ctx, "user")
	assert.Error(t, err)
	assert.Equal(t, "database connection is nil", err.Error())
}

func TestDeleteTag_NilDB(t *testing.T) {
	service := NewNoteService(nil)
	ctx := context.Background()

	err := service.DeleteTag(ctx, "user", "tag")
	assert.Error(t, err)
	assert.Equal(t, "database connection is nil", err.Error())
}
