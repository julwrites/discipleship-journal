---
id: INFRASTRUCTURE-20260211-164121-GYT
status: completed
title: Fix Database Migration Failure: Dirty Database Version
priority: medium
created: 2026-02-11 16:41:21
category: infrastructure
dependencies:
type: task
---

# Fix Database Migration Failure: Dirty Database Version

## Context
The production deployment is failing with:
`failed to apply migrations: Dirty database version 1. Fix and force version.`

This occurs when a migration fails mid-execution (e.g., due to permission errors), leaving the `schema_migrations` table in a "dirty" state. The `golang-migrate` library refuses to run subsequent migrations until this flag is cleared.

## Solution Plan
1. Create a script `scripts/fix_dirty_migration.sql` to clear the dirty flag.
2. Update `DEPLOYMENT.md` with instructions on how to use it.
