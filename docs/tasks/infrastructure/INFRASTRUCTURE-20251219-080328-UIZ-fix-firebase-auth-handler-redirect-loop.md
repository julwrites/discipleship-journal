---
id: INFRASTRUCTURE-20251219-080328-UIZ
status: completed
title: Fix Firebase Auth Handler Redirect Loop
priority: medium
created: 2025-12-19 08:03:28
category: infrastructure
dependencies:
type: task
---

# Fix Firebase Auth Handler Redirect Loop

## Background
Users were experiencing a redirect loop or authentication failure when using `signInWithRedirect`. This often occurs in PWA applications when the Service Worker intercepts the `/__/auth/handler` request (which is used by Firebase Auth for OAuth redirects) and serves the `index.html` fallback instead of allowing the request to hit the Firebase Hosting server.

## Issue Analysis
The `vite-plugin-pwa` configuration uses a navigation fallback to support SPA routing (serving `index.html` for unknown routes). However, without an explicit denylist (blacklist), this fallback also captures the reserved `/__/` paths used by Firebase.

## Solution
Updated `web/vite.config.ts` to include `navigateFallbackDenylist: [/^\/__/]` in the `workbox` configuration.

```typescript
// web/vite.config.ts
VitePWA({
  // ...
  workbox: {
    navigateFallbackDenylist: [/^\/__/],
  },
  // ...
})
```

This ensures that any request starting with `/__/` (including `/__/auth/handler`) bypasses the Service Worker and is handled by the server (Firebase Hosting), allowing the OAuth flow to complete successfully.

## Verification
- Verified `web/dist/sw.js` contains the denylist configuration.
- Verified unit tests for Login page pass (`npm test`).
- Verified build succeeds (`npm run build`).

## Status
Completed.
