# Task: Setup Testing Infrastructure

## Task Information
- **Task ID**: TEST-001
- **Status**: pending
- **Priority**: critical
- **Phase**: 1.5
- **Estimated Effort**: 2 days
- **Dependencies**: FOUNDATION-001

## Task Details

### Description
Establish a comprehensive testing strategy for both Frontend and Backend to ensure code quality and prevent regressions. This includes unit testing frameworks, integration testing support, and end-to-end (E2E) testing setup.

### Acceptance Criteria
- [ ] **Frontend (Unit/Component)**:
    - [ ] Install `vitest`, `@testing-library/react`, `@testing-library/dom`.
    - [ ] Configure `vite.config.ts` for testing.
    - [ ] Create a sample test for a simple component (e.g., `App.tsx` or a Button).
    - [ ] Add `test` script to `web/package.json`.
- [ ] **Backend (Unit/Integration)**:
    - [ ] Verify `go test ./...` works.
    - [ ] Install `testcontainers-go` for integration tests (Postgres).
    - [ ] Create a sample integration test for a database function.
- [ ] **End-to-End (E2E)**:
    - [ ] Initialize Playwright in `web/`.
    - [ ] Configure Playwright to run against local dev server.
    - [ ] Create a basic "smoke test" (load page, check title).

### Implementation Notes
- **Frontend**: Vitest is preferred over Jest for Vite projects due to speed and shared configuration.
- **Backend**: Use table-driven tests for unit tests. Use Testcontainers for DB tests to avoid mocking database logic.
- **CI**: These tests must be runnable in CI (headless mode for browsers).

---

*Created: 2025-05-18*
*Status: pending*
