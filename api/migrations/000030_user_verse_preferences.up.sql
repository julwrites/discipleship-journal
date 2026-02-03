CREATE TABLE IF NOT EXISTS user_verse_preferences (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    verse_id UUID NOT NULL REFERENCES memory_verses(id) ON DELETE CASCADE,
    version_override VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, verse_id)
);

CREATE INDEX idx_user_verse_preferences_user ON user_verse_preferences(user_id);
