---
id: FEATURES-20260131-162514-WAF
status: completed
title: Apply user's default Bible version when cloning memory verse packs
priority: medium
created: 2026-01-31 16:25:14
category: features
dependencies: 
type: task
---

# Apply user's default Bible version when cloning memory verse packs

## Task Information
- **Dependencies**: None (can be implemented independently)
- **Related tasks**: FEATURES-20260131-162336-UYJ (Improve memory verse version selection)

## Task Details
Currently when users clone a memory verse pack (public/system pack), all verses retain their original Bible versions (typically 'ESV'). This task modifies the cloning behavior to apply the user's default Bible version (from settings) to all verses in the cloned pack, providing a more personalized experience.

### Current Behavior:
- **Clone operation**: Copies all verses with their original `version` field values
- **User default version**: Stored in `user.settings.bible_version`
- **Result**: Cloned pack verses keep original versions (e.g., 'ESV'), requiring manual updates if user prefers different version

### Desired Behavior:
When cloning a pack, apply user's default Bible version to all verses in the new pack, unless user explicitly chooses to keep original versions.

### Implementation Options:

#### Option A (Automatic):
- Always apply user's default version during clone
- Simple, consistent behavior
- May surprise users who want to preserve original versions

#### Option B (Prompt):
- Show dialog during clone: "Apply your default version (ESV) to all verses?" with options:
  - "Use my default version (ESV)"
  - "Keep original versions"
- More user-friendly, provides choice
- Requires UI changes

#### Option C (Settings preference):
- Add user preference: "Always apply my default version when cloning packs"
- Respects user choice without prompting each time
- More complex, requires new user setting

### Recommended: Option B (Prompt)
Provides best balance of user control and personalized experience.

### Steps:

#### 1. Backend Changes (`api/services/memory_verse_service.go`):
- Modify `ClonePack()` method (line ~270)
- Add `userDefaultVersion` parameter
- Update verse creation to use `userDefaultVersion` instead of original `version`
- Or implement conditional logic based on user choice

#### 2. Frontend Changes (`web/src/pages/MemoryVersesPage.tsx`):
- Add confirmation dialog when cloning public packs
- Dialog options: "Use my default version" / "Keep original versions"
- Pass user's default version (from settings) to clone API call
- Handle API response appropriately

#### 3. UI/UX Considerations:
- Show user's default version in dialog (e.g., "Use my default version (NIV)")
- Clear labeling of options
- Optional: "Remember my choice" checkbox for future clones
- Ensure dialog works on mobile and desktop

#### 4. Testing:
- Test clone with "Use my default version" option
- Test clone with "Keep original versions" option
- Verify verses have correct versions in database
- Test with different default versions (NIV, KJV, etc.)
- Test error handling and rollback

### Acceptance Criteria
- [ ] Users are prompted for version preference when cloning packs
- [ ] Option "Use my default version" applies user's default Bible version to all cloned verses
- [ ] Option "Keep original versions" preserves original verse versions
- [ ] Cloned packs are created successfully with correct verse versions
- [ ] User's default version is correctly retrieved from settings
- [ ] No regression in existing clone functionality
- [ ] UI is clear and accessible

## Implementation Status
### Completed Work
- ✅ Current cloning behavior analyzed
- ✅ Implementation options evaluated

### Blockers
None yet.

## Notes
- **User Experience**: Option B (prompt) is recommended as it balances automation with user control
- **Performance**: No significant performance impact expected
- **Edge Cases**: Handle null/empty user default version (fallback to 'ESV')
- **Internationalization**: Consider future translation of dialog text
- **Accessibility**: Ensure dialog is keyboard navigable and screen reader friendly
