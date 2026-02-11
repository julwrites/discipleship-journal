#!/bin/bash
set -e

# Load environment variables from .env file if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

DB_URL=${DB_URL:-"postgres://postgres:postgres@localhost:5432/discipleship_journal?sslmode=disable"}
DB_USER=${DB_USER:-"postgres"}

echo "Verifying database setup..."
echo "Database URL: $DB_URL"
echo "Database User: $DB_USER"

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    echo "psql could not be found. Please install postgresql-client."
    exit 1
fi

# Run verification queries
# Note: We use -c with multiple statements separated by semicolons.
psql "$DB_URL" -c "
SELECT current_user, current_database();
SELECT 'Extension check: ' || (CASE WHEN EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp') THEN 'OK' ELSE 'MISSING (uuid-ossp)' END);
SELECT 'Schema usage check: ' || (CASE WHEN has_schema_privilege('$DB_USER', 'public', 'USAGE') THEN 'OK' ELSE 'MISSING (USAGE)' END);
SELECT 'Schema create check: ' || (CASE WHEN has_schema_privilege('$DB_USER', 'public', 'CREATE') THEN 'OK' ELSE 'MISSING (CREATE)' END);
"

echo "Verification complete."
