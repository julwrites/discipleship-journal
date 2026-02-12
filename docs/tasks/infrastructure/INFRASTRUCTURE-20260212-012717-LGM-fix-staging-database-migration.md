---
id: INFRASTRUCTURE-20260212-012717-LGM
status: completed
title: Fix staging database migration
priority: medium
created: 2026-02-12 01:27:17
category: infrastructure
dependencies:
type: task
---

# Fix staging database migration

## Problem
The staging backend service is failing to start due to a dirty database migration state.
Logs:
```
error="failed to apply migrations: Dirty database version 1. Fix and force version."
```
Also, there is a warning about invalid JSON in `LLM_SYSTEM_PROMPTS`:
```
level=WARN msg="Failed to parse system prompts JSON" error="invalid character '\\n' in string literal"
```

## Solution
1. Create a SQL script `scripts/fix_dirty_migration.sql` to clean the dirty state for version 1 and attempt to create the `uuid-ossp` extension.
2. Create a helper shell script `scripts/fix_staging_db.sh` to execute the SQL against the staging database using `gcloud` or a direct connection string.
3. Document the troubleshooting steps for future reference.

## Sub-tasks
- [ ] Create `scripts/fix_dirty_migration.sql`
- [ ] Create `scripts/fix_staging_db.sh`
- [ ] Document in `docs/troubleshooting.md` (or similar) or update `README.md`
