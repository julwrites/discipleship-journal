# Task: User Settings

## Task Information
- **Task ID**: FEATURE-001
- **Status**: in_progress
- **Priority**: high
- **Phase**: 2
- **Estimated Effort**: 2 days
- **Dependencies**: FOUNDATION-002

## Task Details

### Description
Allow users to configure their profile settings.

### Acceptance Criteria
- [ ] Database: Users table with `username` and `settings` columns.
- [ ] API: Endpoint to read/update user profile (`GET/PUT /api/users/me`).
- [ ] UI: Settings page (`/settings`).
- [ ] UI: Form to update Username and Bible Version.

### Implementation Status
- ✅ Database schema defined (`api/database/schema.sql`).
- ✅ Backend handlers scaffolded (`api/handlers/user.go`).
- ✅ Frontend page scaffolded (`web/src/pages/Settings.tsx`).

---
*Created: 2025-05-18*
*Status: in_progress - Code exists, verification pending*
