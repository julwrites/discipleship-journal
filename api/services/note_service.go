package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type NoteService struct {
	db DBInterface
}

type DBInterface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func NewNoteService(db DBInterface) *NoteService {
	return &NoteService{db: db}
}

type Note struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Title     string          `json:"title"`
	Content   json.RawMessage `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (s *NoteService) CreateNote(ctx context.Context, userID, title string, content json.RawMessage) (*Note, error) {
	var note Note
	query := `
		INSERT INTO notes (user_id, title, content)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, title, content, created_at, updated_at
	`
	err := s.db.QueryRow(ctx, query, userID, title, content).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (s *NoteService) DeleteNote(ctx context.Context, userID, noteID string) error {
	commandTag, err := s.db.Exec(ctx, "DELETE FROM notes WHERE id=$1 AND user_id=$2", noteID, userID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
