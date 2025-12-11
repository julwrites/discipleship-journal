---
id: TEST-001
status: completed
title: Setup Testing Infrastructure
priority: critical
created: 2025-12-11 06:09:10
category: unknown
type: task
---

# Setup Testing Infrastructure

### Description
Establish a comprehensive testing strategy for both Frontend and Backend to ensure code quality and prevent regressions. This includes unit testing frameworks, integration testing support, and end-to-end (E2E) testing setup.

### Acceptance Criteria
- [x] **Frontend (Unit/Component)**:
    - [x] Install `vitest`, `@testing-library/react`, `@testing-library/dom`.
    - [x] Configure `vite.config.ts` for testing.
    - [x] Create a sample test for a simple component (e.g., `App.tsx` or a Button).
    - [x] Add `test` script to `web/package.json`.
- [x] **Backend (Unit/Integration)**:
    - [x] Verify `go test ./...` works.
    - [x] Install `testcontainers-go` for integration tests (Postgres).
    - [x] Create a sample integration test for a database function.
- [ ] **End-to-End (E2E)**:
    - [ ] Initialize Playwright in `web/`.
    - [ ] Configure Playwright to run against local dev server.
    - [ ] Create a basic "smoke test" (load page, check title).

### Implementation Notes
- **Frontend**: Vitest is preferred over Jest for Vite projects due to speed and shared configuration.
- **Backend**: Use table-driven tests for unit tests. Use Testcontainers for DB tests to avoid mocking database logic.
- **CI**: These tests must be runnable in CI (headless mode for browsers).

### Status Update
- Frontend testing set up with Vitest and React Testing Library. Dummy test passing.
- Backend testing set up. `testcontainers` installed. Integration test created but commented out in sandbox due to missing Docker. Simple unit test verified.
- E2E setup deferred to a separate task or later step.
