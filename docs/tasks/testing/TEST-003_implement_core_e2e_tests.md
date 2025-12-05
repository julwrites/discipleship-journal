# Task: Implement Core E2E Tests

## Task Information
- **Task ID**: TEST-003
- **Status**: Completed
- **Priority**: high
- **Phase**: 3.5
- **Estimated Effort**: 1 day
- **Dependencies**: TEST-002

## Task Details

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

---
*Created: 2025-05-21*
*Status: Completed - All core E2E tests implemented and verified.*
