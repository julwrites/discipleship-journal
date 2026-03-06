ALTER TABLE users ADD COLUMN username VARCHAR(50);
ALTER TABLE users ADD UNIQUE INDEX idx_users_username (username);
ALTER TABLE users ADD COLUMN settings JSON;
ALTER TABLE journal_entries RENAME TO notes;
