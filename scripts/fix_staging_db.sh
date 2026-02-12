#!/bin/bash
set -e

# Default configuration from known staging setup
DEFAULT_INSTANCE="discipleship-journal-52a2c:asia-southeast1:discipleship-journal-postgresql"
DEFAULT_DB="discipleship_journal_staging"

echo "=== Staging Database Fix Utility ==="
echo "This script helps fix the 'Dirty database version 1' error on staging."
echo ""

# Prompt for Instance
read -p "Enter Cloud SQL Instance Connection Name [$DEFAULT_INSTANCE]: " INSTANCE
INSTANCE=${INSTANCE:-$DEFAULT_INSTANCE}

# Prompt for Database Name
read -p "Enter Database Name [$DEFAULT_DB]: " DB_NAME
DB_NAME=${DB_NAME:-$DEFAULT_DB}

# Prompt for User
read -p "Enter Database User (should be 'postgres' or superuser) [postgres]: " DB_USER
DB_USER=${DB_USER:-postgres}

echo ""
echo "Target:"
echo "  Instance: $INSTANCE"
echo "  Database: $DB_NAME"
echo "  User:     $DB_USER"
echo ""
read -p "Proceed? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

echo "Connecting to database and running fix script..."

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo "Error: 'gcloud' CLI is not installed or not in PATH."
    exit 1
fi

# Run the SQL script using gcloud sql connect
# Note: This requires the Cloud SQL Auth Proxy to be installed or available to gcloud,
# or for the user to be authorized to connect directly.
cat scripts/fix_dirty_migration.sql | gcloud sql connect "$INSTANCE" --user="$DB_USER" --database="$DB_NAME" --quiet

echo ""
echo "✅ Script execution completed."
echo "If you saw 'DELETE 1' and 'CREATE EXTENSION', it succeeded."
echo "If you saw 'permission denied', please try again with a superuser account."
