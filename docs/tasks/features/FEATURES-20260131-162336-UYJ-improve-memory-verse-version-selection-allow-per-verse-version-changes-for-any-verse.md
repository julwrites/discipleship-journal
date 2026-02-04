---
id: FEATURES-20260131-162336-UYJ
status: completed
title: Improve memory verse version selection: allow per-verse version changes for any verse
priority: medium
created: 2026-01-31 16:23:36
category: features
dependencies: 
type: task
---

# Improve memory verse version selection: allow per-verse version changes for any verse

## Task Information
- **Dependencies**: None
- **Related tasks**: None identified

## Task Details
Currently, users can only change Bible versions for memory verses in packs they own. For system/public packs (created by admin), all verses have fixed 'ESV' versions and users cannot modify them. This task enables per-verse version selection for any verse, including those in public/system packs, while maintaining appropriate access controls.

### Current Implementation Analysis:
- **Database**: `memory_verses` table has `version` field (default 'ESV')
- **Backend**: `UpdateVerse()` method in `memory_verse_service.go` restricts updates to pack owners only
- **Frontend**: MemoryVersesPage only shows edit/delete buttons for user-owned packs (`isMyPack`)
- **User Settings**: Default Bible version stored in `user.settings.bible_version`

### Problem:
Users want to customize Bible versions for individual verses in public/system packs without cloning entire packs. Currently, they must clone a pack to modify verses, creating duplicates.

### Proposed Solution (Option A - Minimal):
Allow version updates on public pack verses while keeping other modifications (delete, title, reference) restricted to pack owners.

### Steps:

#### 1. Backend Changes (`api/services/memory_verse_service.go`):
- Modify `UpdateVerse()` method authorization logic
- Current: `AND vp.user_id = $6` (line ~250)
- New: `AND (vp.user_id = $6 OR vp.is_public = true)`
- Ensure only version field can be updated for public packs (optional validation)

#### 2. Frontend Changes (`web/src/pages/MemoryVersesPage.tsx`):
- Update `VersePackDetail` component edit button logic
- Current: Edit button only shown for `isMyPack` (line ~563)
- New: Show edit button for all verses, but disable delete for non-owned packs
- Consider visual distinction: version-only edit vs full edit capabilities

#### 3. UI/UX Considerations:
- Add tooltip explaining "Can only change version for this public verse"
- Keep delete button restricted to pack owners
- Ensure version selector component (`BibleVersionSelector`) works for all verses

#### 4. Testing:
- Test version update on user-owned pack verses (existing functionality)
- Test version update on public/system pack verses (new functionality)
- Verify delete still restricted to pack owners
- Test error handling for unauthorized modifications

### Acceptance Criteria
- [ ] Users can change Bible version for any verse (including public/system packs)
- [ ] Version changes persist correctly in database
- [ ] Delete functionality remains restricted to pack owners
- [ ] Edit UI accessible for all verses (not just user-owned)
- [ ] User's default version still pre-loads when adding new verses
- [ ] No regression in existing memory verse functionality
- [ ] Backward compatible: existing verses maintain their versions

## Implementation Status
### Completed Work
- ✅ Codebase exploration and analysis completed
- ✅ Current limitations identified
- ✅ Proposed solution designed

### Blockers
None yet.

## Notes
- **Security**: Ensure authorization logic correctly distinguishes between version updates vs other modifications
- **UX**: Consider adding visual indicator for "public verse - version only editable"
- **Performance**: No significant performance impact expected
- **Database**: No schema changes needed - `version` field already exists
- **Alternative approaches considered**:
  - Option B (User-specific overrides table): More complex, preserves original system data (see FEATURES-20260131-162601-LPD)
  - Option C (Clone-on-edit): Creates duplicates, simpler but less elegant
- **Recommended**: Option A as it's minimal, uses existing infrastructure, and meets user requirements
- **Related implementation tasks**:
  - FEATURES-20260131-162514-WAF: Apply user's default Bible version when cloning memory verse packs
  - FEATURES-20260131-162601-LPD: Implement user-specific Bible version preferences with default fallback (Option B)
