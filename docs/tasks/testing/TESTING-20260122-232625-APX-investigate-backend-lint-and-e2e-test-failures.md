---
id: TESTING-20260122-232625-APX
status: completed
title: Investigate backend lint and e2e test failures
priority: medium
created: 2026-01-22 23:26:25
category: testing
dependencies:
type: task
---

# Investigate backend lint and e2e test failures

Backend lint and e2e tests are both failing. Need to identify root causes and propose fixes.

## Findings

### Backend Lint Failures (golangci-lint)
Run: `golangci-lint run ./...` in `api/` directory

**14 issues:**
1. **errcheck (13 errors)**: Unchecked error return values
   - `handlers/group.go:573:27`: `json.NewEncoder(w).Encode(map[string]string{"id": groupID})`
   - `main.go:93:27`: `secretLoader.Close()`
   - `services/bible_version_service.go:31:23`: `resp.Body.Close()`
   - `services/notification_service_test.go:33:21`: `mockDB.Close(context.Background())`
   - `services/notification_service_test.go:52:21`: `mockDB.Close(context.Background())`
   - `services/notification_service_test.go:73:21`: `mockDB.Close(context.Background())`
   - `services/secrets_test.go:14:13`: `os.Unsetenv("GOOGLE_CLOUD_PROJECT")`
   - `services/secrets_test.go:15:17`: `os.Setenv("GOOGLE_CLOUD_PROJECT", originalProjectID)`
   - `services/secrets_test.go:35:15`: `loader.Close()`
   - `services/secrets_test.go:43:11`: `os.Setenv("TEST_SECRET", expectedValue)`
   - `services/secrets_test.go:44:19`: `os.Unsetenv("TEST_SECRET")`
   - `services/secrets_test.go:59:13`: `os.Unsetenv("NONEXISTENT_SECRET")`
   - `services/secrets_test.go:82:11`: `os.Setenv("TEST_SECRET", "env-value")`

2. **ineffassign (1 error)**: Ineffectual assignment to `err`
   - `handlers/group.go:538:2`: `err = h.db.QueryRow(...).Scan(&myName)` - assignment not used

### E2E Test Failure (Playwright)
Run: `npm run test:e2e` in `web/` directory

**1 failing test:**
- `e2e/groups.spec.ts:159:3`: "Groups (Mocked) › should add a member as admin"
- **Failure**: Test timeout of 30000ms exceeded waiting for `getByPlaceholder('Search by email, name, or username')`
- **Location**: Line 128: `await page.getByPlaceholder('Search by email, name, or username').fill('newuser');`
- **Context**: The test appears to be waiting for a search input that never becomes available/visible.

**Additional notes:**
- 11 other tests pass successfully
- Firebase messaging initialization errors appear in logs but don't cause test failures (expected in mocked environment)
- Modified e2e test files (not staged): `web/e2e/connections.spec.ts`, `web/e2e/sharing.spec.ts`

## Progress

### Backend Lint Fixes
- Fixed all 14 originally reported lint errors plus 5 additional errors discovered
- Errors included unchecked error returns (`errcheck`) and one ineffectual assignment
- Changes made to: `handlers/group.go`, `main.go`, `services/bible_version_service.go`, `services/notification_service_test.go`, `services/secrets_test.go`, `handlers/group_test.go`
- All fixes follow Go conventions: error checking, ignoring errors with `_` where appropriate, logging encoding errors
- Verification: `golangci-lint run ./...` now reports 0 issues

### E2E Test Fix
- **Root cause**: UI changed from search input to connections list. The Add Member dialog now displays accepted connections instead of a search input.
- **Investigation**: Dialog showed "No connections found" because connections endpoint was not mocked.
- **Fix**: Added mock for `/api/connections` returning one accepted connection. Updated test to click "Add" button next to connection instead of filling placeholder.
- **Verification**: All 12 e2e tests pass, including the previously failing test.

### Additional Verification
- Backend unit tests: all pass (`go test ./...`)
- Frontend unit tests: all pass (`npm test`)
- No regressions introduced by fixes.

## Proposed Fixes

### E2E Test Fix

**Root cause**: The test waits for a placeholder input that may not be visible because:
   - The "Add Member" button click doesn't open the dialog (modal)
   - The placeholder text changed
   - The search input is not rendered or has different selector

**Investigation steps:**
1. Check recent UI changes to the "Add Member" dialog
2. Verify the placeholder text matches `'Search by email, name, or username'`
3. Add proper waiting for dialog to appear before attempting to fill
4. Consider adding `await page.locator('[placeholder="Search by email, name, or username"]').waitFor()` before fill

**Potential fix**: Update test to ensure dialog is visible before interacting:
```typescript
// After clicking "Add Member"
await expect(page.locator('text=Add Member to')).toBeVisible(); // Dialog title
await page.getByPlaceholder('Search by email, name, or username').fill('newuser');
```

**Also check**: The mock data for "newuser" may not be set up correctly, causing search to not return results and test to stall.

### Updated Implementation Plan
1. ✓ Fix backend lint errors individually
2. ✓ Run golangci-lint to verify all errors resolved
3. ✓ Investigate and fix e2e test failure
4. ✓ Run e2e tests to confirm fix
5. ✓ Update task status and document changes (status: review_requested)
