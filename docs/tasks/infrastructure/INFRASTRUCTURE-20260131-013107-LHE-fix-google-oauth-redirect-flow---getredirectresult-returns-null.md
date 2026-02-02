---
id: INFRASTRUCTURE-20260131-013107-LHE
status: completed
title: Fix Google OAuth redirect flow - getRedirectResult returns null
priority: medium
created: 2026-01-31 01:31:07
category: infrastructure
dependencies: 
type: bug
---

# Fix Google OAuth redirect flow - getRedirectResult returns null

## Problem
Users get stuck on 'Loading...' after Google login. The `getRedirectResult()` returns `null` despite `sessionStorage` containing a Firebase redirect event. This prevents authentication state from updating, leaving users on the login page.

## Symptoms from Logs
1. Firebase initializes successfully
2. `sessionStorage` contains `firebase:redirectEvent:` key with OAuth response URL
3. `getRedirectResult()` resolves to `null`
4. `setPersistence` log appears **after** `getRedirectResult` call, indicating race condition
5. No error messages in console; auth state remains `null`

## Root Cause Analysis
The issue is a race condition between Firebase auth persistence setup and redirect result retrieval. The `setPersistence()` call is asynchronous but not awaited before exporting the auth instance. When the login page's `useEffect` calls `getRedirectResult()`, persistence may not yet be configured, causing the redirect state mismatch or failure to retrieve credentials.

Additionally, there may be OAuth 2.0 configuration issues in the Firebase console (staging project), but the logs show Google returning a valid OAuth code, suggesting the OAuth flow itself works.

## Proposed Solution
1. Ensure `setPersistence()` completes before any authentication operations
2. Export a promise from the Firebase module that components can await
3. Modify the login page to wait for persistence before calling `getRedirectResult()` and `signInWithRedirect()`

## Changes Made

### 1. Firebase Module (`web/src/lib/firebase.ts`)
- Added `authReadyPromise: Promise<void> | null = null` variable
- Modified `setPersistence()` call to store the promise in `authReadyPromise`
- Export `authReadyPromise` for components to await
- Ensure `authReadyPromise` is always set (even in error cases) to prevent infinite waiting

### 2. Login Page (`web/src/pages/Login.tsx`)
- Import `authReadyPromise` from Firebase module
- Modified `useEffect` to wait for `authReadyPromise` before calling `getRedirectResult()`
- Modified `handleGoogleLogin` to wait for `authReadyPromise` before calling `signInWithRedirect()`
- Added detailed logging to track the timing of persistence setup

## Testing Required
1. Deploy changes to staging environment
2. Test Google OAuth login flow:
   - Click "Google" button on login page
   - Complete Google OAuth consent screen
   - Should redirect back to app and automatically log in (no "Loading..." stuck)
3. Verify `getRedirectResult()` returns a valid user credential
4. Verify `onAuthStateChanged` fires with authenticated user
5. Test other auth methods (email/password, magic link) still work

## Acceptance Criteria
- [ ] Google OAuth login completes successfully and user is redirected to dashboard
- [ ] No "Loading..." infinite loop after OAuth redirect
- [ ] `getRedirectResult()` returns a valid user credential (not `null`)
- [ ] Auth state updates correctly and `useAuth()` hook returns `user` object
- [ ] Other authentication methods continue to work (email/password, magic link)
- [ ] No regression in existing functionality

## Notes
- If this fix doesn't resolve the issue, further investigation into Firebase console OAuth configuration (authorized domains, OAuth client IDs) will be needed.
- The production environment likely has the same issue; fix should be applied there as well after verification in staging.
