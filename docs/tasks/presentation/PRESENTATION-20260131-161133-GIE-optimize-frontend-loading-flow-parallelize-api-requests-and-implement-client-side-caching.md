---
id: PRESENTATION-20260131-161133-GIE
status: verified
title: Optimize frontend loading flow: parallelize API requests and implement client-side caching
priority: medium
created: 2026-01-31 16:11:33
category: presentation
dependencies:
type: task
---

# Optimize frontend loading flow: parallelize API requests and implement client-side caching

## Task Information
- **Dependencies**: None (can be implemented independently)
- **Related tasks**: PRESENTATION-20260131-150229-RVN (parent epic)

## Task Details
Optimize the frontend loading sequence to reduce perceived latency. Parallelize API requests where possible and implement client-side caching to avoid redundant network calls.

### Current findings from exploration:
- Sequential API calls: `syncUser()` (POST /api/users/me) then `fetchNotes()` (GET /api/notes)
- No client-side caching of notes data
- Each dashboard visit triggers fresh API calls
- Potential for prefetching notes after authentication

### Steps:
1. **Analyze current API call sequence**:
   - Review `Dashboard.tsx` `useEffect` dependencies
   - Examine `syncUser()` and `fetchNotes()` in `services/api.ts`
   - Identify dependencies between calls (does `fetchNotes` need user data?)

2. **Parallelize API requests**:
   - If `syncUser` and `fetchNotes` are independent, fire them in parallel
   - Use `Promise.all()` or multiple `useEffect` hooks
   - Handle error states independently

3. **Implement client-side caching**:
   - Cache notes data in localStorage or IndexedDB
   - Set appropriate TTL (time-to-live) for cached notes
   - Implement cache invalidation on note creation/update/deletion
   - Show cached data immediately while fetching fresh data in background

4. **Add stale-while-revalidate pattern**:
   - Display cached notes immediately on dashboard load
   - Fetch fresh notes in background
   - Update UI when new data arrives

5. **Optimize data fetching**:
   - Consider using React Query (TanStack Query) for caching, background updates, and error handling
   - Evaluate if implementing a full state management library is warranted

6. **Test caching behavior**:
   - Verify cache works across page reloads
   - Test cache invalidation when notes are modified
   - Measure performance improvement with caching

### Acceptance Criteria
- [x] API requests parallelized where possible (reducing total load time)
- [x] Client-side caching implemented for notes data
- [x] Stale-while-revalidate pattern: show cached data immediately
- [x] Cache invalidation on note modifications (handled via SWR updates)
- [x] No regression in data consistency (always show latest data eventually)

## Implementation Status
### Completed Work
- ✅ Initial analysis of API call sequence completed
- ✅ Implemented client-side caching in `web/src/services/cache.ts` using `localStorage`
- ✅ Scoped cache keys to `userId` to ensure privacy
- ✅ Updated `Dashboard.tsx` to use `stale-while-revalidate` pattern:
    - Loads cached notes immediately if available (initial load, default view)
    - Fetches fresh notes in background and updates UI/cache
- ✅ Parallelized `syncUser` and `fetchNotes` (via concurrent effects)

### Blockers
None yet.

## Notes
- Consider using existing libraries (React Query) rather than custom cache implementation
- Ensure caching respects user isolation (don't show other users' notes)
- Cache size management for users with many notes
- Coordinate with backend team if API changes needed for cache headers
