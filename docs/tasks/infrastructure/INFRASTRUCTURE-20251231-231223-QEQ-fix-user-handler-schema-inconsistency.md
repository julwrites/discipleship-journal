---
id: INFRASTRUCTURE-20251231-231223-QEQ
status: verified
title: Fix User handler schema inconsistency
priority: medium
created: 2025-12-31 23:12:23
category: infrastructure
dependencies: 
type: bug
---

# Fix User handler schema inconsistency

## Problem
The `CreateOrUpdateUser` handler was overwriting the `settings` JSONB column instead of merging new settings with existing ones. This caused data loss when updating specific settings (e.g. `bible_version` update would wipe out `theme`).

## Solution
Modified `api/handlers/user.go` to:
1. Unmarshal existing `settings` from the database.
2. Merge incoming settings from the request into the existing map.
3. Save the merged settings map back to the database.

## Verification
- Added a regression test `TestUserHandler_CreateOrUpdateUser_MergesSettings` in `api/handlers/user_test.go`.
- Verified that the test passes and existing tests still pass.
