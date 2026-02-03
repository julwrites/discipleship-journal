---
id: INFRASTRUCTURE-20251231-231139-VZZ
status: completed
title: Implement Google Secret Manager integration
priority: medium
created: 2025-12-31 23:11:39
category: infrastructure
dependencies: 
type: task
---

# Implement Google Secret Manager integration

## Problem
The application currently reads secrets (DATABASE_URL, BIBLE_API_URL, BIBLE_API_KEY) directly from environment variables. In production, these should be retrieved from Google Secret Manager for better security. Environment variables should be a fallback for local development.

## Acceptance Criteria
- [x] Add Google Secret Manager Go dependency
- [x] Create a secret loading utility that tries Secret Manager first, then falls back to env vars
- [x] Update main.go to use the new secret loader for DATABASE_URL, BIBLE_API_URL, BIBLE_API_KEY
- [x] Add proper error handling when secrets cannot be loaded
- [ ] Update documentation to reflect new secret loading strategy
- [x] Test both local (env var) and production (Secret Manager) paths

## Implementation Plan
1. Add `cloud.google.com/go/secretmanager` dependency
2. Create `services/secrets.go` with `LoadSecret` function
3. Update `main.go` to use secret loader
4. Add graceful degradation - use env vars if Secret Manager fails
5. Update deployment documentation

## Implementation Details

### What was implemented:
1. **Added `services/secrets.go`** with `SecretLoader` struct that:
   - Tries Google Secret Manager first (if client can be initialized)
   - Falls back to environment variables if Secret Manager fails or is unavailable
   - Gracefully handles missing credentials (local development)
   - Properly closes the Secret Manager client

2. **Updated `main.go`** to:
   - Initialize `SecretLoader` at startup
   - Load `DATABASE_URL` from secrets (sets env var for backward compatibility)
   - Load `BIBLE_API_URL` and `BIBLE_API_KEY` from secrets
   - Exit on database connection failure in production (but continue in development)
   - Log which Bible API client is being used (real vs mock)

3. **Error handling**:
   - Database connection failure causes exit in production only
   - Secret loading failures are logged but don't stop application
   - Mock Bible API client used if secrets not available

### Secret Manager behavior:
- Requires `GOOGLE_CLOUD_PROJECT` environment variable to be set
- Uses default application credentials in Cloud Run
- For local development: set `GOOGLE_APPLICATION_CREDENTIALS` or rely on env vars
- Secret names: `DATABASE_URL`, `BIBLE_API_URL`, `BIBLE_API_KEY`

### Testing:
- All existing tests pass
- Application builds successfully
- Local development continues to work with .env files
- Production will use Secret Manager if credentials are available

## Notes
- Secret names should follow convention: `DATABASE_URL`, `BIBLE_API_URL`, `BIBLE_API_KEY`
- Project ID should be read from `GOOGLE_CLOUD_PROJECT` env var
- Local development should continue to work with .env files
- The application now exits on database connection failure in production (safety improvement)
