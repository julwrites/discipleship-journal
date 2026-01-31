---
id: FOUNDATION-20260131-161127-UQK
status: pending
title: Optimize database queries for notes loading: add missing indexes and improve COUNT performance
priority: medium
created: 2026-01-31 16:11:27
category: foundation
dependencies: 
type: task
---

# Optimize database queries for notes loading: add missing indexes and improve COUNT performance

## Task Information
- **Dependencies**: FOUNDATION-20260131-161114-QCR (performance analysis should identify specific indexes needed)
- **Related tasks**: PRESENTATION-20260131-150229-RVN (parent epic)

## Task Details
Implement database optimizations for notes loading based on performance analysis findings. Add missing indexes, optimize COUNT queries, and improve query performance for the `/api/notes` endpoint.

### Current findings from exploration:
- Existing indexes: `idx_notes_user_id` (user_id), `idx_notes_content` (GIN on content JSONB)
- Missing composite indexes for common query patterns:
  - `(user_id, deleted_at, updated_at)` for default ordering and filtering
  - `(user_id, deleted_at, title)` for search queries
  - `(updated_at)` for date range queries
- `COUNT(*)` may perform full table scans; consider approximate counts or materialized views
- Query filters: `WHERE user_id = ? AND deleted_at IS NULL ORDER BY updated_at DESC`

### Steps:
1. **Review performance analysis findings** from FOUNDATION-20260131-161114-QCR
2. **Create migration for new indexes**:
   - Composite index on `(user_id, deleted_at, updated_at)` for default queries
   - Index on `title` for search performance (maybe using pg_trgm for ILIKE)
   - Consider partial index `WHERE deleted_at IS NULL` if most notes are not deleted
3. **Optimize COUNT query**:
   - Evaluate `COUNT(*)` vs `COUNT(id)` performance
   - Consider using `EXPLAIN` to estimate row counts for pagination
   - Implement cached count in separate table if needed
4. **Test query performance**:
   - Run benchmarks before and after index creation
   - Verify query plans use new indexes
   - Ensure no negative impact on write performance
5. **Update database schema documentation** with new indexes and their purposes

### Acceptance Criteria
- [ ] New database indexes created via migration
- [ ] `EXPLAIN ANALYZE` shows improved query plans (index scans vs sequential scans)
- [ ] COUNT query performance improved or optimized
- [ ] No regression in write performance (insert/update/delete)
- [ ] Database schema documentation updated

## Implementation Status
### Completed Work
- ✅ Initial analysis of missing indexes completed

### Blockers
- Waiting on performance analysis to confirm specific indexes needed

## Notes
- Coordinate with infrastructure team for production database migration
- Consider downtime implications for index creation on large tables
- Test with representative data volume
