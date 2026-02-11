# Cloud SQL Setup Guide

This guide details the process for setting up a new PostgreSQL database within an existing Cloud SQL instance for the Discipleship Journal application.

## Prerequisites

- **Cloud SQL Instance**: An existing PostgreSQL 15+ instance.
- **Admin Access**: You must be able to connect as the `postgres` user (or another superuser).
- **Client Tool**: `psql` or Cloud Shell is recommended.

## 1. Connect to the Database

Connect to your Cloud SQL instance using the `postgres` user.

```bash
# Using Cloud SQL Proxy (recommended for local access)
./cloud_sql_proxy -instances=<INSTANCE_CONNECTION_NAME>=tcp:5432

# In another terminal:
psql "host=127.0.0.1 port=5432 user=postgres password=<POSTGRES_PASSWORD> dbname=postgres sslmode=disable"
```

## 2. Create Database and User

If they don't already exist, create the application database and a dedicated user.

```sql
-- Create the database
CREATE DATABASE "discipleship_journal";

-- Create the application user (replace with actual password)
CREATE USER "app_user" WITH PASSWORD 'strong_password';
```

## 3. Configure Extensions & Permissions (CRITICAL)

The application requires the `uuid-ossp` extension for generating UUIDs. **This extension must be enabled by a superuser.** The application user also needs permission to create tables in the `public` schema.

Connect to the **newly created database**:

```bash
psql "host=127.0.0.1 port=5432 user=postgres password=<POSTGRES_PASSWORD> dbname=discipleship_journal sslmode=disable"
```

Run the following SQL commands:

```sql
-- 1. Enable the UUID extension (Must be done by superuser)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Grant usage on the public schema to the application user
GRANT USAGE ON SCHEMA public TO "app_user";

-- 3. Grant table creation privileges on the public schema
GRANT CREATE ON SCHEMA public TO "app_user";

-- 4. Grant privileges on all existing tables and sequences (if any exist from failed attempts)
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "app_user";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "app_user";

-- 5. Ensure future tables created by app_user are owned by app_user (default behavior, but good to verify)
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO "app_user";
```

## 4. Fixing "Dirty database version 1"

If the application failed to start with the error `Dirty database version 1`, it means the initial migration was interrupted (likely due to missing permissions or extension issues).

To fix this state:

1.  **Run the Permission Fixes** from Step 3 above.
2.  **Clean the Migration State**:

    ```sql
    -- Connect to the database as postgres or app_user
    -- Remove the dirty migration entry
    DELETE FROM schema_migrations WHERE version = 1;

    -- Drop any partially created tables to allow a clean start
    DROP TABLE IF EXISTS journal_entries;
    DROP TABLE IF EXISTS users;
    ```

3.  **Restart the Application**: The application will now automatically apply the migrations on startup.

## 5. Verification

To verify the setup is correct, run the following as the `app_user`:

```bash
psql "host=127.0.0.1 port=5432 user=app_user password=<APP_PASSWORD> dbname=discipleship_journal sslmode=disable"
```

```sql
-- Should return a UUID
SELECT uuid_generate_v4();

-- Should allow table creation
CREATE TABLE test_perm (id serial PRIMARY KEY);
DROP TABLE test_perm;
```
