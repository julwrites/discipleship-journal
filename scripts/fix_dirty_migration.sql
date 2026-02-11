-- This script cleans up a "Dirty database version 1" state from a failed initial migration.
-- Run this script as the application user or a superuser (e.g., 'postgres').
-- Usage: psql -h <HOST> -U <USER> -d <DB_NAME> -f scripts/fix_dirty_migration.sql

\echo 'Cleaning up dirty migration state for version 1...'

-- 1. Remove the dirty migration entry from schema_migrations
DELETE FROM schema_migrations WHERE version = 1 AND dirty = true;

-- 2. Drop tables created by migration 1 if they exist (to allow a fresh start)
--    The initial migration creates 'users' and 'journal_entries'.
DROP TABLE IF EXISTS journal_entries CASCADE;
DROP TABLE IF EXISTS users CASCADE;

\echo 'Dirty migration state cleared. You can now re-run migrations.'
