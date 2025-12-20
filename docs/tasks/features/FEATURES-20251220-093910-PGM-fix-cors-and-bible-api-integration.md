---
id: FEATURES-20251220-093910-PGM
status: completed
title: Fix CORS and Bible API Integration
priority: high
created: 2025-12-20 09:39:10
category: features
dependencies: INFRA-004
type: bug
---

# Fix CORS and Bible API Integration

## Bug 1: CORS Errors
Users are experiencing CORS 403 errors when accessing the API from the frontend.
- `GET /api/bible/passage`
- `POST /api/notes`

**Symptoms:**
- `Cross-Origin Request Blocked ... (Reason: CORS preflight response did not succeed). Status code: 403.`

**Potential Causes:**
- Cloud Run service might not be set to allow unauthenticated invocations (blocking OPTIONS requests).
- CORS middleware in Go backend might be misconfigured (allowed origins).

## Bug 2: Bible API Configuration
The backend should integrate with `BibleAIAPI` using `BIBLE_API_URL` and `BIBLE_API_KEY` secrets.

**Requirements:**
- Ensure backend service implementation uses these env vars.
- Ensure deployment passes these secrets to the container.
