#!/bin/bash
set -e

# Prompt for database connection details
read -p "Database Host [localhost]: " DB_HOST
DB_HOST=${DB_HOST:-localhost}

read -p "Database Port [5432]: " DB_PORT
DB_PORT=${DB_PORT:-5432}

read -p "Database Name [journal]: " DB_NAME
DB_NAME=${DB_NAME:-journal}

read -p "Superuser (postgres) Name [postgres]: " DB_SUPERUSER
DB_SUPERUSER=${DB_SUPERUSER:-postgres}

read -p "Application User Name [discipleship-journal]: " DB_USER
DB_USER=${DB_USER:-discipleship-journal}

# Read SQL template
SQL_FILE="scripts/fix_dirty_migration.sql"
if [ ! -f "$SQL_FILE" ]; then
    echo "Error: $SQL_FILE not found."
    exit 1
fi

echo "Running permissions fix script as '$DB_SUPERUSER' for user '$DB_USER' on database '$DB_NAME'..."

# Replace placeholders and pipe to psql
sed "s/YOUR_DB_USER/$DB_USER/g" "$SQL_FILE" | \
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_SUPERUSER" -d "$DB_NAME" -f -

echo "Permissions fixed successfully. You can now retry your migration."
