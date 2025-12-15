---
id: FEATURES-20251213-050215-QMJ
status: completed
title: Implement Delete Note
priority: high
created: 2025-12-13 05:02:15
category: features
dependencies:
type: story
---

# Implement Delete Note

## Context
Users currently cannot delete notes they have created. This is a basic CRUD requirement for a journaling application.

## Objectives
- Allow authenticated users to delete their own notes.
- Ensure proper cleanup of related data (though `ON DELETE CASCADE` in DB handles this, API should confirm).

## Requirements
### Backend
1.  Create `DELETE /api/notes/{id}` endpoint.
2.  Verify the user owns the note before deleting.
3.  Return 200 OK or 204 No Content on success.

### Frontend
1.  Add a "Delete" button to the `NoteEditor` (e.g., in the toolbar or a dropdown).
2.  Add a "Delete" option to the note cards on the Dashboard (optional but good).
3.  Implement a confirmation dialog ("Are you sure?").
4.  On success, redirect to the Dashboard or remove the note from the list.

## Acceptance Criteria
- [x] `DELETE /api/notes/{id}` deletes the note from the database.
- [x] Users cannot delete notes belonging to others (return 403).
- [x] "Delete" button is visible in the UI.
- [x] Confirmation dialog prevents accidental deletion.
- [x] Deleting a note redirects user appropriately.
