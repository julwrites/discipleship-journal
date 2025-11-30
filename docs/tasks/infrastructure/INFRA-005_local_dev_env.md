# Task: Local Development Environment

## Task Information
- **Task ID**: INFRA-005
- **Status**: completed
- **Priority**: high
- **Phase**: 1.5
- **Estimated Effort**: 0.5 days
- **Dependencies**: FOUNDATION-001

## Task Details

### Description
Setup Docker Compose for local development and testing. This allows running the full stack (DB, API, Web) locally in a containerized environment.

### Acceptance Criteria
- [x] **Docker**:
    - [x] Create `web/Dockerfile` for production-like build.
    - [x] Create `docker-compose.yml` defining services: `db`, `api`, `web`, `migrator`.
    - [x] Included `migrator` service to run migrations.
- [x] **Configuration**:
    - [x] Ensure `api` can connect to `db` (`DATABASE_URL`).
    - [x] Ensure `web` can connect to `api` (via nginx proxy).
- [x] **Documentation**:
    - [x] Update `README.md` with "How to run locally" instructions.

---

*Created: 2025-05-18*
*Status: completed*
