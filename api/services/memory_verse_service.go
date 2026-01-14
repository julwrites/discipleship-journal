package services

import (
	"context"
	"encoding/json"
	"time"

	"discipleship_journal_api/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MemoryVerseService interface {
	GetPacks(ctx context.Context, userID uuid.UUID, typeFilter string) ([]*models.VersePack, error)
	GetPack(ctx context.Context, packID uuid.UUID, userID uuid.UUID) (*models.VersePack, error)
	CreatePack(ctx context.Context, pack *models.VersePack) (*models.VersePack, error)
	GetVerses(ctx context.Context, packID uuid.UUID) ([]*models.MemoryVerse, error)
	CreateVerse(ctx context.Context, verse *models.MemoryVerse) (*models.MemoryVerse, error)
	ClonePack(ctx context.Context, packID uuid.UUID, userID uuid.UUID, newTitle string) (*models.VersePack, error)
	DeletePack(ctx context.Context, packID uuid.UUID, userID uuid.UUID) error
	UpdateVerse(ctx context.Context, verse *models.MemoryVerse, userID uuid.UUID) error
	DeleteVerse(ctx context.Context, verseID uuid.UUID, userID uuid.UUID) error
	// SearchVerses searches for verses across all accessible packs (user's or system's)
	SearchVerses(ctx context.Context, userID uuid.UUID, query string) ([]*models.MemoryVerse, error)
}

type memoryVerseService struct {
	db DBInterfaceWithQuery
}

type pgxRows = pgx.Rows

func NewMemoryVerseService(db DBInterfaceWithQuery) MemoryVerseService {
	return &memoryVerseService{db: db}
}

func (s *memoryVerseService) GetPacks(ctx context.Context, userID uuid.UUID, typeFilter string) ([]*models.VersePack, error) {
	// Base query with computed verse count
	baseQuery := `
		SELECT vp.id, vp.user_id, vp.title, vp.identifier, vp.description, vp.is_public, vp.created_at, vp.updated_at,
		       (SELECT count(*) FROM memory_verses mv WHERE mv.verse_pack_id = vp.id) as verse_count
		FROM verse_packs vp
	`

	var rowsRows pgxRows
	var err error

	if typeFilter == "system" {
		baseQuery += " WHERE vp.user_id IS NULL ORDER BY vp.title"
		rowsRows, err = s.db.Query(ctx, baseQuery)
	} else {
		// User packs
		baseQuery += " WHERE vp.user_id = $1 ORDER BY vp.created_at DESC"
		rowsRows, err = s.db.Query(ctx, baseQuery, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rowsRows.Close()

	var packs []*models.VersePack
	for rowsRows.Next() {
		var p models.VersePack
		var identifier, description *string
		if err := rowsRows.Scan(&p.ID, &p.UserID, &p.Title, &identifier, &description, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt, &p.VerseCount); err != nil {
			return nil, err
		}
		if identifier != nil {
			p.Identifier = *identifier
		}
		if description != nil {
			p.Description = *description
		}
		packs = append(packs, &p)
	}
	return packs, nil
}

func (s *memoryVerseService) GetPack(ctx context.Context, packID uuid.UUID, userID uuid.UUID) (*models.VersePack, error) {
	query := `
		SELECT vp.id, vp.user_id, vp.title, vp.identifier, vp.description, vp.is_public, vp.created_at, vp.updated_at,
		       (SELECT count(*) FROM memory_verses mv WHERE mv.verse_pack_id = vp.id) as verse_count
		FROM verse_packs vp
		WHERE vp.id = $1 AND (vp.user_id = $2 OR vp.is_public = true)
	`
	row := s.db.QueryRow(ctx, query, packID, userID)

	var p models.VersePack
	var identifier, description *string
	if err := row.Scan(&p.ID, &p.UserID, &p.Title, &identifier, &description, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt, &p.VerseCount); err != nil {
		return nil, err
	}
	if identifier != nil {
		p.Identifier = *identifier
	}
	if description != nil {
		p.Description = *description
	}
	return &p, nil
}

func (s *memoryVerseService) CreatePack(ctx context.Context, pack *models.VersePack) (*models.VersePack, error) {
	pack.ID = uuid.New()
	pack.CreatedAt = time.Now()
	pack.UpdatedAt = time.Now()

	query := `
		INSERT INTO verse_packs (id, user_id, title, identifier, description, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.db.Exec(ctx, query,
		pack.ID, pack.UserID, pack.Title, pack.Identifier, pack.Description, pack.IsPublic, pack.CreatedAt, pack.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (s *memoryVerseService) GetVerses(ctx context.Context, packID uuid.UUID) ([]*models.MemoryVerse, error) {
	query := `
		SELECT id, verse_pack_id, reference, version, tags, created_at, updated_at
		FROM memory_verses
		WHERE verse_pack_id = $1
		ORDER BY created_at ASC
	`
	rows, err := s.db.Query(ctx, query, packID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []*models.MemoryVerse
	for rows.Next() {
		var v models.MemoryVerse
		var tagsBytes []byte
		if err := rows.Scan(&v.ID, &v.VersePackID, &v.Reference, &v.Version, &tagsBytes, &v.CreatedAt, &v.UpdatedAt); err != nil {
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
		INSERT INTO memory_verses (id, verse_pack_id, reference, version, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.db.Exec(ctx, query,
		verse.ID, verse.VersePackID, verse.Reference, verse.Version, tagsJSON, verse.CreatedAt, verse.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return verse, nil
}

func (s *memoryVerseService) ClonePack(ctx context.Context, packID uuid.UUID, userID uuid.UUID, newTitle string) (*models.VersePack, error) {
	// 1. Get original pack
	// We can use GetPack but we need to fetch without ownership check if it's public (GetPack handles public)
	original, err := s.GetPack(ctx, packID, userID)
	if err != nil {
		return nil, err
	}

	// 2. Create new pack
	newPack := &models.VersePack{
		UserID:      &userID,
		Title:       newTitle,
		Identifier:  original.Identifier,
		Description: original.Description,
		IsPublic:    false,
	}
	if newPack.Title == "" {
		newPack.Title = original.Title // Keep original if not renamed
	}

	createdPack, err := s.CreatePack(ctx, newPack)
	if err != nil {
		return nil, err
	}

	// 3. Copy verses
	verses, err := s.GetVerses(ctx, packID)
	if err != nil {
		return nil, err
	}

	for _, v := range verses {
		newVerse := &models.MemoryVerse{
			VersePackID: createdPack.ID,
			Reference:   v.Reference,
			Version:     v.Version,
			Tags:        v.Tags,
		}
		if _, err := s.CreateVerse(ctx, newVerse); err != nil {
			return nil, err
		}
	}

	// Update verse count in response
	createdPack.VerseCount = len(verses)

	return createdPack, nil
}

func (s *memoryVerseService) DeletePack(ctx context.Context, packID uuid.UUID, userID uuid.UUID) error {
	// Check ownership
	query := `DELETE FROM verse_packs WHERE id = $1 AND user_id = $2`
	res, err := s.db.Exec(ctx, query, packID, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *memoryVerseService) UpdateVerse(ctx context.Context, verse *models.MemoryVerse, userID uuid.UUID) error {
	tagsJSON, _ := json.Marshal(verse.Tags)

	query := `
		UPDATE memory_verses mv
		SET reference = $2, version = $3, tags = $4, updated_at = NOW()
		FROM verse_packs vp
		WHERE mv.verse_pack_id = vp.id
		AND mv.id = $1
		AND vp.user_id = $5
	`
	res, err := s.db.Exec(ctx, query,
		verse.ID, verse.Reference, verse.Version, tagsJSON, userID,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *memoryVerseService) DeleteVerse(ctx context.Context, verseID uuid.UUID, userID uuid.UUID) error {
	query := `
		DELETE FROM memory_verses mv
		USING verse_packs vp
		WHERE mv.verse_pack_id = vp.id
		AND mv.id = $1
		AND vp.user_id = $2
	`
	res, err := s.db.Exec(ctx, query, verseID, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *memoryVerseService) SearchVerses(ctx context.Context, userID uuid.UUID, queryStr string) ([]*models.MemoryVerse, error) {
	// Search in user's packs OR public packs
	query := `
		SELECT mv.id, mv.verse_pack_id, mv.reference, mv.version, mv.tags, vp.title, mv.created_at, mv.updated_at
		FROM memory_verses mv
		JOIN verse_packs vp ON mv.verse_pack_id = vp.id
		WHERE (vp.user_id = $1 OR vp.is_public = true)
		AND (mv.reference ILIKE $2 OR vp.title ILIKE $2)
		ORDER BY mv.reference ASC
		LIMIT 20
	`
	rows, err := s.db.Query(ctx, query, userID, "%"+queryStr+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []*models.MemoryVerse
	for rows.Next() {
		var v models.MemoryVerse
		var tagsBytes []byte
		if err := rows.Scan(&v.ID, &v.VersePackID, &v.Reference, &v.Version, &tagsBytes, &v.PackTitle, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		if len(tagsBytes) > 0 {
			_ = json.Unmarshal(tagsBytes, &v.Tags)
		}
		verses = append(verses, &v)
	}
	return verses, nil
}
