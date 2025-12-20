---
id: DOMAIN-20251213-100827-HRO
status: completed
title: Standardize Connections API
priority: low
created: 2025-12-13 10:08:27
category: domain
dependencies:
type: task
---

# Standardize Connections API

## Context
The current implementation of the `/api/connections/{id}` endpoint uses a query parameter `action` (e.g., `?action=accept` or `?action=reject`) to determine the operation. This deviates from standard RESTful design where HTTP verbs should indicate the action.

## Objectives
- Refactor the Connections API to use standard HTTP verbs.
- Simplify the handler logic.

## Requirements
1.  **Refactor Backend**:
    -   Modify `api/handlers/connection.go`.
    -   Implement `PUT /api/connections/{id}` to **accept** a connection request.
    -   Implement `DELETE /api/connections/{id}` to **reject** (or delete) a connection request.
    -   Remove the reliance on the `action` query parameter.
2.  **Update Frontend**:
    -   Update `web/src/pages/ConnectionsPage.tsx` (and/or `web/src/services/api.ts`) to use the new endpoints.

## Acceptance Criteria
- [ ] `PUT /api/connections/{id}` successfully updates the status to `accepted`.
- [ ] `DELETE /api/connections/{id}` successfully removes the connection record.
- [ ] The "Connections" page in the frontend functions correctly with the new API structure.
