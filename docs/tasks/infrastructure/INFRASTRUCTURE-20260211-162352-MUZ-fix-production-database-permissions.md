---
id: INFRASTRUCTURE-20260211-162352-MUZ
status: completed
title: Fix Production Database Permissions
priority: medium
created: 2026-02-11 16:23:52
category: infrastructure
dependencies:
type: task
---

# Fix Production Database Permissions

## Context
The production deployment on Cloud SQL is failing migrations with the error:
`ERROR: permission denied for schema public (SQLSTATE 42501)`

This is caused by PostgreSQL 15+ changing default permissions. The `public` schema is no longer writable by non-owner users by default.

## Solution Plan
1. Create a script `scripts/fix_postgres_permissions.sql` containing the SQL commands to grant `USAGE` and `CREATE` on schema `public` to the database user.
2. Update `DEPLOYMENT.md` to document this requirement and provide instructions on how to run the fix.
