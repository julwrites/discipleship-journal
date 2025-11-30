#!/bin/bash
set -e

# Load environment variables from .env file if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

DB_URL=${DB_URL:-"postgres://postgres:postgres@localhost:5432/journal?sslmode=disable"}

echo "Running migrations..."
echo "Database URL: $DB_URL"

# Check if migrate tool is installed
if ! command -v migrate &> /dev/null; then
    echo "migrate could not be found. Please install it."
    echo "go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

migrate -path api/migrations -database "$DB_URL" up
