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

// NoteFilter defines the criteria for filtering notes.
type NoteFilter struct {
	SearchQuery string
	StartDate   *time.Time
	EndDate     *time.Time
	SortBy      string // "updated_at", "created_at", "title"
	SortOrder   string // "asc", "desc"
}

// NoteServiceInterface defines the methods for note operations.
type NoteServiceInterface interface {
	CreateNote(ctx context.Context, userID, title string, content json.RawMessage) (*Note, error)
	DeleteNote(ctx context.Context, userID, noteID string) error
	UpdateNote(ctx context.Context, userID, noteID, title string, content json.RawMessage) error
	GetNote(ctx context.Context, userID, noteID string) (*Note, error)
	GetNotes(ctx context.Context, userID string, page, limit int, filter NoteFilter) ([]Note, int, error)
}

type NoteService struct {
	db DBInterface
}

type DBInterface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewNoteService(db DBInterface) *NoteService {
	return &NoteService{db: db}
}

type Note struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Title     string          `json:"title"`
	Content   json.RawMessage `json:"content,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty"`
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
	query := `UPDATE notes SET deleted_at=NOW() WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL`
	commandTag, err := s.db.Exec(ctx, query, noteID, userID)
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
		WHERE id=$3 AND user_id=$4 AND deleted_at IS NULL
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
	query := "SELECT id, user_id, title, content, created_at, updated_at, deleted_at FROM notes WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL"
	err := s.db.QueryRow(ctx, query, noteID, userID).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &note, nil
}

func (s *NoteService) GetNotes(ctx context.Context, userID string, page, limit int, filter NoteFilter) ([]Note, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int
	var err error

	// Base logic for building the WHERE clause
	// We'll use a slice of args and string building
	whereClause := " WHERE user_id=$1 AND deleted_at IS NULL"
	args := []interface{}{userID}
	argIdx := 2

	if filter.SearchQuery != "" {
		whereClause += fmt.Sprintf(" AND (title ILIKE $%d OR content::text ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+filter.SearchQuery+"%")
		argIdx++
	}

	if filter.StartDate != nil {
		whereClause += fmt.Sprintf(" AND updated_at >= $%d", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}

	if filter.EndDate != nil {
		whereClause += fmt.Sprintf(" AND updated_at <= $%d", argIdx)
		args = append(args, *filter.EndDate)
	}

	// Count total notes
	countQuery := "SELECT COUNT(*) FROM notes" + whereClause
	err = s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Determine sort
	orderBy := "updated_at"
	orderDir := "DESC"

	switch filter.SortBy {
	case "created_at":
		orderBy = "created_at"
	case "title":
		orderBy = "title"
	case "updated_at":
		orderBy = "updated_at"
	}

	if filter.SortOrder == "asc" {
		orderDir = "ASC"
	}

	// Fetch notes
	baseQuery := "SELECT id, user_id, title, created_at, updated_at, deleted_at FROM notes" + whereClause
	baseQuery += fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d", orderBy, orderDir, limit, offset)

	rows, err := s.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.CreatedAt, &n.UpdatedAt, &n.DeletedAt); err != nil {
			continue
		}
		notes = append(notes, n)
	}

	if notes == nil {
		notes = []Note{}
	}

	return notes, total, nil
}
