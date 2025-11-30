# Task: Setup Database Migrations

## Task Information
- **Task ID**: FOUNDATION-003
- **Status**: completed
- **Priority**: high
- **Phase**: 1.5
- **Estimated Effort**: 0.5 days
- **Dependencies**: FOUNDATION-001

## Task Details

### Description
Establish a robust system for managing database schema changes using version control. This ensures all developers and environments are in sync.

### Acceptance Criteria
- [x] **Tool Selection**:
    - [x] Install `golang-migrate/migrate` CLI tool (or include instructions).
- [x] **Migration Setup**:
    - [x] Create `api/migrations` directory.
    - [x] Create initial migration (`000001_init_schema.up.sql` / `.down.sql`).
    - [x] Define Users table and Journal Notes table structure (relational + JSONB).
- [x] **Execution**:
    - [x] Create a `Makefile` or script to run `migrate up` and `migrate down`.
    - [x] Document how to run migrations locally.

### Implementation Notes
- Migration files should be timestamped or sequentially numbered.
- The `up` script applies changes; the `down` script reverts them.
- Ensure the app can connect to the DB with the correct credentials.

### Instructions
To run migrations locally:
1. Ensure Postgres is running.
2. Install migrate: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
3. Run: `./scripts/migrate_up.sh` (Set `DB_URL` env var if different from default).

---

*Created: 2025-05-18*
*Status: completed*
