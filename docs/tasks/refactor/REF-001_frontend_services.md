# Refactor: Frontend Service Layer for Groups

## Context
The current implementation of `GroupsPage.tsx` makes direct `fetch` calls to the backend API. This violates the architectural decision to use a dedicated Service Layer (`web/src/services/`) for all API interactions.

## Objectives
1.  Create `web/src/services/group.service.ts`.
2.  Move all API calls from `GroupsPage.tsx` to the new service.
3.  Ensure typing is shared/imported correctly.
4.  Update `GroupsPage.tsx` to use the service functions.

## Detailed Tasks
- [ ] Create `web/src/services/group.service.ts`.
- [ ] Define interfaces for Requests/Responses in the service file (or a shared `types.ts`).
- [ ] Implement functions:
    - `getMyGroups(token: string)`
    - `searchGroups(token: string, query: string)`
    - `createGroup(token: string, data: CreateGroupRequest)`
    - `joinGroup(token: string, groupId: string)`
    - `leaveGroup(token: string, groupId: string)`
    - `getGroupMembers(token: string, groupId: string)`
    - `addGroupMember(token: string, groupId: string, userId: string)`
    - `removeGroupMember(token: string, groupId: string, userId: string)`
    - `getGroupShares(token: string, groupId: string)`
    - `getSharedNote(token: string, groupId: string, shareId: string)`
    - `searchUsers(token: string, query: string)` (Note: This might belong in `user.service.ts` if it exists, otherwise `group.service.ts` or `connection.service.ts` is fine for now).
- [ ] Update `GroupsPage.tsx` to import and use these functions.

## Acceptance Criteria
- [ ] No `fetch` calls remain in `GroupsPage.tsx`.
- [ ] Application functionality remains unchanged (verify manually or via tests).
- [ ] Code is linted and passes build.
