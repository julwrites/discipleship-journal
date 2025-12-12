---
id: TEST-002
status: completed
title: Setup End-to-End (E2E) Testing
priority: high
created: 2025-12-11 06:09:10
category: testing
type: task
---

# Setup End-to-End (E2E) Testing

### Description
Implement End-to-End (E2E) testing using Playwright to verify the application's critical paths from a user's perspective. This was deferred from TEST-001.

### Acceptance Criteria
- [x] **Installation**:
    - [x] Playwright installed in `web/` directory.
    - [x] Playwright dependencies/browsers configured.
- [x] **Configuration**:
    - [x] `playwright.config.ts` created and configured to run against the local development server (or a built preview).
    - [x] Added `test:e2e` script to `web/package.json`.
- [x] **Tests**:
    - [x] Create a basic smoke test that:
        - Loads the application.
        - Verifies the page title or landing page content.

### Implementation Notes
- Run tests in headless mode for CI compatibility.
- Use `npm init playwright@latest` or install manually if needed.
- Ensure the dev server is running before tests start (Playwright's `webServer` config).
- Added fallback logic in `useAuth` hook to handle mock environment variables during tests to prevent infinite loading state.
- Created `web/.env` with mock keys for local testing.
