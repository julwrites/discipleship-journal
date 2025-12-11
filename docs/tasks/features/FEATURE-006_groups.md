---
id: FEATURE-006
status: unknown
title: Social Groups
priority: high
created: 2025-12-11 06:09:10
category: unknown
type: task
---

# Social Groups

### Description
Implement the ability for users to create and join groups. This allows for small group discipleship and content sharing.

### Acceptance Criteria
- [x] **Database**:
    - [x] `groups` table (id, name, description, created_by).
    - [x] `group_members` table (group_id, user_id, role, joined_at).
- [x] **API**:
    - [x] Create Group (`POST /api/groups`).
    - [x] List My Groups (`GET /api/groups`).
    - [x] Search Groups (`GET /api/groups/search`).
    - [x] Join Group (`POST /api/groups/{id}/join`).
    - [x] List Group Members (`GET /api/groups/{id}/members`).
- [x] **UI**:
    - [x] "Groups" Page with "My Groups" and "Find Groups" tabs.
    - [x] Create Group Dialog.
    - [x] Group Details View (member list).

### Implementation Notes
- Use `pgx` for database interactions.
- Reuse `validate.go` for input validation.
- Role can be 'admin' or 'member'.
- Groups should be public for now (searchable).
