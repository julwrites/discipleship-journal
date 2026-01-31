---
id: FEATURES-20260131-162601-LPD
status: pending
title: Implement user-specific Bible version preferences with default fallback for memory verses
priority: medium
created: 2026-01-31 16:26:01
category: features
dependencies: 
type: task
---

# Implement user-specific Bible version preferences with default fallback for memory verses

## Task Information
- **Dependencies**: None (alternative approach to FEATURES-20260131-162336-UYJ)
- **Related tasks**:
  - FEATURES-20260131-162336-UYJ (minimal version editing solution)
  - FEATURES-20260131-162514-WAF (default version on clone)

## Task Details
Implement a comprehensive solution for Bible version preferences where users can set version overrides for individual verses, with automatic fallback to their default version when no override exists. This preserves original system pack data while allowing personalized version preferences.

### Current Limitations:
- Verse versions are stored directly in `memory_verses.version`
- System pack verses have fixed 'ESV' versions
- Users cannot personalize versions without modifying original data
- No distinction between "user override" vs "inherited default"

### Proposed Solution (Option B - User-specific overrides):
1. **New table**: `user_verse_preferences`
   - `user_id` (UUID, references `users.id`)
   - `verse_id` (UUID, references `memory_verses.id`)
   - `version_override` (VARCHAR(50), nullable)
   - `created_at`, `updated_at` (timestamps)
   - Primary key: `(user_id, verse_id)`

2. **Version resolution logic**:
   - Check `user_verse_preferences` for user-specific override
   - If override exists, use it
   - Otherwise, use user's default version from `user.settings.bible_version`
   - Fallback to verse's original `version` if no user default (edge case)

3. **Benefits**:
   - Preserves original system pack data
   - Allows per-verse user customization
   - Clean separation of concerns
   - Supports "reset to default" functionality

### Steps:

#### 1. Database Migration:
- Create `user_verse_preferences` table
- Add foreign key constraints
- Consider indexes on `(user_id, verse_id)` for performance

#### 2. Backend Changes:
- **New service**: `VersePreferenceService` with CRUD operations
- **Update `MemoryVerseService`**: Modify version resolution logic in `GetVerses()` and related methods
- **API endpoints**:
  - `PUT /api/memory-verses/{id}/version` to set user override
  - `DELETE /api/memory-verses/{id}/version` to clear override
- **Authorization**: Ensure users can only manage their own preferences

#### 3. Frontend Changes:
- **Version selector**: Update to save/clear user preferences
- **Verse display**: Show "Your version: NIV" vs "Default: ESV" indicators
- **Settings integration**: Link to user's default version settings
- **Bulk operations**: "Apply my default version to all verses in this pack"

#### 4. Data Migration (Optional):
- Migrate existing user-modified verses to preference table
- Or keep existing `memory_verses.version` for backward compatibility

#### 5. Testing:
- Test version resolution: override → user default → original
- Test API endpoints for preferences
- Verify authorization prevents accessing others' preferences
- Test bulk operations and performance with many verses

### Acceptance Criteria
- [ ] `user_verse_preferences` table created with proper constraints
- [ ] Version resolution logic implemented (override → user default → original)
- [ ] API endpoints for managing version preferences
- [ ] Frontend allows setting/clearing per-verse version overrides
- [ ] Clear visual indication of "your version" vs "default"
- [ ] Bulk operation: "Apply my default version to all verses in pack"
- [ ] Backward compatibility: existing verses display correctly
- [ ] Performance: version resolution doesn't significantly impact load times

## Implementation Status
### Completed Work
- ✅ Architectural analysis completed
- ✅ Database design proposed

### Blockers
None yet.

## Notes
- **Complexity**: Higher than minimal solution (Option A) but more robust
- **Performance**: Consider eager loading preferences for user's verses
- **Migration Strategy**:
  - Option 1: Migrate existing modified verses to preferences table
  - Option 2: Dual system: check `memory_verses.version` if user owns verse, else use preferences
- **UI/UX**: Need clear indicators for version source (override/default/original)
- **Future Extensions**: Could extend to other preferences (font size, highlighting, etc.)
