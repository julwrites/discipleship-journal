---
id: INFRASTRUCTURE-20251223-033044-KMQ
status: completed
title: Clean up Vite configuration
priority: low
created: 2025-12-23 03:30:44
category: infrastructure
dependencies:
type: chore
---

# Clean up Vite configuration

## Description
In `web/vite.config.ts`, the configuration object is cast to `any` at the end:
```typescript
} as any)
```
This defeats type safety for the configuration.

## Acceptance Criteria
- [x] Remove `as any` casting.
- [x] Fix any type errors that arise (likely related to `Vitest` types or `VitePWA` plugin types).
- [x] Ensure `npm run build` and `npm run test` still work.
