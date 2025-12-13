package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"discipleship_journal_api/services"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateNote(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	service := services.NewNoteService(mock)
	userID := "user-123"
	title := "Test Note"
	content := json.RawMessage(`{"text": "content"}`)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO notes").
			WithArgs(userID, title, content).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "updated_at"}).
				AddRow("note-123", userID, title, content, time.Now(), time.Now()))

		note, err := service.CreateNote(context.Background(), userID, title, content)
		assert.NoError(t, err)
		assert.Equal(t, "note-123", note.ID)
		assert.Equal(t, userID, note.UserID)
		assert.Equal(t, title, note.Title)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO notes").
			WithArgs(userID, title, content).
			WillReturnError(errors.New("db error"))

		note, err := service.CreateNote(context.Background(), userID, title, content)
		assert.Error(t, err)
		assert.Nil(t, note)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
