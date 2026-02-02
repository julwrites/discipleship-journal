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
	CreateNote(ctx context.Context, userID, title string, content json.RawMessage, tags []string, status ...string) (*Note, error)
	DeleteNote(ctx context.Context, userID, noteID string) error
	UpdateNote(ctx context.Context, userID, noteID, title string, content json.RawMessage, tags []string, status ...string) error
	GetNote(ctx context.Context, userID, noteID string) (*Note, error)
	GetNotes(ctx context.Context, userID string, page, limit int, filter NoteFilter) ([]Note, int, error)
	CreateTag(ctx context.Context, userID, name string) (*Tag, error)
	GetUserTags(ctx context.Context, userID string) ([]Tag, error)
	DeleteTag(ctx context.Context, userID, tagID string) error
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
	Status    string          `json:"status"` // "active", "pending", "processing", "failed"
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty"`
	Tags      []Tag           `json:"tags"`
}

type Tag struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *NoteService) CreateNote(ctx context.Context, userID, title string, content json.RawMessage, tags []string, status ...string) (*Note, error) {
	var note Note
	statusVal := "active"
	if len(status) > 0 {
		statusVal = status[0]
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
		INSERT INTO notes (user_id, title, content, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, content, status, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query, userID, title, content, statusVal).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.Status, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	note.Tags = []Tag{}
	for _, tagName := range tags {
		// Ensure tag exists
		var tag Tag
		tagQuery := `
			INSERT INTO tags (user_id, name) VALUES ($1, $2)
			ON CONFLICT (user_id, name) DO UPDATE SET name=EXCLUDED.name
			RETURNING id, user_id, name, created_at
		`
		err := tx.QueryRow(ctx, tagQuery, userID, tagName).Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Link tag to note
		linkQuery := `INSERT INTO note_tags (note_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err = tx.Exec(ctx, linkQuery, note.ID, tag.ID)
		if err != nil {
			return nil, err
		}
		note.Tags = append(note.Tags, tag)
	}

	if err := tx.Commit(ctx); err != nil {
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

func (s *NoteService) UpdateNote(ctx context.Context, userID, noteID, title string, content json.RawMessage, tags []string, status ...string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var query string
	var args []interface{}

	if len(status) > 0 {
		query = `
			UPDATE notes
			SET title=$1, content=$2, status=$3, updated_at=NOW()
			WHERE id=$4 AND user_id=$5 AND deleted_at IS NULL
		`
		args = []interface{}{title, content, status[0], noteID, userID}
	} else {
		query = `
			UPDATE notes
			SET title=$1, content=$2, updated_at=NOW()
			WHERE id=$3 AND user_id=$4 AND deleted_at IS NULL
		`
		args = []interface{}{title, content, noteID, userID}
	}

	commandTag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	// Update tags
	// First, remove existing tags that are not in the new list?
	// Or just wipe and recreate?
	// Wipe and recreate is easiest for now, but preserve IDs?
	// note_tags table only has (note_id, tag_id).
	// So deleting all from note_tags where note_id=$1 is fine.
	_, err = tx.Exec(ctx, "DELETE FROM note_tags WHERE note_id=$1", noteID)
	if err != nil {
		return err
	}

	for _, tagName := range tags {
		var tag Tag
		tagQuery := `
			INSERT INTO tags (user_id, name) VALUES ($1, $2)
			ON CONFLICT (user_id, name) DO UPDATE SET name=EXCLUDED.name
			RETURNING id, user_id, name, created_at
		`
		err := tx.QueryRow(ctx, tagQuery, userID, tagName).Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt)
		if err != nil {
			return err
		}

		linkQuery := `INSERT INTO note_tags (note_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err = tx.Exec(ctx, linkQuery, noteID, tag.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *NoteService) GetNote(ctx context.Context, userID, noteID string) (*Note, error) {
	var note Note
	query := "SELECT id, user_id, title, content, status, created_at, updated_at, deleted_at FROM notes WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL"
	err := s.db.QueryRow(ctx, query, noteID, userID).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.Status, &note.CreatedAt, &note.UpdatedAt, &note.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	// Fetch tags
	tagsQuery := `
		SELECT t.id, t.user_id, t.name, t.created_at
		FROM tags t
		JOIN note_tags nt ON t.id = nt.tag_id
		WHERE nt.note_id = $1
		ORDER BY t.name
	`
	rows, err := s.db.Query(ctx, tagsQuery, note.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	note.Tags = []Tag{}
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		note.Tags = append(note.Tags, t)
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
	baseQuery := "SELECT id, user_id, title, content, status, created_at, updated_at, deleted_at FROM notes" + whereClause
	baseQuery += fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d", orderBy, orderDir, limit, offset)

	rows, err := s.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		// Scan status as well
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.Status, &n.CreatedAt, &n.UpdatedAt, &n.DeletedAt); err != nil {
			continue
		}
		notes = append(notes, n)
	}

	if notes == nil {
		notes = []Note{}
	}

	// Fetch tags for all notes (optimize to avoid N+1)
	if len(notes) > 0 {
		tagsQuery := `
			SELECT t.id, t.user_id, t.name, t.created_at, nt.note_id
			FROM tags t
			JOIN note_tags nt ON t.id = nt.tag_id
			WHERE nt.note_id = ANY($1)
			ORDER BY t.name
		`
		// Extract IDs as string array
		ids := make([]string, len(notes))
		for i, n := range notes {
			ids[i] = n.ID
		}

		rowsTags, err := s.db.Query(ctx, tagsQuery, ids)
		if err != nil {
			return nil, 0, err
		}
		defer rowsTags.Close()

		tagsMap := make(map[string][]Tag)
		for rowsTags.Next() {
			var t Tag
			var noteID string
			if err := rowsTags.Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt, &noteID); err != nil {
				return nil, 0, err
			}
			tagsMap[noteID] = append(tagsMap[noteID], t)
		}

		for i := range notes {
			if tags, ok := tagsMap[notes[i].ID]; ok {
				notes[i].Tags = tags
			} else {
				notes[i].Tags = []Tag{}
			}
		}
	}

	return notes, total, nil
}

func (s *NoteService) CreateTag(ctx context.Context, userID, name string) (*Tag, error) {
	var tag Tag
	query := `
		INSERT INTO tags (user_id, name)
		VALUES ($1, $2)
		ON CONFLICT (user_id, name) DO UPDATE SET name=EXCLUDED.name
		RETURNING id, user_id, name, created_at
	`
	err := s.db.QueryRow(ctx, query, userID, name).Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *NoteService) GetUserTags(ctx context.Context, userID string) ([]Tag, error) {
	query := "SELECT id, user_id, name, created_at FROM tags WHERE user_id=$1 ORDER BY name"
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	if tags == nil {
		tags = []Tag{}
	}
	return tags, nil
}

func (s *NoteService) DeleteTag(ctx context.Context, userID, tagID string) error {
	query := "DELETE FROM tags WHERE id=$1 AND user_id=$2"
	commandTag, err := s.db.Exec(ctx, query, tagID, userID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
