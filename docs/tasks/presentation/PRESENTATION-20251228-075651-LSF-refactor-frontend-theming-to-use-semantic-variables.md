---
id: PRESENTATION-20251228-075651-LSF
status: completed
title: Refactor Frontend Theming to use Semantic Variables
priority: medium
created: 2025-12-28 07:56:51
category: presentation
dependencies:
type: task
---

# Refactor Frontend Theming to use Semantic Variables

## Context
The frontend currently uses a mix of semantic Shadcn/UI variables (e.g., `bg-muted`) and hardcoded Tailwind colors (e.g., `bg-slate-50`). This causes inconsistency in Dark Mode and makes theming difficult.

## Objectives
- Replace hardcoded color classes with semantic equivalents.
- Ensure Dark Mode looks consistent.

## Implementation Plan
- [x] Refactor `web/src/pages/NoteEditor.tsx`
  - `bg-slate-50` -> `bg-muted`
  - `text-gray-500` -> `text-muted-foreground`
  - `prose-slate` -> `prose-neutral` (or remove if relying on base)
- [x] Refactor `web/src/components/RichTextEditor.tsx`
  - `bg-slate-50` -> `bg-muted`
  - `bg-slate-300` -> `bg-border`
- [x] Refactor `web/src/pages/GroupsPage.tsx`
  - `bg-slate-50` -> `bg-muted`
  - `bg-slate-100` -> `bg-accent` or `bg-muted/50`
  - `text-gray-600/500` -> `text-muted-foreground`
  - `bg-gray-100` -> `bg-muted`
- [x] Refactor `web/src/pages/ChatPage.tsx`
  - `bg-slate-100` -> `bg-muted` or `bg-card`
  - `text-gray-500` -> `text-muted-foreground`
- [x] Refactor `web/src/pages/ConnectionsPage.tsx`
  - `text-gray-500` -> `text-muted-foreground`
- [x] Refactor `web/src/pages/Login.tsx`
  - **Decision:** Kept `bg-slate-900` for the hero section as it is a deliberate design choice (fixed dark mode). Right side auth form uses semantic variables correctly.
- [x] Refactor `web/src/pages/Settings.tsx`
  - Added dark mode support for green status buttons: `dark:bg-green-950 dark:text-green-400`.
- [x] Refactor `web/src/pages/ReadingPlanDetail.tsx`
  - Added dark mode support for green completion indicators: `dark:bg-green-900 dark:text-green-300`.

## Verification
- Lint checks: `npm run lint` (Passed)
- Tests: `npm run test` (Passed)
- Visual check: Ensure components render correctly in both Light and Dark modes.
