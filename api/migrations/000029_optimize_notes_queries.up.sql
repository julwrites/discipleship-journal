CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_notes_user_deleted_updated
    ON notes (user_id, deleted_at, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_notes_user_deleted_title
    ON notes (user_id, deleted_at, title);

CREATE INDEX IF NOT EXISTS idx_notes_title_trgm
    ON notes USING gin (title gin_trgm_ops);
