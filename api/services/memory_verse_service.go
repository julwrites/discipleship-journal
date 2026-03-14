package services

import (
	"context"
	"encoding/json"
	"time"

	"database/sql"
	"discipleship_journal_api/models"

	"github.com/google/uuid"
)

type MemoryVerseService interface {
	GetPacks(ctx context.Context, userID uuid.UUID, typeFilter string) ([]*models.VersePack, error)
	GetPack(ctx context.Context, packID uuid.UUID, userID uuid.UUID) (*models.VersePack, error)
	CreatePack(ctx context.Context, pack *models.VersePack) (*models.VersePack, error)
	GetVerses(ctx context.Context, packID uuid.UUID, userID uuid.UUID) ([]*models.MemoryVerse, error)
	GetOriginalVerses(ctx context.Context, packID uuid.UUID) ([]*models.MemoryVerse, error)
	CreateVerse(ctx context.Context, verse *models.MemoryVerse) (*models.MemoryVerse, error)
	ClonePack(ctx context.Context, packID uuid.UUID, userID uuid.UUID, newTitle string, useUserDefault bool) (*models.VersePack, error)
	DeletePack(ctx context.Context, packID uuid.UUID, userID uuid.UUID) error
	UpdateVerse(ctx context.Context, verse *models.MemoryVerse, userID uuid.UUID) error
	DeleteVerse(ctx context.Context, verseID uuid.UUID, userID uuid.UUID) error
	// SearchVerses searches for verses across all accessible packs (user's or system's)
	SearchVerses(ctx context.Context, userID uuid.UUID, query string) ([]*models.MemoryVerse, error)
	SetVersePreference(ctx context.Context, userID, verseID uuid.UUID, version string) error
	SetVersePreferencesBatch(ctx context.Context, userID uuid.UUID, verseIDs []uuid.UUID, version string) error
	RemoveVersePreference(ctx context.Context, userID, verseID uuid.UUID) error
}

type memoryVerseService struct {
	db DBInterfaceWithQuery
}

type pgxRows = *sql.Rows

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
		rowsRows, err = s.db.QueryContext(ctx, baseQuery)
	} else {
		// User packs
		baseQuery += " WHERE vp.user_id = ? ORDER BY vp.created_at DESC"
		rowsRows, err = s.db.QueryContext(ctx, baseQuery, userID)
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
		WHERE vp.id = ? AND (vp.user_id = ? OR vp.is_public = true)
	`
	row := s.db.QueryRowContext(ctx, query, packID, userID)

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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.ExecContext(ctx, query,
		pack.ID, pack.UserID, pack.Title, pack.Identifier, pack.Description, pack.IsPublic, pack.CreatedAt, pack.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (s *memoryVerseService) GetVerses(ctx context.Context, packID uuid.UUID, userID uuid.UUID) ([]*models.MemoryVerse, error) {
	// Query to fetch verses along with user preferences and default settings
	// Priority: 1. User Override (user_verse_preferences)
	//           2. User Default (users.settings->>'bible_version')
	//           3. Original Verse Version (memory_verses.version)
	query := `
		SELECT
			mv.id,
			mv.verse_pack_id,
			mv.reference,
			mv.title,
			COALESCE(uvp.version_override, u.settings->>'$.bible_version', mv.version) as effective_version,
			CASE
				WHEN uvp.version_override IS NOT NULL THEN 'override'
				WHEN u.settings->>'$.bible_version' IS NOT NULL THEN 'user_default'
				ELSE 'original'
			END as version_source,
			mv.tags,
			mv.created_at,
			mv.updated_at
		FROM memory_verses mv
		LEFT JOIN user_verse_preferences uvp ON mv.id = uvp.verse_id AND uvp.user_id = ?
		LEFT JOIN users u ON u.id = ?
		WHERE mv.verse_pack_id = ?
		ORDER BY mv.created_at ASC
	`
	rows, err := s.db.QueryContext(ctx, query, userID, userID, packID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []*models.MemoryVerse
	for rows.Next() {
		var v models.MemoryVerse
		var tagsBytes []byte
		var title *string
		if err := rows.Scan(&v.ID, &v.VersePackID, &v.Reference, &title, &v.Version, &v.VersionSource, &tagsBytes, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		if title != nil {
			v.Title = *title
		}
		if len(tagsBytes) > 0 {
			_ = json.Unmarshal(tagsBytes, &v.Tags)
		}
		verses = append(verses, &v)
	}
	return verses, nil
}

func (s *memoryVerseService) GetOriginalVerses(ctx context.Context, packID uuid.UUID) ([]*models.MemoryVerse, error) {
	// Query to fetch verses exactly as they are in the database, ignoring user preferences
	query := `
		SELECT
			mv.id,
			mv.verse_pack_id,
			mv.reference,
			mv.title,
			mv.version,
			'original' as version_source,
			mv.tags,
			mv.created_at,
			mv.updated_at
		FROM memory_verses mv
		WHERE mv.verse_pack_id = ?
		ORDER BY mv.created_at ASC
	`
	rows, err := s.db.QueryContext(ctx, query, packID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []*models.MemoryVerse
	for rows.Next() {
		var v models.MemoryVerse
		var tagsBytes []byte
		var title *string
		if err := rows.Scan(&v.ID, &v.VersePackID, &v.Reference, &title, &v.Version, &v.VersionSource, &tagsBytes, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		if title != nil {
			v.Title = *title
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
		INSERT INTO memory_verses (id, verse_pack_id, reference, title, version, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.ExecContext(ctx, query,
		verse.ID, verse.VersePackID, verse.Reference, verse.Title, verse.Version, tagsJSON, verse.CreatedAt, verse.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return verse, nil
}

func (s *memoryVerseService) ClonePack(ctx context.Context, packID uuid.UUID, userID uuid.UUID, newTitle string, useUserDefault bool) (*models.VersePack, error) {
	// 1. Get original pack
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
	var verses []*models.MemoryVerse
	if useUserDefault {
		// Use GetVerses which resolves effective version based on user prefs/defaults
		verses, err = s.GetVerses(ctx, packID, userID)
	} else {
		// Use GetOriginalVerses which gets raw versions
		verses, err = s.GetOriginalVerses(ctx, packID)
	}

	if err != nil {
		return nil, err
	}

	for _, v := range verses {
		newVerse := &models.MemoryVerse{
			VersePackID: createdPack.ID,
			Reference:   v.Reference,
			Title:       v.Title,
			Version:     v.Version,
			Tags:        v.Tags,
		}
		if _, err := s.CreateVerse(ctx, newVerse); err != nil {
			return nil, err
		}
	}

	createdPack.VerseCount = len(verses)
	return createdPack, nil
}

func (s *memoryVerseService) DeletePack(ctx context.Context, packID uuid.UUID, userID uuid.UUID) error {
	// Check ownership
	query := `DELETE FROM verse_packs WHERE id = ? AND user_id = ?`
	res, err := s.db.ExecContext(ctx, query, packID, userID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *memoryVerseService) UpdateVerse(ctx context.Context, verse *models.MemoryVerse, userID uuid.UUID) error {
	tagsJSON, _ := json.Marshal(verse.Tags)

	query := `
		UPDATE memory_verses mv
		JOIN verse_packs vp ON mv.verse_pack_id = vp.id
		SET mv.reference = ?, mv.title = ?, mv.version = ?, mv.tags = ?, mv.updated_at = NOW()
		WHERE mv.id = ?
		AND vp.user_id = ?
	`
	res, err := s.db.ExecContext(ctx, query,
		verse.Reference, verse.Title, verse.Version, tagsJSON, verse.ID, userID,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *memoryVerseService) DeleteVerse(ctx context.Context, verseID uuid.UUID, userID uuid.UUID) error {
	query := `
		DELETE mv FROM memory_verses mv
		JOIN verse_packs vp ON mv.verse_pack_id = vp.id
		WHERE mv.id = ?
		AND vp.user_id = ?
	`
	res, err := s.db.ExecContext(ctx, query, verseID, userID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *memoryVerseService) SearchVerses(ctx context.Context, userID uuid.UUID, queryStr string) ([]*models.MemoryVerse, error) {
	// Search in user's packs OR public packs
	// We also apply version resolution here
	query := `
		SELECT
			mv.id,
			mv.verse_pack_id,
			mv.reference,
			mv.title,
			COALESCE(uvp.version_override, u.settings->>'$.bible_version', mv.version) as effective_version,
			CASE
				WHEN uvp.version_override IS NOT NULL THEN 'override'
				WHEN u.settings->>'$.bible_version' IS NOT NULL THEN 'user_default'
				ELSE 'original'
			END as version_source,
			mv.tags,
			vp.title,
			mv.created_at,
			mv.updated_at
		FROM memory_verses mv
		JOIN verse_packs vp ON mv.verse_pack_id = vp.id
		LEFT JOIN user_verse_preferences uvp ON mv.id = uvp.verse_id AND uvp.user_id = ?
		LEFT JOIN users u ON u.id = ?
		WHERE (vp.user_id = ? OR vp.is_public = true)
		AND (mv.reference LIKE ? OR vp.title LIKE ? OR mv.title LIKE ?)
		ORDER BY mv.reference ASC
		LIMIT 20
	`
	searchStr := "%" + queryStr + "%"
	rows, err := s.db.QueryContext(ctx, query, userID, userID, userID, searchStr, searchStr, searchStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []*models.MemoryVerse
	for rows.Next() {
		var v models.MemoryVerse
		var tagsBytes []byte
		var title *string
		if err := rows.Scan(&v.ID, &v.VersePackID, &v.Reference, &title, &v.Version, &v.VersionSource, &tagsBytes, &v.PackTitle, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		if title != nil {
			v.Title = *title
		}
		if len(tagsBytes) > 0 {
			_ = json.Unmarshal(tagsBytes, &v.Tags)
		}
		verses = append(verses, &v)
	}
	return verses, nil
}

func (s *memoryVerseService) SetVersePreference(ctx context.Context, userID, verseID uuid.UUID, version string) error {
	query := `
		INSERT INTO user_verse_preferences (user_id, verse_id, version_override)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE version_override = VALUES(version_override), updated_at = NOW()
	`
	_, err := s.db.ExecContext(ctx, query, userID, verseID, version)
	return err
}

func (s *memoryVerseService) SetVersePreferencesBatch(ctx context.Context, userID uuid.UUID, verseIDs []uuid.UUID, version string) error {
	if len(verseIDs) == 0 {
		return nil
	}

	query := "INSERT INTO user_verse_preferences (user_id, verse_id, version_override) VALUES "
	var args []interface{}

	for i, verseID := range verseIDs {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?)"
		args = append(args, userID, verseID, version)
	}

	query += " ON DUPLICATE KEY UPDATE version_override = VALUES(version_override), updated_at = NOW()"

	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *memoryVerseService) RemoveVersePreference(ctx context.Context, userID, verseID uuid.UUID) error {
	query := `DELETE FROM user_verse_preferences WHERE user_id = ? AND verse_id = ?`
	_, err := s.db.ExecContext(ctx, query, userID, verseID)
	return err
}
