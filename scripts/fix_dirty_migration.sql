-- scripts/fix_dirty_migration.sql
-- Run this script as a SUPERUSER (e.g., 'postgres') to fix the 'Dirty database version 1' error and grant permissions.
--
-- IMPORTANT: Replace 'YOUR_DB_USER' with the actual database username used by your application.
-- For example: 'postgres', 'discipleship-journal', etc.

-- 1. Grant necessary permissions to the application user
GRANT ALL ON SCHEMA public TO "YOUR_DB_USER";
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "YOUR_DB_USER";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "YOUR_DB_USER";

-- 2. Fix the dirty migration state
DELETE FROM schema_migrations WHERE version = 1;

-- 3. Clean up potentially partial migration artifacts so migration 1 can run cleanly again
DROP TABLE IF EXISTS journal_entries;
DROP TABLE IF EXISTS users;

-- 4. Verify cleanup
SELECT * FROM schema_migrations;
