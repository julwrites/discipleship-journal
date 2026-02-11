-- scripts/fix_dirty_migration.sql
-- Run this script as a SUPERUSER (e.g., 'postgres') to fix the 'Dirty database version 1' error and grant permissions.

-- IMPORTANT: The wrapper script (scripts/fix_db_permissions.sh) will replace YOUR_DB_USER with the actual username.

DO $$
BEGIN
    -- 1. Grant necessary permissions to the application user
    EXECUTE 'GRANT ALL ON SCHEMA public TO "YOUR_DB_USER"';
    EXECUTE 'GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "YOUR_DB_USER"';
    EXECUTE 'GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "YOUR_DB_USER"';

    -- 2. Ensure schema_migrations table is owned by the application user
    -- This prevents future "permission denied" errors on migrations.
    IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'schema_migrations') THEN
        EXECUTE 'ALTER TABLE schema_migrations OWNER TO "YOUR_DB_USER"';

        -- 3. Fix the dirty migration state
        DELETE FROM schema_migrations WHERE version = 1;
    END IF;

    -- 4. Clean up potentially partial migration artifacts so migration 1 can run cleanly again
    -- We drop these tables to allow a fresh start.
    DROP TABLE IF EXISTS journal_entries;
    DROP TABLE IF EXISTS users;

END $$;

-- 5. Verify cleanup (will print output if run interactively)
SELECT * FROM schema_migrations;
