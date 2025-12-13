---
id: PRESENTATION-20251213-050215-EQD
status: pending
title: Refactor Frontend API
priority: medium
created: 2025-12-13 05:02:15
category: presentation
dependencies:
type: task
---

# Refactor Frontend API

## Context
`NoteEditor.tsx` contains inline `fetch` calls. Error handling is done via `alert()`. This makes the code hard to maintain and the UX poor.

## Objectives
- Centralize all API calls in `web/src/services/api.ts`.
- Replace `alert()` with a Toast notification system (e.g., `sonner` or `react-hot-toast` or Shadcn/UI Toast).

## Requirements
1.  Add `getNote`, `updateNote`, `deleteNote` to `web/src/services/api.ts`.
2.  Refactor `NoteEditor.tsx` to use these functions.
3.  Install and configure a Toast component.
4.  Replace all `alert("Success")` and `alert("Error")` calls with Toast notifications.

## Acceptance Criteria
- [ ] No direct `fetch` calls in `NoteEditor.tsx`.
- [ ] Application uses Toasts for success/error feedback.
