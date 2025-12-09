# Task: Group Sharing

## Task Information
- **Task ID**: FEATURE-007
- **Status**: [x] Completed
- **Priority**: high
- **Phase**: 4 (Post-MVP)
- **Estimated Effort**: 2 days
- **Dependencies**: FEATURE-006 (Groups)

## Task Details

### Description
Allow users to share their journal notes with a group. This facilitates group discussions and accountability.

### Acceptance Criteria
- [x] **Database**:
    - [x] `group_shares` table (id, group_id, note_id, shared_by, shared_at).
- [x] **API**:
    - [x] Share Note (`POST /api/groups/{id}/shares`).
    - [x] List Group Shares (`GET /api/groups/{id}/shares`).
    - [x] Get Shared Note (`GET /api/groups/{id}/shares/{shareId}`).
- [x] **UI**:
    - [x] "Share" button in Note Editor / Dashboard.
    - [x] Select Group Dialog.
    - [x] "Shared Notes" tab in Group Details.
    - [x] View Shared Note page (Read-only for members).

### Implementation Notes
- A shared note is a reference to an existing note.
- If the original note is updated, the share reflects the update (unless we snapshot it, but reference is better for now).
- If the original note is deleted, the share should be deleted (CASCADE).
- Only members of the group can view the shared notes.

---
*Created: 2025-05-18*
*Completed: 2025-05-18*
