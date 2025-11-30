# Task: Dashboard & Journaling

## Task Information
- **Task ID**: FEATURE-002
- **Status**: in_progress
- **Priority**: critical
- **Phase**: 2
- **Estimated Effort**: 3 days
- **Dependencies**: FOUNDATION-002, FEATURE-001

## Task Details

### Description
Core journaling functionality: Dashboard view, Create Note, Edit Note.

### Acceptance Criteria
- [ ] Database: Notes table (JSONB content).
- [ ] API: CRUD endpoints for notes.
- [ ] UI: Dashboard with Note Cards (MRU sorted).
- [ ] UI: Search functionality.
- [ ] UI: Note Editor (Markdown support).

### Implementation Status
- ✅ Database schema defined.
- ✅ Backend handlers scaffolded (`api/handlers/notes.go`).
- ✅ Frontend pages scaffolded (`Dashboard.tsx`, `NoteEditor.tsx`).

---
*Created: 2025-05-18*
*Status: in_progress - Code exists, verification pending*
