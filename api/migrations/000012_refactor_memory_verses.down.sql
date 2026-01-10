-- Reverse group_shares changes
ALTER TABLE group_shares DROP CONSTRAINT chk_share_target;
ALTER TABLE group_shares DROP COLUMN verse_pack_id;
-- Cannot easily restore note_id NOT NULL if we inserted shares with only verse_pack_id, but assuming we want to rollback schema:
DELETE FROM group_shares WHERE note_id IS NULL;
ALTER TABLE group_shares ALTER COLUMN note_id SET NOT NULL;

-- Reverse memory_verses changes
DELETE FROM memory_verses;
ALTER TABLE memory_verses DROP COLUMN verse_pack_id;
ALTER TABLE memory_verses ADD COLUMN pack_name VARCHAR(100);
ALTER TABLE memory_verses ADD COLUMN text TEXT NOT NULL DEFAULT '';
ALTER TABLE memory_verses ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;

-- Drop verse_packs
DROP TABLE verse_packs;
