# Task: Setup Database Migrations

## Task Information
- **Task ID**: FOUNDATION-003
- **Status**: pending
- **Priority**: high
- **Phase**: 1.5
- **Estimated Effort**: 0.5 days
- **Dependencies**: FOUNDATION-001

## Task Details

### Description
Establish a robust system for managing database schema changes using version control. This ensures all developers and environments are in sync.

### Acceptance Criteria
- [ ] **Tool Selection**:
    - [ ] Install `golang-migrate/migrate` CLI tool (or include instructions).
- [ ] **Migration Setup**:
    - [ ] Create `api/migrations` directory.
    - [ ] Create initial migration (`000001_init_schema.up.sql` / `.down.sql`).
    - [ ] Define Users table and Journal Notes table structure (relational + JSONB).
- [ ] **Execution**:
    - [ ] Create a `Makefile` or script to run `migrate up` and `migrate down`.
    - [ ] Document how to run migrations locally.

### Implementation Notes
- Migration files should be timestamped or sequentially numbered.
- The `up` script applies changes; the `down` script reverts them.
- Ensure the app can connect to the DB with the correct credentials.

---

*Created: 2025-05-18*
*Status: pending*
