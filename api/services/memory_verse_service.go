package services

import (
	"context"
	"encoding/json"
	"time"

	"discipleship_journal_api/models"
	"github.com/google/uuid"
)

type MemoryVerseService interface {
	SearchVerses(ctx context.Context, userID uuid.UUID, query string, tags []string) ([]*models.MemoryVerse, error)
	CreateVerse(ctx context.Context, verse *models.MemoryVerse) (*models.MemoryVerse, error)
	GetSystemPacks(ctx context.Context) ([]string, error)
}

type memoryVerseService struct {
	db DBInterfaceWithQuery
}

func NewMemoryVerseService(db DBInterfaceWithQuery) MemoryVerseService {
	return &memoryVerseService{db: db}
}

func (s *memoryVerseService) SearchVerses(ctx context.Context, userID uuid.UUID, query string, tags []string) ([]*models.MemoryVerse, error) {
	// Base query
	baseQuery := `
		SELECT id, user_id, pack_name, reference, text, version, tags, created_at, updated_at
		FROM memory_verses
		WHERE (user_id IS NULL OR user_id = $1)
	`
	args := []interface{}{userID}
	argIdx := 2

	if query != "" {
		// PostgreSQL standard is $2, $3 etc.
		// Since we append args, the index logic needs to match
		baseQuery += " AND (reference ILIKE $2 OR text ILIKE $2 OR pack_name ILIKE $2)"
		args = append(args, "%"+query+"%")
		argIdx++
	}

	baseQuery += " ORDER BY pack_name ASC, reference ASC LIMIT 50"

	rows, err := s.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []*models.MemoryVerse
	for rows.Next() {
		var v models.MemoryVerse
		var tagsBytes []byte
		if err := rows.Scan(&v.ID, &v.UserID, &v.PackName, &v.Reference, &v.Text, &v.Version, &tagsBytes, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		if len(tagsBytes) > 0 {
			_ = json.Unmarshal(tagsBytes, &v.Tags)
		}
		verses = append(verses, &v)
	}
	return verses, nil
}

func (s *memoryVerseService) CreateVerse(ctx context.Context, verse *models.MemoryVerse) (*models.MemoryVerse, error) {
	verse.ID = uuid.New()
	verse.CreatedAt = time.Now()
	verse.UpdatedAt = time.Now()

	tagsJSON, _ := json.Marshal(verse.Tags)

	query := `
		INSERT INTO memory_verses (id, user_id, pack_name, reference, text, version, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	_, err := s.db.Exec(ctx, query,
		verse.ID, verse.UserID, verse.PackName, verse.Reference, verse.Text, verse.Version, tagsJSON, verse.CreatedAt, verse.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return verse, nil
}

func (s *memoryVerseService) GetSystemPacks(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT pack_name FROM memory_verses WHERE user_id IS NULL ORDER BY pack_name`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packs []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		packs = append(packs, p)
	}
	return packs, nil
}
