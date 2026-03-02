ALTER TABLE notes RENAME TO journal_entries;
ALTER TABLE users DROP COLUMN username;
ALTER TABLE users DROP COLUMN settings;
