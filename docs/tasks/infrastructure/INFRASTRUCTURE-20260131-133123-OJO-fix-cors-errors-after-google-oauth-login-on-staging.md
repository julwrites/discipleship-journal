---
id: INFRASTRUCTURE-20260131-133123-OJO
status: completed
title: Fix CORS errors after Google OAuth login on staging
priority: medium
created: 2026-01-31 13:31:23
category: infrastructure
dependencies: 
type: task
---

# Fix CORS errors after Google OAuth login on staging

## Problem
After successful Google OAuth login via popup flow on Edge/Chrome, the frontend receives CORS errors when trying to call backend API endpoints.

## Symptoms from Logs
1. Google OAuth popup login succeeds: `"Popup login successful: jump4jupiter@gmail.com Provider: google.com"`
2. Auth state updates correctly: `"Auth state changed, user: jump4jupiter@gmail.com"`
3. Frontend attempts to fetch user profile: `POST https://discipleship-journal-api-staging-995319008345.asia-southeast1.run.app/api/users/me`
4. CORS error: `Access to fetch at 'https://discipleship-journal-api-staging-995319008345.asia-southeast1.run.app/api/users/me' from origin 'https://discipleship-journal-staging.firebaseapp.com' has been blocked by CORS policy: Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource.`
5. Similar error for `GET /api/notes` endpoint.

## Root Cause Analysis
The backend API CORS middleware (`github.com/go-chi/cors`) validates origins using the `CORS_ALLOWED_ORIGINS` environment variable loaded from Google Secret Manager. The `AllowOriginFunc` returns `false` for origins not in the allowed list, causing preflight OPTIONS requests to fail with missing CORS headers.

Current configuration in `api/main.go`:
```go
AllowOriginFunc: func(r *http.Request, origin string) bool {
    // Check env var first for explicit allowed origins
    if allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); allowedOrigins != "" {
        for _, allowed := range strings.Split(allowedOrigins, ",") {
            if strings.TrimSpace(allowed) == origin {
                return true
            }
        }
    }
    return false
},
```

The staging frontend origin `https://discipleship-journal-staging.firebaseapp.com` is likely not included in the `CORS_ALLOWED_ORIGINS` secret.

## Investigation Findings
1. **CORS Configuration**: Backend uses chi/cors middleware with dynamic origin validation
2. **Secret Management**: `CORS_ALLOWED_ORIGINS` should be stored in Google Secret Manager as per `docs/setup/secrets.md`
3. **Deployment**: Staging backend deployment (`deploy-staging.yml`) does not explicitly set `CORS_ALLOWED_ORIGINS` as environment variable; relies on Secret Manager
4. **Related Task**: `FEATURES-20251220-093910-PGM-fix-cors-and-bible-api-integration.md` (completed) addressed similar CORS issues previously

## Required Actions
1. **Update Google Secret Manager**: Add staging frontend origin to `CORS_ALLOWED_ORIGINS` secret
   - Current value unknown; need to check existing secret
   - Should include: `https://discipleship-journal-staging.firebaseapp.com`
   - Consider also including production origin if not already present
   - Local development origins: `http://localhost:5173`, `http://localhost:8081`
2. **Verify Secret Access**: Ensure backend service account has `roles/secretmanager.secretAccessor` permission
3. **Test Configuration**: After updating secret, restart backend or wait for next deployment
4. **Documentation**: Update any configuration documentation

## Implementation Steps
1. Check current `CORS_ALLOWED_ORIGINS` secret value:
   ```bash
   gcloud secrets versions access latest --secret=CORS_ALLOWED_ORIGINS --project=discipleship-journal-pwa
   ```
2. Update secret with staging origin (append if comma-separated list):
   ```bash
   echo "https://discipleship-journal-staging.firebaseapp.com,http://localhost:5173,http://localhost:8081" | gcloud secrets versions add CORS_ALLOWED_ORIGINS --data-file=- --project=discipleship-journal-pwa
   ```
   Or update existing list while preserving other origins.
3. Restart Cloud Run service to pick up new secret:
   ```bash
   gcloud run services update discipleship-journal-api --region=asia-southeast1 --project=discipleship-journal-pwa
   ```
4. Test by logging in again and checking browser console for CORS errors.

## Acceptance Criteria
- [x] CORS preflight requests succeed for staging frontend origin
- [x] API endpoints (`/api/users/me`, `/api/notes`, etc.) return proper CORS headers
- [x] No CORS errors in browser console after authentication
- [x] User profile and notes load successfully after login

## Implementation Details
- Modified `api/main.go` to support `CORS_EXTRA_ORIGINS` environment variable. This allows appending additional allowed origins to those loaded from Secret Manager (which normally overrides environment variables).
- Updated `.github/workflows/deploy-staging.yml` to inject `CORS_EXTRA_ORIGINS` with the staging frontend URL (`https://discipleship-journal-staging.firebaseapp.com`).
- This approach avoids the need to manually update Google Secret Manager secrets for environment-specific configurations and allows Staging to be configured via Infrastructure-as-Code.

## Notes
- This is a staging-specific issue but production may have similar configuration
- COOP warnings about `window.closed` calls are separate and non-critical (Firebase SDK internal)
- Consider adding local development origins (`http://localhost:5173`, `http://localhost:8081`) to secret for development convenience
