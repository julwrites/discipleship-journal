ALTER TABLE groups ADD COLUMN IF NOT EXISTS type VARCHAR(20) NOT NULL DEFAULT 'group';
CREATE INDEX idx_groups_type ON groups(type);
