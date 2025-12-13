---
id: DOMAIN-20251213-100746-VAK
status: pending
title: Introduce Backend Service Layer
priority: high
created: 2025-12-13 10:07:46
category: domain
dependencies:
type: task
---

# Introduce Backend Service Layer

## Context
The application's business logic is currently tightly coupled with the HTTP handlers. Specifically, the logic for creating a note is duplicated in `api/handlers/note.go` (standard note creation) and `api/handlers/chat.go` (saving AI chat logs). Both handlers construct manual SQL `INSERT` queries. This duplication makes the code brittle and harder to test.

## Objectives
- Extract business logic from handlers into a dedicated `services` package.
- Eliminate code duplication for note creation.
- Improve testability by decoupling logic from HTTP concerns.

## Requirements
1.  **Create `api/services` Package**:
    -   Initialize a new package for service logic.
2.  **Implement `NoteService`**:
    -   Create `api/services/note_service.go`.
    -   Implement a `CreateNote` method that handles the database insertion.
    -   (Optional) Define an interface for the service to facilitate mocking.
3.  **Refactor Handlers**:
    -   Update `api/handlers/note.go`: Inject or instantiate `NoteService` and use it in `CreateNote`.
    -   Update `api/handlers/chat.go`: Use `NoteService` to save the chat log instead of raw SQL.

## Acceptance Criteria
- [ ] `api/services/note_service.go` exists and contains the note creation logic.
- [ ] `api/handlers/note.go` and `api/handlers/chat.go` no longer contain raw SQL for inserting notes.
- [ ] Both "New Note" and "Chat with AI" features function correctly.
