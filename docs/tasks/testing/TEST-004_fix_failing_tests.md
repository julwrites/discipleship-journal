# Task: Fix Failing Tests

## Task Information
- **Task ID**: TEST-004
- **Status**: completed
- **Priority**: high
- **Phase**: 4
- **Estimated Effort**: 0.5 days
- **Actual Effort**: 0.1 days
- **Dependencies**: TEST-003
- **Completed**: 2025-02-11

## Task Details

### Description
Fix failing tests reported in CI/CD and local environments.
1. Fix Playwright E2E test failure in `web/e2e/sharing.spec.ts` (`dialog.accept: Target page, context or browser has been closed`).
2. Fix Backend unit test error log in `api/handlers/group_test.go` (`Failed to rollback transaction`).

### Acceptance Criteria
- [x] `web/e2e/sharing.spec.ts` passes consistently.
- [x] `api/handlers/group_test.go` (and all backend tests) pass without error logs.
- [x] All tests pass in local environment.

### Implementation Status

### Completed Work
- ✅ Updated `web/e2e/sharing.spec.ts` to use `page.waitForEvent('dialog')` pattern, ensuring dialog is handled before test completion.
- ✅ Updated `web/playwright.config.ts` to inject `VITE_API_URL` in `webServer` command, ensuring correct API endpoints in tests.
- ✅ Refactored `CreateGroup` in `api/handlers/group.go` to use `committed` flag, preventing `Rollback` calls after successful `Commit`.

---
*Created: 2025-02-11*
*Status: completed*
