---
id: TESTING-20260201-144015-ZQI
status: completed
title: Fix backend tests: expand group contract tests
priority: medium
created: 2026-02-01 14:40:15
category: testing
dependencies:
type: task
---

# Fix backend tests: expand group contract tests

## Context
Backend contract tests for groups were failing due to missing route registrations for group share endpoints in the contract test setup and authentication token handling issues in the group share handler.

## Changes Made
1. Added group share handler initialization in `api/tests/contract/setup.go`.
2. Registered group share routes (`/api/groups/{id}/shares`, `/api/groups/{id}/shares/{shareId}`) in the contract test router.
3. Updated `api/handlers/group_share.go` to use `GetUserUUIDFromContext` instead of directly accessing auth token, enabling test overrides.
4. Removed unused imports (`middleware`, `firebase.google.com/go/v4/auth`) from group share handler.

## Verification
- All contract tests pass (including group share tests).
- All unit tests pass.
- No regression in existing functionality.

## Notes
The fix ensures that contract tests for group shares are properly executed and the handlers support both production authentication and test overrides.
