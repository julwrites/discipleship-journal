---
id: INFRASTRUCTURE-20260131-101320-SOU
status: review_requested
title: Fix Google OAuth on Edge/Chrome - getRedirectResult returns null
priority: medium
created: 2026-01-31 10:13:20
category: infrastructure
dependencies: 
type: task
---

# Fix Google OAuth on Edge/Chrome - getRedirectResult returns null

## Problem
Google OAuth works on Safari but fails on Microsoft Edge and Chrome. After redirect from Google, `getRedirectResult()` returns `null` despite `sessionStorage` containing a Firebase redirect event. Users get stuck on 'Loading...' screen.

## Symptoms from Logs
1. Firebase initializes successfully, persistence set to `browserLocalPersistence`.
2. `sessionStorage` contains `firebase:redirectEvent:` key with OAuth response URL.
3. `getRedirectResult()` resolves to `null`.
4. Auth state remains `null`.
5. No error messages; OAuth flow appears to complete on Google side.

## Root Cause Analysis
The recent fix for race condition between `setPersistence()` and `getRedirectResult()` (task INFRASTRUCTURE-20260131-013107-LHE) has been implemented but may not fully resolve the issue on Edge/Chrome due to browser-specific behaviors:

1. **Service Worker Interception**: The service worker denylist includes `/^\/__\/auth\/handler/` but Edge/Chrome may still intercept due to different service worker lifecycle.
2. **Cookie Policies**: Edge/Chrome have stricter default cookie settings (blocking third-party cookies, SameSite enforcement) that may affect OAuth flow.
3. **Cross-Origin Headers**: The `Cross-Origin-Opener-Policy: same-origin-allow-popups` header may affect sessionStorage accessibility across redirects.
4. **Stale Redirect Events**: Previous failed OAuth attempts may leave stale redirect events in `sessionStorage` that Firebase cannot process.
5. **Browser-Specific Firebase SDK Bugs**: Potential bugs in Firebase v12.6.0 with Edge/Chrome's implementation of `sessionStorage` or redirect handling.

## Changes Made (Implementation)
1. **SessionStorage Cleanup**: Modified `web/src/pages/Login.tsx` to automatically clear stale Firebase redirect events from `sessionStorage` when `getRedirectResult()` returns `null`. This prevents stale events from blocking subsequent OAuth attempts.

## Proposed Solutions
Test the following hypotheses in order:

### 1. Clear Stale Redirect Events (IMPLEMENTED)
Implemented in `Login.tsx`. Automatically clears stale Firebase redirect events from `sessionStorage` when `getRedirectResult()` returns `null`.

### 2. Disable Service Worker for Testing
Temporarily disable service worker registration to rule out interference. Could comment out PWA plugin in `vite.config.ts` or unregister service workers.

### 3. Adjust Cross-Origin Headers
Temporarily remove `Cross-Origin-Opener-Policy` header to test if it affects sessionStorage isolation.

### 4. Use Popup Flow Instead of Redirect
Test `signInWithPopup()` as alternative (though popups may be blocked by browsers). This bypasses redirect flow issues.

### 5. Cookie Configuration
Ensure Firebase OAuth client is configured with correct authorized domains and that cookies are allowed with appropriate SameSite settings.

## Implementation Plan
1. **Add sessionStorage cleanup** in `Login.tsx` after `getRedirectResult()` returns `null`.
2. **Deploy to staging** and test on Edge/Chrome.
3. If still failing, **disable service worker** and retest.
4. If still failing, **adjust headers** (temporarily remove COOP).
5. If still failing, **implement popup fallback** with user option.
6. **Update Firebase console** OAuth configuration to ensure authorized domains include staging URL.

## Testing Required
- Test Google OAuth on Microsoft Edge, Chrome, and Safari (control).
- Verify `getRedirectResult()` returns valid user credential.
- Verify `onAuthStateChanged` fires with authenticated user.
- Ensure other authentication methods (email/password, magic link) still work.

## Acceptance Criteria
- [ ] Google OAuth completes successfully on Edge and Chrome.
- [ ] No "Loading..." infinite loop after OAuth redirect.
- [ ] `getRedirectResult()` returns a valid user credential (not `null`).
- [ ] Auth state updates correctly and user is redirected to dashboard.
- [ ] No regression on Safari or other browsers.

## Notes
- The production environment likely has the same issue; fix should be applied there after verification in staging.
- Consider logging the full redirect event structure for deeper debugging.
