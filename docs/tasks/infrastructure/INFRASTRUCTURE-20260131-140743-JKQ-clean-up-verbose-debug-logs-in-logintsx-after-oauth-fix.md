---
id: INFRASTRUCTURE-20260131-140743-JKQ
status: completed
title: Clean up verbose debug logs in Login.tsx after OAuth fix
priority: medium
created: 2026-01-31 14:07:43
category: infrastructure
dependencies:
type: task
---

# Clean up verbose debug logs in Login.tsx after OAuth fix

## Problem
After fixing Google OAuth popup flow for Edge/Chrome, the `Login.tsx` component contains extensive debug logging that clutters the console in production. These logs were added during debugging but should be cleaned up now that the OAuth flow is working.

## Symptoms
1. Verbose `console.log` statements for every authentication step
2. Full URL parameter dumps and sessionStorage inspection
3. Service worker registration/unregistration logs
4. Multiple redundant logs about auth persistence and redirect handling

## Goals
1. Reduce console noise in production while preserving debugging capability
2. Keep error and warning logs for actual issues
3. Address service worker unregistration errors
4. Maintain development-time debugging when needed

## Implementation Plan

### 1. Development-Only Logging
Wrap debug `console.log` statements with `import.meta.env.DEV` checks:
```typescript
if (import.meta.env.DEV) {
  console.log("Debug message");
}
```

### 2. Log Categories to Handle
- **Keep (always)**: `console.error` for actual failures, `console.warn` for warnings
- **Development-only**: Debug `console.log` statements for:
  - URL parameter inspection
  - sessionStorage dumps
  - Auth state changes
  - Service worker operations
  - OAuth flow steps
- **Remove entirely**: Redundant or overly verbose logs

### 3. Service Worker Error Handling
The `InvalidStateError: Only the active worker can claim clients` occurs when trying to unregister active service workers. Should:
- Catch and suppress this specific error
- Or skip unregistration if not needed (popup flow works without it)

### 4. Error Message Formatting
Keep helpful error formatting for auth failures but reduce verbosity.

## Acceptance Criteria
- [ ] Production builds show minimal debug logging
- [ ] Development builds retain debugging capability
- [ ] Error and warning logs remain for troubleshooting
- [ ] No `InvalidStateError` from service worker unregistration
- [ ] Console output is clean and professional in production

## Notes
- Use Vite's `import.meta.env.DEV` for environment detection
- Consider creating a simple logger utility if similar cleanup needed elsewhere
- Test both development and production builds
