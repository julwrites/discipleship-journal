---
id: PRESENTATION-20251222-145053-BKU
status: completed
title: Fix Google Sign-In Cross-Origin-Opener-Policy Error
priority: medium
created: 2025-12-22 14:50:53
category: presentation
dependencies:
type: bug
---

# Fix Google Sign-In Cross-Origin-Opener-Policy Error

## Description
Google Sign-In with popup requires the `Cross-Origin-Opener-Policy` (COOP) header to be set to `same-origin-allow-popups` to properly communicate between the popup and the main window. Without this, the popup might be blocked or fail to return the credential to the opener.

## Implementation
The header `Cross-Origin-Opener-Policy: same-origin-allow-popups` has been added to:
1.  `firebase.json`: For production hosting on Firebase.
2.  `web/nginx.conf`: For containerized deployment (e.g., Cloud Run or local Docker).
3.  `web/vite.config.ts`: For local development server.

## Verification
- **Automated Tests**: Frontend tests pass (`npm run test`).
- **Manual Verification**:
    - Started local dev server (`npm run dev`).
    - Verified header presence using `curl -I http://localhost:5173`:
      ```
      HTTP/1.1 200 OK
      ...
      Cross-Origin-Opener-Policy: same-origin-allow-popups
      ...
      ```
    - This confirms the header is correctly served in the development environment.
    - Production and Docker environments use similar configurations which are visually verified in the config files.

## Status
- [x] Update `firebase.json`
- [x] Update `web/nginx.conf`
- [x] Update `web/vite.config.ts`
- [x] Verify in local environment
