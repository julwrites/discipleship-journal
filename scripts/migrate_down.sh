#!/bin/bash

# Load environment variables from .env file if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Determine migrate command
# Force check for file existence first because 'command -v' might give false positives or behave unexpectedly in some environments if not in PATH
if [ -f "/home/user/go/bin/migrate" ]; then
    MIGRATE_CMD="/home/user/go/bin/migrate"
elif [ -f "$HOME/go/bin/migrate" ]; then
    MIGRATE_CMD="$HOME/go/bin/migrate"
elif command -v migrate &> /dev/null; then
    MIGRATE_CMD="migrate"
else
    echo "Error: migrate tool not found in PATH or standard locations."
    exit 1
fi

# Default values if not set in environment
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-postgres}

# Construct the database URL
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "Rolling back migrations for ${DB_NAME} at ${DB_HOST}..."
echo "Using migrate command: $MIGRATE_CMD"

# Run migrations using golang-migrate
$MIGRATE_CMD -path api/migrations -database "${DATABASE_URL}" down 1

if [ $? -eq 0 ]; then
  echo "Migration rollback successful!"
else
  echo "Migration rollback failed."
  exit 1
fi
