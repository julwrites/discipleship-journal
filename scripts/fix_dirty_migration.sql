-- Fix for "Dirty database version X. Fix and force version."
-- Run this script as the database user (or a superuser if needed) on Cloud SQL.

-- 1. Connect to the database:
--    gcloud sql connect YOUR_INSTANCE_NAME --user=postgres --quiet

-- 2. Run the following command to check the current migration state:
-- SELECT * FROM schema_migrations;

-- 3. If you see a row with 'dirty = true', execute the following to clear the error state:
UPDATE schema_migrations SET dirty = false;

-- NOTE:
-- If the failed migration was partially applied (e.g., some tables created but others failed),
-- you may need to manually drop the created objects before re-running the migration.
-- For version 1 failures (initial schema), you can safely drop all tables and start over:
-- DELETE FROM schema_migrations;
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP TABLE IF EXISTS notes CASCADE;
-- ... (other tables from 000001_init_schema.up.sql)
