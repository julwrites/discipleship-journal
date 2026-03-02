

-- CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_notes_user_deleted_updated
    ON notes (user_id, deleted_at, updated_at DESC);

CREATE INDEX idx_notes_user_deleted_title
    ON notes (user_id, deleted_at, title);

-- DROP INDEX IF EXISTS idx_notes_title_trgm;
