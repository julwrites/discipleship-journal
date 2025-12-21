package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"discipleship_journal_api/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// NoteServiceInterface defines the methods for note operations.
type NoteServiceInterface interface {
	CreateNote(ctx context.Context, userID, title string, content json.RawMessage) (*Note, error)
	DeleteNote(ctx context.Context, userID, noteID string) error
	UpdateNote(ctx context.Context, userID, noteID, title string, content json.RawMessage) error
	GetNote(ctx context.Context, userID, noteID string) (*Note, error)
	GetNotes(ctx context.Context, userID string, page, limit int, searchQuery string) ([]Note, int, error)
}

type NoteService struct {
	db DBInterface
}

type DBInterface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
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
		return models.ErrNotFound
	}
	return nil
}

func (s *NoteService) UpdateNote(ctx context.Context, userID, noteID, title string, content json.RawMessage) error {
	query := `
		UPDATE notes
		SET title=$1, content=$2, updated_at=NOW()
		WHERE id=$3 AND user_id=$4
	`
	commandTag, err := s.db.Exec(ctx, query, title, content, noteID, userID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *NoteService) GetNote(ctx context.Context, userID, noteID string) (*Note, error) {
	var note Note
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id=$1 AND user_id=$2"
	err := s.db.QueryRow(ctx, query, noteID, userID).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &note, nil
}

func (s *NoteService) GetNotes(ctx context.Context, userID string, page, limit int, searchQuery string) ([]Note, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int
	var err error

	// Count total notes
	countQuery := "SELECT COUNT(*) FROM notes WHERE user_id=$1"
	var countArgs []interface{}
	countArgs = append(countArgs, userID)

	if searchQuery != "" {
		countQuery += " AND (title ILIKE $2 OR content::text ILIKE $2)"
		countArgs = append(countArgs, "%"+searchQuery+"%")
	}

	err = s.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch notes
	baseQuery := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id=$1"
	var queryArgs []interface{}
	queryArgs = append(queryArgs, userID)

	if searchQuery != "" {
		baseQuery += " AND (title ILIKE $2 OR content::text ILIKE $2)"
		queryArgs = append(queryArgs, "%"+searchQuery+"%")
	}

	baseQuery += fmt.Sprintf(" ORDER BY updated_at DESC LIMIT %d OFFSET %d", limit, offset)

	rows, err := s.db.Query(ctx, baseQuery, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			continue
		}
		notes = append(notes, n)
	}

	if notes == nil {
		notes = []Note{}
	}

	return notes, total, nil
}
