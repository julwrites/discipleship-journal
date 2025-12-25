---
id: INFRASTRUCTURE-20251224-182428-LQJ
status: completed
title: Refactor Group Feature to Service Pattern
priority: medium
created: 2025-12-24 18:24:28
category: infrastructure
dependencies: []
type: task
---

# Refactor Group Feature to Service Pattern

## Context
The `GroupHandler` (`api/handlers/group.go`) and `GroupShareHandler` (`api/handlers/group_share.go`) currently contain direct database logic and business rules. This violates the architectural principle outlined in `AGENTS.md`: "Keep handler logic separate from business logic (services)."

## Goal
Extract business logic from the handlers into a new `GroupService` (`api/services/group_service.go`). The handlers should responsible for HTTP request parsing, validation, and response formatting, while the service handles database interactions and business rules (e.g., membership checks, notification triggering).

## Scope
1.  **Create `GroupService` interface and implementation**:
    *   Define methods for Group management (Create, List, Search, Join, Leave, GetMembers, AddMember, RemoveMember).
    *   Define methods for Group Sharing (ShareNote, ListShares, GetShareDetails).
2.  **Refactor `GroupHandler`**:
    *   Inject `GroupService` instead of `DBInterface`.
    *   Replace DB calls with Service calls.
3.  **Refactor `GroupShareHandler`**:
    *   Inject `GroupService`.
    *   Replace DB calls with Service calls.
4.  **Update Tests**:
    *   Update handler tests to use a mock `GroupService`.
    *   (Optional but recommended) Add unit tests for `GroupService`.

## Plan
- [x] Create `api/services/group_service.go` defining the interface and implementation.
- [x] Create `api/services/group_service_mock.go` for testing.
- [x] Refactor `api/handlers/group.go` to use `GroupService`.
- [x] Refactor `api/handlers/group_share.go` to use `GroupService`.
- [x] Update `api/main.go` to wire up the new service.
- [x] Verify functionality with tests.
