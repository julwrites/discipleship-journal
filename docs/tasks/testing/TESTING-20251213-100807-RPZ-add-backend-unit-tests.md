---
id: TESTING-20251213-100807-RPZ
status: completed
title: Add Backend Unit Tests
priority: medium
created: 2025-12-13 10:08:07
category: testing
dependencies: DOMAIN-20251213-100746-VAK
type: task
---

# Add Backend Unit Tests

## Context
The backend codebase (`api/`) currently lacks unit tests, particularly in the `handlers` package. This increases the risk of regressions during refactoring or feature development. The project architecture supports mocking database interactions via `pgxmock`.

## Objectives
- Establish a pattern for unit testing backend components.
- Increase code coverage for critical business logic.

## Requirements
1.  **Setup Test Infrastructure**:
    -   Ensure `github.com/pashagolub/pgxmock/v4` is available and configured.
2.  **Write Tests**:
    -   Create unit tests for the new `NoteService` (created in task `DOMAIN-20251213-100746-VAK`).
    -   (Optional) Create unit tests for `api/handlers/note.go` using a mocked service or database.
3.  **CI Integration**:
    -   Verify that `go test ./...` correctly discovers and runs these tests in the CI pipeline.

## Acceptance Criteria
- [ ] At least one `*_test.go` file is added to `api/services/` or `api/handlers/`.
- [ ] Tests cover success and failure scenarios for Note creation.
- [ ] `go test ./...` passes locally and in CI.
