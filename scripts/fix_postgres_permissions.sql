-- Fix for PostgreSQL 15+ "permission denied for schema public"
-- Run this script as the 'postgres' superuser or the instance owner on Cloud SQL.

-- 1. Connect to the database using the 'postgres' user:
--    gcloud sql connect YOUR_INSTANCE_NAME --user=postgres --quiet

-- 2. Once connected, run the following commands replacing 'YOUR_DB_USER' with the application user (e.g., 'discipleship_journal_user'):

-- Grant usage on the public schema (allows access to objects in it)
GRANT USAGE ON SCHEMA public TO "YOUR_DB_USER";

-- Grant create permission on the public schema (allows creating tables/migrations)
GRANT CREATE ON SCHEMA public TO "YOUR_DB_USER";

-- (Optional) If you want the user to have full control over all future tables in public:
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO "YOUR_DB_USER";
