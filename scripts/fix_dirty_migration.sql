-- Fix dirty state for version 1
-- This removes the "dirty" entry from the schema_migrations table, allowing the migration to be re-attempted.
DELETE FROM schema_migrations WHERE version = 1;

-- Attempt to create the uuid-ossp extension
-- This is often the cause of failure for the first migration if the user lacks permissions.
-- If this fails with "permission denied", you must run this script as a superuser (e.g., 'postgres').
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
