---
id: FEATURE-001
status: completed
title: User Settings
priority: high
created: 2025-12-11 06:09:10
category: unknown
type: task
---

# User Settings

### Description
Allow users to configure their profile settings.

### Acceptance Criteria
- [x] Database: Users table with `username` and `settings` columns.
- [x] API: Endpoint to read/update user profile (`GET/PUT /api/users/me`).
- [x] UI: Settings page (`/settings`).
- [x] UI: Form to update Username and Bible Version.

### Implementation Status
- ✅ Database schema updated (`api/migrations/000002_schema_updates.up.sql`).
- ✅ Backend handlers implemented and verified (`api/handlers/user.go`).
- ✅ Frontend page matches API (`web/src/pages/Settings.tsx`).
