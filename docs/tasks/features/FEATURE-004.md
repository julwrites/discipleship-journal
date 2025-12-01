# Task: User Connections & Sharing

## Task Information
- **Task ID**: FEATURE-004
- **Status**: in_progress
- **Priority**: high
- **Phase**: 3
- **Estimated Effort**: 3 days
- **Dependencies**: FEATURE-001 (User Profiles)

## Task Details

### Description
Implement the ability for users to find each other, send connection requests, and view their connections. This is the foundational step for P2P sharing and discipleship features.

### Acceptance Criteria
- [x] Database: `connections` table to store relationships and status (pending/accepted).
- [x] API: Endpoint to search users by email or username (`GET /api/users/search`).
- [x] API: Endpoint to send connection request (`POST /api/connections/request`).
- [x] API: Endpoint to list connections (`GET /api/connections`).
- [x] API: Endpoint to accept/reject request (`PUT /api/connections/{id}`).
- [x] UI: "Connections" page with "Find User" and "My Connections" tabs.
- [x] UI: Display pending requests and allow action.

### Implementation Status
- ✅ Database migration created (`api/migrations/000003_connections_schema.up.sql`).
- ✅ Backend handlers implemented in `api/handlers/connection.go`.
- ✅ Backend routes added in `api/main.go`.
- ✅ Frontend page created `web/src/pages/ConnectionsPage.tsx`.
- ✅ Frontend routing updated in `web/src/App.tsx`.
- ✅ Dashboard updated to link to Connections.

---
*Created: 2025-05-18*
*Status: completed - Implemented and ready for verification*
