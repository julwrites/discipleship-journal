# API Testing Strategy

## Overview
This document outlines the testing strategy for the Discipleship Journal API endpoints to ensure reliability before deployment.

## Critical Issues Fixed

### 1. **Authentication Middleware Failure**
**Problem**: When Firebase initialization fails (no credentials in local development), `authMiddleware` was nil, causing:
- Protected routes registered without auth middleware
- Handlers panicked when trying to access auth token from context

**Solution**:
- Created mock auth middleware for development
- Injects `TestUserKey` into context (matches `GetUserUUID` helper)
- Accepts `Authorization: Bearer test` header for development

### 2. **Group Handler Panic**
**Problem**: `group.go:116` attempted to cast context value to `*auth.Token` without type assertion check.

**Solution**: Updated handler to:
1. Check for `TestUserKey` first (for development/testing)
2. Use type assertion for `*auth.Token` (for production)
3. Return 401 if neither exists

### 3. **Bible API 500 Error**
**Problem**: `.env` had `BIBLE_API_URL` set but invalid `BIBLE_API_KEY`, causing real API calls to fail.

**Solution**:
- Clear `BIBLE_API_URL` in `.env` to use mock client
- For production: set valid credentials in Cloud Run environment variables

## API Endpoint Summary

**Total Endpoints**: 36
- **Public**: 2 (`/health`, `/swagger/*`)
- **Protected**: 34 (require authentication)

**Test Coverage**: 94.4% (34/36 endpoints have unit tests)
- Only public endpoints lack explicit tests (acceptable)

## Testing Methodology

### 1. **Local Development Testing**
```bash
# Start services
docker-compose up --build -d

# Run comprehensive API test
./scripts/test_apis.sh

# Or test manually:
curl -H "Authorization: Bearer test" http://localhost:8081/api/notes
```

### 2. **Unit Tests**
```bash
cd api
go test ./handlers -v      # Handler tests
go test ./tests/... -v     # Contract/integration tests
```

### 3. **Pre-Deployment Checklist**
- [ ] All unit tests pass
- [ ] API test script runs successfully
- [ ] No 500 errors in basic endpoint tests
- [ ] Authentication works (401 without auth, 200 with `Bearer test`)
- [ ] Critical endpoints tested:
  - `GET /api/notes` - Notes listing
  - `POST /api/notes` - Note creation
  - `GET /api/bible/passage` - Bible API
  - `GET /api/groups` - Groups listing
  - `GET /api/reading-plans` - Reading plans

## Environment Configuration

### Development (.env)
```env
# Use mock Bible API
BIBLE_API_URL=
BIBLE_API_KEY=

# Use mock authentication
# (automatically falls back when Firebase fails)
```

### Production (Cloud Run)
```env
# Real Firebase credentials
GOOGLE_APPLICATION_CREDENTIALS=auto

# Real Bible API (if configured)
BIBLE_API_URL=https://api.scripture.api.bible
BIBLE_API_KEY=your_actual_key_here
```

## Common Issues & Solutions

### 1. **500 Errors in Production**
**Check**:
- Firebase service account credentials
- Database connection string
- Bible API credentials (if using real API)
- Environment variables in Cloud Run

### 2. **401 Errors with Valid Token**
**Check**:
- User exists in database (first-time users need `POST /api/users/me`)
- Firebase UID matches database record

### 3. **Panic Stack Traces**
**Check**:
- Nil pointer dereferences in handlers
- Type assertions without checks (`.(*auth.Token)` vs `, ok`)
- Database connection issues

## Automated Testing Script

See `scripts/test_apis.sh` for comprehensive endpoint testing.

## Future Improvements

1. **Integration Test Suite**: End-to-end tests with real database
2. **Load Testing**: Performance testing for critical endpoints
3. **Monitoring**: Alert on 500 errors in production
4. **API Contract Tests**: Verify OpenAPI specification matches implementation

## Emergency Rollback

If APIs fail in production:
1. Check Cloud Run logs for panic stack traces
2. Rollback to previous working version
3. Enable maintenance mode if needed
4. Fix root cause in development, test thoroughly, redeploy