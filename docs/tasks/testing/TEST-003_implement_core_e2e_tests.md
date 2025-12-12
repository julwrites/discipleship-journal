---
id: TEST-003
status: completed
title: Implement Core E2E Tests
priority: high
created: 2025-12-11 06:09:10
category: testing
type: task
---

# Implement Core E2E Tests

### Description
Implement comprehensive End-to-End (E2E) tests for the core features (Journaling, Connections, Chat) using Playwright. These tests will mock the API and Authentication to verify the Frontend UI logic and flows without needing a full backend environment.

### Acceptance Criteria
- [x] **Auth Mocking**:
    - [x] `useAuth` hook updated to accept a simulated user from `localStorage` in mock mode.
- [x] **Tests Implemented**:
    - [x] `auth.spec.ts`: Verify Login flow (mocked).
    - [x] `journal.spec.ts`: Verify CRUD operations for Notes.
    - [x] `connections.spec.ts`: Verify User Search and Connection requests.
    - [x] `chat.spec.ts`: Verify Chat and AI interaction.
- [x] **Verification**:
    - [x] All E2E tests pass locally (`npm run test:e2e`).

### Implementation Notes
- Use `page.route` in Playwright to intercept API calls (`/api/*`) and return mock data.
- Use `localStorage` to inject the mock user state.
- Ensure tests cover happy paths for all major features.
- Implemented robust mocking for Firebase SDK in `web/src/lib/firebase.ts` to prevent crashes during tests.
