CREATE TABLE IF NOT EXISTS user_verse_preferences (
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    verse_id VARCHAR(36) NOT NULL REFERENCES memory_verses(id) ON DELETE CASCADE,
    version_override VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, verse_id)
);

CREATE INDEX idx_user_verse_preferences_user ON user_verse_preferences(user_id);
