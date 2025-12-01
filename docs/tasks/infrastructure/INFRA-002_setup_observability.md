# Task: Setup Observability (Logging & Metrics)

## Task Information
- **Task ID**: INFRA-002
- **Status**: completed
- **Priority**: medium
- **Phase**: 1.5
- **Estimated Effort**: 1 day
- **Dependencies**: FOUNDATION-001

## Task Details

### Description
Implement structured logging and basic error tracking to enable effective debugging and monitoring of the application in production.

### Acceptance Criteria
- [x] **Backend Logging**:
    - [x] Configure `log/slog` to output JSON format in production.
    - [x] Add Request ID middleware to correlate logs.
    - [x] Log HTTP requests (Method, Path, Status, Duration).
- [x] **Frontend Logging**:
    - [x] Implement a global error boundary to catch React errors.
    - [ ] (Optional) Setup a service like Sentry (mockable for now).
- [x] **Health Checks**:
    - [x] Implement `/health` endpoint on Backend.
    - [x] Include DB connection check in health response.

### Implementation Notes
- Used standard `log/slog` in Go.
- `api/middleware/logging.go` logs request details.
- `web/src/components/ErrorBoundary.tsx` catches React errors.
- `api/main.go` sets up the logger and health check.
- Sentry setup is deferred/optional, currently logging to console in ErrorBoundary.

---

*Created: 2025-05-18*
*Status: completed*
