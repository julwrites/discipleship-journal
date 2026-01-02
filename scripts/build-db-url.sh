#!/bin/bash
# Script to build DATABASE_URL from individual components
# Used by the migrator service which expects a single DATABASE_URL

set -e

# Load environment variables from .env files if they exist
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi
if [ -f .env.local ]; then
    export $(grep -v '^#' .env.local | xargs)
fi

# Check required variables
if [ -z "$DB_USERNAME" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    echo "Error: DB_USERNAME, DB_PASSWORD, and DB_NAME are required" >&2
    exit 1
fi

# Build connection string
if [ -n "$CLOUD_SQL_INSTANCE" ]; then
    # Cloud SQL Unix socket connection
    DATABASE_URL="postgres://$DB_USERNAME:$DB_PASSWORD@/$DB_NAME?host=/cloudsql/$CLOUD_SQL_INSTANCE&sslmode=disable"
else
    # Standard TCP connection
    DB_HOST=${DB_HOST:-localhost}
    DB_PORT=${DB_PORT:-5432}
    DATABASE_URL="postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
fi

echo "$DATABASE_URL"