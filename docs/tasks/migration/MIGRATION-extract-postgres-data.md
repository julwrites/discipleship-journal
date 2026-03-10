---
id: MIGRATION-extract-postgres-data
status: todo
title: Extract PostgreSQL User Data to MySQL-Compatible CSVs
priority: high
created: 2026-03-03
category: migration
dependencies: 
type: task
---

# Extract PostgreSQL User Data to MySQL-Compatible CSVs

## Objective
With the TiDB infrastructure fully provisioned, the next immediate task is bridging the data compatibility gap between PostgreSQL and TiDB (MySQL). `pg_dump` cannot be directly imported. We need a script to extract the data formats correctly.

## Requirements
*   Create a data extraction script (e.g., in Python `scripts/extract_postgres_to_csv.py`) that connects to our existing Google Cloud SQL Postgres database.
*   The script should extract data from user-centric tables (e.g., `users`, `journal_entries`, `tags`, `user_devices`, etc.).
*   **Data Conversion:** 
    *   Convert PostgreSQL `JSONB` to standard MySQL `JSON` strings.
    *   Ensure any Postgres `UUID` types are exported cleanly as `VARCHAR(36)` strings.
*   Output the formatted data into `.csv` files that can be directly consumed by `TiDB Lightning` or standard `LOAD DATA INFILE` tools.
*   The script should *not* mutate any data on the Postgres database, it is pure extraction.

## Next Steps
*   [ ] Write `scripts/extract_postgres_to_csv.py`
*   [ ] Test extraction against a local Postgres dump or staging database configuration.
*   [ ] Verify the output CSVs match the expected TiDB MySQL schema layout.
