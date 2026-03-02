CREATE TABLE IF NOT EXISTS group_shares (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    group_id VARCHAR(36) REFERENCES `groups`(id) ON DELETE CASCADE,
    note_id VARCHAR(36) REFERENCES notes(id) ON DELETE CASCADE,
    shared_by VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    shared_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    comment TEXT,
    UNIQUE(group_id, note_id)
);

CREATE INDEX idx_group_shares_group_id ON group_shares(group_id);
