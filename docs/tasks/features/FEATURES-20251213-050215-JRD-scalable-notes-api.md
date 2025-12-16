---
id: FEATURES-20251213-050215-JRD
status: completed
title: Scalable Notes API
priority: medium
created: 2025-12-13 05:02:15
category: features
dependencies:
type: task
---

# Scalable Notes API

## Context
The `GET /api/notes` endpoint returns all notes for a user. This will not scale. The Dashboard fetches all notes and filters them client-side.

## Objectives
- Implement pagination and server-side filtering.
- Update Frontend to use these new parameters.

## Requirements
### Backend
1.  Update `GET /api/notes` to accept query parameters:
    - `page` (default 1)
    - `limit` (default 20, max 100)
    - `q` (search query)
2.  Implement SQL pagination (`OFFSET`, `LIMIT`) and search (`WHERE title ILIKE ...`).
3.  Return metadata (total count, total pages) in the response (either in headers or a wrapped envelope).

### Frontend
1.  Update `Dashboard.tsx` and `api.ts` to pass these parameters.
2.  Implement "Load More" or "Next Page" UI in the Dashboard.
3.  Use the server-side search when the user types in the search bar.

## Acceptance Criteria
- [ ] `GET /api/notes?limit=10` returns only 10 notes.
- [ ] `GET /api/notes?q=test` returns only matching notes.
- [ ] Dashboard displays notes progressively or by page.
- [ ] Search functionality works via the API.
