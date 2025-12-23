---
id: FOUNDATION-20251223-033044-LUO
status: pending
title: Optimize GetNotes to exclude content
priority: medium
created: 2025-12-23 03:30:44
category: foundation
dependencies:
type: optimization
---

# Optimize GetNotes to exclude content

## Description
The `GET /api/notes` endpoint currently returns the full content (`jsonb`) of every note in the list.
The frontend Dashboard only displays the Title and UpdatedAt date.
Fetching the full content is inefficient, especially as notes grow in size.

## Acceptance Criteria
- [ ] Modify `NoteService.GetNotes` to only fetch necessary columns (ID, UserID, Title, CreatedAt, UpdatedAt, DeletedAt) or make content fetching optional.
- [ ] Modify `NoteHandler.GetNotes` to respect this change.
- [ ] Ensure `Dashboard` still works correctly.
- [ ] Verify `NoteEditor` (which fetches single note via `GetNote`) still gets the content.
