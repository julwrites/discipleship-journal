---
id: INFRASTRUCTURE-20251213-050200-BRV
status: pending
title: Sync Database Schema
priority: high
created: 2025-12-13 05:02:00
category: infrastructure
dependencies:
type: task
---

# Sync Database Schema

## Context
The `api/database/schema.sql` file is intended to be a snapshot of the current database schema. However, it is currently missing the `connections` table definition, even though the migration `api/migrations/000003_connections_schema.up.sql` exists and has been applied in development.

## Objectives
- Ensure `api/database/schema.sql` accurately reflects the state of the database after all migrations.
- Establish a process or check to keep it in sync (optional, but good to note).

## Requirements
1.  Add the `connections` table definition to `api/database/schema.sql`.
2.  Verify other tables match the migrations.

## Acceptance Criteria
- [ ] `api/database/schema.sql` contains the `CREATE TABLE connections` statement.
- [ ] New environments initialized from `schema.sql` (if any) work correctly with the application.
