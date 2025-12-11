---
id: INFRA-005
status: completed
title: Local Development Environment
priority: high
created: 2025-12-11 06:09:10
category: unknown
type: task
---

# Local Development Environment

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
