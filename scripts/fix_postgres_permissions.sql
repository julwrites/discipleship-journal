-- This script fixes permissions for the application user on a Cloud SQL PostgreSQL database.
-- Run this script as a SUPERUSER (e.g., 'postgres').
-- Usage: psql -h <HOST> -U postgres -d <DB_NAME> -v user=<APP_USER> -f scripts/fix_postgres_permissions.sql
-- Example: psql -h 127.0.0.1 -U postgres -d discipleship_journal -v user=dj_app_user -f scripts/fix_postgres_permissions.sql

\echo 'Granting permissions for user:' :user

-- 1. Ensure the UUID extension is enabled (requires superuser)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Grant usage on the public schema
GRANT USAGE ON SCHEMA public TO :user;

-- 3. Grant table creation privileges on the public schema
GRANT CREATE ON SCHEMA public TO :user;

-- 4. Grant privileges on all existing tables and sequences
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO :user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO :user;

-- 5. Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO :user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO :user;

\echo 'Permissions granted successfully.'
