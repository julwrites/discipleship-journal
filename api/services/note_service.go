package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"discipleship_journal_api/models"

	"database/sql"

	"github.com/google/uuid"
)

// NoteFilter defines the criteria for filtering notes.
type NoteFilter struct {
	SearchQuery string
	Tag         string
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
	QueryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, sql string, args ...any) *sql.Row
	ExecContext(ctx context.Context, sql string, arguments ...any) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
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
	if s.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var note Note
	statusVal := "active"
	if len(status) > 0 {
		statusVal = status[0]
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	note.ID = uuid.New().String()
	note.UserID = userID
	note.Title = title
	note.Content = content
	note.Status = statusVal

	query := `
		INSERT INTO notes (id, user_id, title, content, status)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, query, note.ID, userID, title, content, statusVal)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(ctx, "SELECT created_at, updated_at FROM notes WHERE id = ?", note.ID).Scan(&note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	note.Tags = []Tag{}
	for _, tagName := range tags {
		var tag Tag
		tag.ID = uuid.New().String()
		tag.UserID = userID
		tag.Name = tagName

		tagQuery := `
			INSERT INTO tags (id, user_id, name) VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE name=VALUES(name)
		`
		_, err := tx.ExecContext(ctx, tagQuery, tag.ID, userID, tagName)
		if err != nil {
			return nil, err
		}

		err = tx.QueryRowContext(ctx, "SELECT id, created_at FROM tags WHERE user_id=? AND name=?", userID, tagName).Scan(&tag.ID, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}

		linkQuery := `INSERT IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)`
		_, err = tx.ExecContext(ctx, linkQuery, note.ID, tag.ID)
		if err != nil {
			return nil, err
		}
		note.Tags = append(note.Tags, tag)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &note, nil
}

func (s *NoteService) DeleteNote(ctx context.Context, userID, noteID string) error {
	if s.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	query := `UPDATE notes SET deleted_at=NOW() WHERE id=? AND user_id=? AND deleted_at IS NULL`
	commandTag, err := s.db.ExecContext(ctx, query, noteID, userID)
	if err != nil {
		return err
	}
	rowsAffected, _ := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *NoteService) UpdateNote(ctx context.Context, userID, noteID, title string, content json.RawMessage, tags []string, status ...string) error {
	if s.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var query string
	var args []interface{}

	if len(status) > 0 {
		query = `
			UPDATE notes
			SET title=?, content=?, status=?, updated_at=NOW()
			WHERE id=? AND user_id=? AND deleted_at IS NULL
		`
		args = []interface{}{title, content, status[0], noteID, userID}
	} else {
		query = `
			UPDATE notes
			SET title=?, content=?, updated_at=NOW()
			WHERE id=? AND user_id=? AND deleted_at IS NULL
		`
		args = []interface{}{title, content, noteID, userID}
	}

	commandTag, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, _ := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	// Update tags
	// First, remove existing tags that are not in the new list?
	// Or just wipe and recreate?
	// Wipe and recreate is easiest for now, but preserve IDs?
	// note_tags table only has (note_id, tag_id).
	// So deleting all from note_tags where note_id=? is fine.
	_, err = tx.ExecContext(ctx, "DELETE FROM note_tags WHERE note_id=?", noteID)
	if err != nil {
		return err
	}

	for _, tagName := range tags {
		var tag Tag
		tag.ID = uuid.New().String()
		tagQuery := `
			INSERT INTO tags (id, user_id, name) VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE name=VALUES(name)
		`
		_, err := tx.ExecContext(ctx, tagQuery, tag.ID, userID, tagName)
		if err != nil {
			return err
		}

		err = tx.QueryRowContext(ctx, "SELECT id, created_at FROM tags WHERE user_id=? AND name=?", userID, tagName).Scan(&tag.ID, &tag.CreatedAt)
		if err != nil {
			return err
		}

		linkQuery := `INSERT IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)`
		_, err = tx.ExecContext(ctx, linkQuery, noteID, tag.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *NoteService) GetNote(ctx context.Context, userID, noteID string) (*Note, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	var note Note
	query := "SELECT id, user_id, title, content, status, created_at, updated_at, deleted_at FROM notes WHERE id=? AND user_id=? AND deleted_at IS NULL"
	err := s.db.QueryRowContext(ctx, query, noteID, userID).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.Status, &note.CreatedAt, &note.UpdatedAt, &note.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	// Fetch tags
	tagsQuery := `
		SELECT t.id, t.user_id, t.name, t.created_at
		FROM tags t
		JOIN note_tags nt ON t.id = nt.tag_id
		WHERE nt.note_id = ?
		ORDER BY t.name
	`
	rows, err := s.db.QueryContext(ctx, tagsQuery, note.ID)
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
	if s.db == nil {
		return nil, 0, fmt.Errorf("database connection is nil")
	}
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
	whereClause := " WHERE user_id=? AND deleted_at IS NULL"
	args := []interface{}{userID}

	if filter.SearchQuery != "" {
		whereClause += " AND (title LIKE ? OR content LIKE ?)"
		args = append(args, "%"+filter.SearchQuery+"%", "%"+filter.SearchQuery+"%")
	}

	if filter.Tag != "" {
		whereClause += " AND EXISTS (SELECT 1 FROM note_tags nt JOIN tags t ON nt.tag_id = t.id WHERE nt.note_id = notes.id AND t.name = ?)"
		args = append(args, filter.Tag)
	}

	if filter.StartDate != nil {
		whereClause += " AND updated_at >= ?"
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		whereClause += " AND updated_at <= ?"
		args = append(args, *filter.EndDate)
	}

	// Count total notes
	countQuery := "SELECT COUNT(*) FROM notes" + whereClause
	err = s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
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

	rows, err := s.db.QueryContext(ctx, baseQuery, args...)
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
		placeholders := make([]string, len(notes))
		args := make([]interface{}, len(notes))
		for i, n := range notes {
			placeholders[i] = "?"
			args[i] = n.ID
		}

		tagsQuery := fmt.Sprintf(`
			SELECT t.id, t.user_id, t.name, t.created_at, nt.note_id
			FROM tags t
			JOIN note_tags nt ON t.id = nt.tag_id
			WHERE nt.note_id IN (%s)
			ORDER BY t.name
		`, strings.Join(placeholders, ","))

		rowsTags, err := s.db.QueryContext(ctx, tagsQuery, args...)
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
	if s.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	var tag Tag
	tag.ID = uuid.New().String()
	tag.UserID = userID
	tag.Name = name

	query := `
		INSERT INTO tags (id, user_id, name)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE name=VALUES(name)
	`
	_, err := s.db.ExecContext(ctx, query, tag.ID, userID, name)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, "SELECT id, created_at FROM tags WHERE user_id=? AND name=?", userID, name).Scan(&tag.ID, &tag.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *NoteService) GetUserTags(ctx context.Context, userID string) ([]Tag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	query := "SELECT id, user_id, name, created_at FROM tags WHERE user_id=? ORDER BY name"
	rows, err := s.db.QueryContext(ctx, query, userID)
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
	if s.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	query := "DELETE FROM tags WHERE id=? AND user_id=?"
	commandTag, err := s.db.ExecContext(ctx, query, tagID, userID)
	if err != nil {
		return err
	}
	rowsAffected, _ := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}
