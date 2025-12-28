---
id: TESTING-20251225-033202-MIH
status: completed
title: Add unit tests for Groups and Notifications
priority: high
created: 2025-12-25 03:32:02
category: testing
dependencies: []
type: task
---

# Add unit tests for Groups and Notifications

## Context
The "Groups" and "Notifications" features are fully implemented and documented, but they lack backend unit tests. This leaves the codebase vulnerable to regressions.
Additionally, the `GroupHandler` currently mixes dependency injection (`h.db`) with global state (`database.DB`) for transactions, making it difficult to test.

## Objectives
1.  **Refactor GroupHandler**: Ensure `CreateGroup` uses the injected `DBInterface` (via type assertion or interface expansion) to begin transactions, rather than `database.DB.Begin()`.
2.  **Add GroupHandler Tests**: Create `api/handlers/group_test.go` covering:
    *   Create Group (Success/Failure)
    *   List My Groups
    *   Search Groups
    *   Join/Leave Group
    *   Member Management (Get, Add, Remove)
3.  **Add NotificationHandler Tests**: Create `api/handlers/notification_test.go`.
4.  **Add NotificationService Tests**: Create `api/services/notification_service_test.go`.

## Implementation Details
*   **Transactions in Mocks**: The `DBInterface` in `api/handlers/interfaces.go` might need to include `Begin(ctx context.Context) (pgx.Tx, error)` to support mocking transactions.
*   **Mocking**: Use `pgxmock` for DB interactions and `MockNotificationService` for notification logic.
*   **Test Data**: Use `testUserKey` context injection for authentication bypassing in tests.

## Verification
*   Run `go test ./api/handlers/...` and ensure all new tests pass.
*   Run `go test ./api/services/...`.
