# Task: Setup Observability (Logging & Metrics)

## Task Information
- **Task ID**: INFRA-002
- **Status**: pending
- **Priority**: medium
- **Phase**: 1.5
- **Estimated Effort**: 1 day
- **Dependencies**: FOUNDATION-001

## Task Details

### Description
Implement structured logging and basic error tracking to enable effective debugging and monitoring of the application in production.

### Acceptance Criteria
- [ ] **Backend Logging**:
    - [ ] Configure `log/slog` to output JSON format in production.
    - [ ] Add Request ID middleware to correlate logs.
    - [ ] Log HTTP requests (Method, Path, Status, Duration).
- [ ] **Frontend Logging**:
    - [ ] Implement a global error boundary to catch React errors.
    - [ ] (Optional) Setup a service like Sentry (mockable for now).
- [ ] **Health Checks**:
    - [ ] Implement `/health` endpoint on Backend.
    - [ ] Include DB connection check in health response.

### Implementation Notes
- Use standard `log/slog` in Go.
- Ensure sensitive data (passwords, tokens) is never logged.

---

*Created: 2025-05-18*
*Status: pending*
