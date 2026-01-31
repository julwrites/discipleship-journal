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

## Test Results Summary

### SessionStorage Cleanup (Test 1)
- ✅ **Working**: Successfully detects and clears stale Firebase redirect events from `sessionStorage`
- ❌ **Issue persists**: `getRedirectResult()` returns `null` even with fresh redirect event in `sessionStorage`
- 🔍 **Observation**: Redirect event has `eventId: null`, Firebase writes but cannot read/process it

### Service Worker Unregistration (Test 2)
- ✅ **Working**: Service worker unregistration successful ("Found 1 service worker registration(s)", "Service worker unregistration successful")
- ❌ **Issue persists**: `getRedirectResult()` still returns `null` after service worker removal
- 🔍 **Conclusion**: Service workers are NOT the root cause

### COOP Header Removal (Test 3)
- ✅ **Tested**: Removed `Cross-Origin-Opener-Policy: same-origin-allow-popups` header from `firebase.json`
- ❌ **Issue persists**: `getRedirectResult()` still returns `null` without COOP header
- 🔍 **Conclusion**: COOP header is NOT the root cause

### Popup Flow Implementation (Test 4 - SUCCESS!)
- ✅ **Implemented**: Popup-first authentication flow in `Login.tsx`
- ✅ **Result**: **POPUP FLOW WORKS** on Edge/Chrome! User successfully authenticated: `jump4jupiter@gmail.com`
- ✅ **Key success logs**: `"Popup login successful: jump4jupiter@gmail.com Provider: google.com"`, `"Auth state changed, user: jump4jupiter@gmail.com"`
- 🔍 **Observations**:
  - Popup appears and login completes successfully
  - Redirect flow still fails (`getRedirectResult()` returns `null`)
  - Service worker unregistration causes error: `InvalidStateError: Only the active worker can claim clients.`
  - CORS errors appear after authentication (separate backend issue)
- ✅ **Conclusion**: **Popup flow successfully bypasses Edge/Chrome redirect flow compatibility issue**

### Root Issue Identified
Firebase successfully creates redirect event with `eventId: null` in `sessionStorage` on Edge/Chrome, but `getRedirectResult()` cannot read/process it. This appears to be a Firebase SDK compatibility issue with Edge/Chrome's implementation of OAuth redirect flow and sessionStorage isolation.

## Root Cause Update
Based on test results, the primary issue appears to be:
1. **Firebase SDK/Edge-Chrome compatibility**: Firebase successfully writes redirect event to `sessionStorage` but fails to read/process it
2. **Possible `eventId: null` issue**: The redirect event has `eventId: null` which might indicate Google OAuth isn't providing proper event ID on Edge/Chrome
3. **Storage isolation**: Browser security policies may isolate `sessionStorage` across redirects despite event being present

## Proposed Solutions - Implementation Status

### 1. Clear Stale Redirect Events (IMPLEMENTED & TESTED)
✅ **Implemented**: Automatic cleanup of stale Firebase redirect events in `Login.tsx`
❌ **Result**: Helps clean up but doesn't fix root issue - `getRedirectResult()` still returns `null`

### 2. Disable Service Worker for Testing (IMPLEMENTED & TESTED)
✅ **Implemented**: Service worker unregistration code in `Login.tsx`
❌ **Result**: Service workers successfully unregistered but `getRedirectResult()` still returns `null`
🔍 **Conclusion**: Service workers are NOT the root cause

### 3. Adjust Cross-Origin Headers (IMPLEMENTED & TESTED)
✅ **Tested**: Temporarily removed `Cross-Origin-Opener-Policy` header from `firebase.json`
❌ **Result**: `getRedirectResult()` still returns `null` without COOP header
🔍 **Conclusion**: COOP header is NOT the root cause for redirect flow
🔄 **Update**: COOP header restored to `same-origin-allow-popups` for popup flow compatibility

### 4. Use Popup Flow Instead of Redirect (IMPLEMENTED & TESTED - SUCCESS!)
✅ **Implemented**: Popup-first authentication flow in `Login.tsx`
✅ **Result**: **POPUP FLOW WORKS SUCCESSFULLY** on Edge/Chrome
🔍 **Test Results**:
  - User successfully authenticated via popup: `jump4jupiter@gmail.com`
  - Auth state updated correctly: `"Auth state changed, user: jump4jupiter@gmail.com"`
  - Popup flow bypasses Edge/Chrome redirect compatibility issue
⚠️ **Side Effects**:
  - Service worker unregistration causes error (needs adjustment)
  - CORS errors after authentication (separate backend issue)
📝 **Note**: COOP header `same-origin-allow-popups` required for popup compatibility

### 5. Cookie Configuration
Ensure Firebase OAuth client is configured with correct authorized domains and that cookies are allowed with appropriate SameSite settings.

## Implementation Plan - Current Status
1. ✅ **Add sessionStorage cleanup** in `Login.tsx` after `getRedirectResult()` returns `null`. (COMPLETED & TESTED)
2. ✅ **Deploy to staging** and test on Edge/Chrome. (COMPLETED - issue persists)
3. ✅ **Disable service worker** and retest. (COMPLETED & TESTED - not the root cause)
4. ✅ **Adjust headers** (temporarily remove COOP). (COMPLETED & TESTED - not the root cause, header restored for popups)
5. 🔄 **Implement popup fallback** with user option. (IN PROGRESS - popup-first flow implemented)
6. **Update Firebase console** OAuth configuration to ensure authorized domains include staging URL. (PENDING - if popup also fails)

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
