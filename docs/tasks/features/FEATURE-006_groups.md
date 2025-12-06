# Task: Social Groups

## Task Information
- **Task ID**: FEATURE-006
- **Status**: [ ] Pending
- **Priority**: high
- **Phase**: 4 (Post-MVP)
- **Estimated Effort**: 3 days
- **Dependencies**: FEATURE-004 (Connections)

## Task Details

### Description
Implement the ability for users to create and join groups. This allows for small group discipleship and content sharing.

### Acceptance Criteria
- [ ] **Database**:
    - [ ] `groups` table (id, name, description, created_by).
    - [ ] `group_members` table (group_id, user_id, role, joined_at).
- [ ] **API**:
    - [ ] Create Group (`POST /api/groups`).
    - [ ] List My Groups (`GET /api/groups`).
    - [ ] Search Groups (`GET /api/groups/search`).
    - [ ] Join Group (`POST /api/groups/{id}/join`).
    - [ ] List Group Members (`GET /api/groups/{id}/members`).
- [ ] **UI**:
    - [ ] "Groups" Page with "My Groups" and "Find Groups" tabs.
    - [ ] Create Group Dialog.
    - [ ] Group Details View (member list).

### Implementation Notes
- Use `pgx` for database interactions.
- Reuse `validate.go` for input validation.
- Role can be 'admin' or 'member'.
- Groups should be public for now (searchable).

---
*Created: 2025-12-05*
