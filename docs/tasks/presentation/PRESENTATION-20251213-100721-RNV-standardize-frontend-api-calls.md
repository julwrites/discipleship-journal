---
id: PRESENTATION-20251213-100721-RNV
status: completed
title: Standardize Frontend API Calls
priority: medium
created: 2025-12-13 10:07:21
category: presentation
dependencies:
type: task
---

# Standardize Frontend API Calls

## Context
Several key frontend components (`NoteEditor.tsx`, `GroupsPage.tsx`, `ConnectionsPage.tsx`) currently perform inline `fetch` calls to the backend. This leads to code duplication, inconsistent error handling, and scattered authentication logic. A centralized service layer (`web/src/services/api.ts`) exists but is not being used for these features.

## Objectives
- Refactor the identified components to remove all inline `fetch` calls.
- Extend `web/src/services/api.ts` to cover all required API endpoints.
- Ensure consistent authentication token handling.

## Requirements
1.  **Update `web/src/services/api.ts`**:
    -   Add methods for Note operations: `getNote`, `updateNote`.
    -   Add methods for Group operations: `getGroups`, `createGroup`, `joinGroup`, `leaveGroup`, `getGroupMembers`, `addGroupMember`, `removeGroupMember`, `getGroupShares`, `shareNote`, `getSharedNote`.
    -   Add methods for Connection operations: `getConnections`, `sendConnectionRequest`, `respondToConnectionRequest`, `searchUsers`.
2.  **Refactor Components**:
    -   Update `web/src/pages/NoteEditor.tsx` to use the new service methods.
    -   Update `web/src/pages/GroupsPage.tsx` to use the new service methods.
    -   Update `web/src/pages/ConnectionsPage.tsx` to use the new service methods.

## Acceptance Criteria
- [ ] No inline `fetch` calls remain in `NoteEditor.tsx`, `GroupsPage.tsx`, or `ConnectionsPage.tsx`.
- [ ] All API interactions continue to function correctly (verified by E2E tests).
- [ ] `web/src/services/api.ts` is the single source of truth for API calls.
