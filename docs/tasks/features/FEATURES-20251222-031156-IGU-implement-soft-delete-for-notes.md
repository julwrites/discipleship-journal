---
id: FEATURES-20251222-031156-IGU
status: completed
title: Implement Soft Delete for Notes
priority: medium
created: 2025-12-22 03:11:56
category: features
dependencies:
type: task
---

# Implement Soft Delete for Notes

Implemented soft delete mechanism for notes to allow recovery and prevent accidental permanent data loss.

## Implementation
- Added `deleted_at` column to `notes` table (migration `000007_soft_delete_notes`).
- Updated `NoteService` to:
    - `DeleteNote`: Set `deleted_at` to `NOW()`.
    - `GetNote`, `GetNotes`: Filter by `deleted_at IS NULL`.
    - `UpdateNote`: Ensure `deleted_at IS NULL`.
    - `Note` struct: Added `DeletedAt` field.
- Updated `GroupShareHandler` to filter shared notes that are soft deleted.

## Verification
- Updated `api/services/note_service_test.go` to verify soft delete logic (UPDATE instead of DELETE, filtering queries).
- Ran backend tests `go test ./...` - All passed.
