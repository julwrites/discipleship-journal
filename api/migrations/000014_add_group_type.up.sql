ALTER TABLE `groups` ADD COLUMN type VARCHAR(20) NOT NULL DEFAULT 'group';
CREATE INDEX idx_groups_type ON `groups`(type);
