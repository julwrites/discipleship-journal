---
id: FEATURE-002
status: completed
title: Dashboard & Journaling
priority: critical
created: 2025-12-11 06:09:10
category: unknown
type: task
---

# Dashboard & Journaling

### Description
Core journaling functionality: Dashboard view, Create Note, Edit Note.

### Acceptance Criteria
- [x] Database: Notes table (JSONB content) (Implemented as `journal_entries`).
- [x] API: CRUD endpoints for notes.
- [x] UI: Dashboard with Note Cards (MRU sorted).
- [x] UI: Search functionality.
- [x] UI: Note Editor (Markdown support).

### Implementation Status
- ✅ Database schema aligned (`api/migrations/000002_schema_updates.up.sql`).
- ✅ Backend handlers implemented and verified (`api/handlers/note.go`).
- ✅ Frontend pages scaffolded (`Dashboard.tsx`, `NoteEditor.tsx`).
