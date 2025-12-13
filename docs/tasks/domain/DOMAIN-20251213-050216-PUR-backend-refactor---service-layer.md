---
id: DOMAIN-20251213-050216-PUR
status: pending
title: Backend Refactor - Service Layer
priority: medium
created: 2025-12-13 05:02:16
category: domain
dependencies:
type: task
---

# Backend Refactor - Service Layer

## Context
The `ChatWithAI` handler in `api/handlers/chat.go` duplicates the logic for creating a note (inserting into DB). Business logic is tightly coupled with HTTP handlers.

## Objectives
- Extract business logic into a "Service" layer (or at least reusable functions).
- Ensure `CreateNote` and `ChatWithAI` (and any future handlers) use the same logic for creating notes.

## Requirements
1.  Create a `services` package in `api/` (or `internal/services`).
2.  Move `CreateNote` logic (DB insertion, validation) to `services.NoteService`.
3.  Update `handlers.CreateNote` and `handlers.ChatWithAI` to call the service.

## Acceptance Criteria
- [ ] `ChatWithAI` no longer contains SQL `INSERT INTO notes`.
- [ ] `handlers.CreateNote` delegates to the service.
- [ ] Unit tests (if any) are updated.
