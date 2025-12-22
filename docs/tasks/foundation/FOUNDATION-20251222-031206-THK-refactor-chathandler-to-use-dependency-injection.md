---
id: FOUNDATION-20251222-031206-THK
status: completed
title: Refactor ChatHandler to use Dependency Injection
priority: medium
created: 2025-12-22 03:12:06
category: foundation
dependencies:
type: task
---

# Refactor ChatHandler to use Dependency Injection

The `ChatHandler` currently relies on the global `database.DB` instance via the `GetUserUUID` helper function. This makes it difficult to unit test the handler in isolation, as it requires a live database connection or global state manipulation.

To improve testability and architectural cleanliness, we will refactor `ChatHandler` to accept a `DBInterface` via dependency injection.

## Plan

1.  **Modify `ChatHandler` struct**: Add a `DB` field of type `DBInterface`.
2.  **Update `NewChatHandler`**: Change the signature to accept `DBInterface` and initialize the struct field.
3.  **Refactor `ChatWithAI`**:
    *   Remove the call to `GetUserUUID`.
    *   Implement the user lookup logic directly in the handler using `h.DB.QueryRow`.
4.  **Update `main.go`**: Inject the `database.DB` instance when creating `NewChatHandler`.
5.  **Update Tests**:
    *   Update `api/handlers/chat_test.go` to mock the DB interactions using `pgxmock`.
    *   Add a test case for `ChatWithAI` that verifies the DB query and Note creation.
