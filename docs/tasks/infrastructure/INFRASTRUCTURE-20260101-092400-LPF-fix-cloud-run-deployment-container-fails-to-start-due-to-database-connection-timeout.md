---
id: INFRASTRUCTURE-20260101-092400-LPF
status: completed
title: Fix Cloud Run deployment: Container fails to start due to database connection timeout
priority: medium
created: 2026-01-01 09:24:00
category: infrastructure
dependencies: 
type: task
---

# Fix Cloud Run deployment: Container fails to start due to database connection timeout

## Problem
Cloud Run deployments were failing with error:
```
ERROR: (gcloud.run.deploy) The user-provided container failed to start and listen on the port defined provided by the PORT=8080 environment variable within the allocated timeout.
```

The root cause was that in production mode (`APP_ENV=production`), the application would exit immediately if the database connection failed (`main.go:84-87`). Cloud Run expects containers to start listening on the port within a timeout, and if the application exits before binding to the port, the deployment fails.

## Root Cause Analysis
1. **Database dependency at startup**: The application tried to connect to the database before starting the HTTP server
2. **Immediate exit in production**: If `database.Connect()` failed in production mode, the application would `os.Exit(1)`
3. **Cloud Run timeout**: Cloud Run has a startup timeout; if the container doesn't bind to the port within this time, deployment fails

## Solution
Modified the startup sequence in `api/main.go`:

1. **Don't exit on database failure in production**: Changed the logic to continue starting the server even if database connection fails
2. **Track database connection state**: Added `dbConnected` boolean to track connection status
3. **Enhanced health check**: Updated `/health` endpoint to return appropriate status based on database connectivity:
   - `503 DATABASE_CONNECTION_FAILED` if database never connected
   - `503 DATABASE_PING_FAILED` if connection lost
   - `200 OK` if database is healthy

## Changes Made

### `api/main.go:81-94`
```go
// Try to connect to database, but don't exit immediately in production
// Cloud Run needs the container to start listening on the port first
dbConnected := false
if err := database.Connect(); err != nil {
    logger.Error("Database connection failed", "error", err)
    // Don't exit immediately - let the server start and health check will fail
    // This allows Cloud Run to properly start the container
    if os.Getenv("APP_ENV") == "production" {
        logger.Warn("Database connection failed in production, but continuing to start server")
    }
} else {
    dbConnected = true
    defer database.Close()
}
```

### `api/main.go:226-249`
```go
r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    if !dbConnected {
        w.WriteHeader(http.StatusServiceUnavailable)
        logger.Error("Health check failed: DB not connected")
        if _, err := w.Write([]byte("DATABASE_CONNECTION_FAILED")); err != nil {
            logger.Error("Failed to write health check response", "error", err)
        }
        return
    }

    if err := database.DB.Ping(r.Context()); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        logger.Error("Health check failed: DB ping failed", "error", err)
        if _, err := w.Write([]byte("DATABASE_PING_FAILED")); err != nil {
            logger.Error("Failed to write health check response", "error", err)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
    if _, err := w.Write([]byte("OK")); err != nil {
        logger.Error("Failed to write health check response", "error", err)
    }
})
```

## Testing
1. **Local test with database**: Container starts successfully, health endpoint returns `200 OK`
2. **Local test without database**: Container starts successfully, health endpoint returns `503 DATABASE_CONNECTION_FAILED`
3. **Both scenarios**: Container binds to port 8080 within Cloud Run timeout requirements

## Impact
- **Positive**: Cloud Run deployments will now succeed even if database connection fails initially
- **Positive**: Health checks accurately reflect database status
- **Consideration**: API endpoints that require database access will fail gracefully (return 500 errors) when database is unavailable
- **Next step**: Consider implementing retry logic for database connections in background

## Verification
- [x] Container starts successfully with valid database connection
- [x] Container starts successfully without database connection (production mode)
- [x] Health endpoint returns appropriate status codes
- [x] Server binds to port within expected timeframe
