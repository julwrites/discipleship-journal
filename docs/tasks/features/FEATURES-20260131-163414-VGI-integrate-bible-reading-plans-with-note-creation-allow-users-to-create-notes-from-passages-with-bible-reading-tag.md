---
id: FEATURES-20260131-163414-VGI
status: pending
title: Integrate Bible reading plans with note creation: allow users to create notes from passages with 'Bible Reading' tag
priority: medium
created: 2026-01-31 16:34:14
category: features
dependencies: FEATURES-20260131-163408-ODS 
type: task
---

# Integrate Bible reading plans with note creation: allow users to create notes from passages with 'Bible Reading' tag

## Task Information
- **Dependencies**: FEATURES-20260131-163408-ODS (tagging system must be implemented first)
- **Related tasks**:
  - FEATURES-20260131-163357-YNQ (passage retrieval)
  - FEATURES-20251220-133329-KLU (reading plans backend)

## Task Details
Allow users to create notes directly from Bible reading plan passages, automatically tagged with 'Bible Reading' for easy organization and search. This integrates reading plans with the note-taking system, enabling users to journal about their daily readings.

### Current Implementation Analysis:
- **Reading plans**: Show passage references, no connection to notes
- **Notes system**: Separate from reading plans
- **Tagging**: Will be implemented in FEATURES-20260131-163408-ODS
- **Bible API**: Passage retrieval available via `getBiblePassage()`

### User Workflow:
1. User views reading plan day with passage reference
2. User can click "Create Note" button
3. System creates note pre-populated with:
   - Title: "Bible Reading: [Passage Reference]"
   - Content: Passage text (if available) or reference
   - Tags: "Bible Reading" (automatically added)
   - Optional: Date, reading plan info
4. User can add personal notes/comments to the note
5. Note appears in user's notes, searchable by "Bible Reading" tag

### Steps:

#### 1. Frontend Integration (`web/src/pages/ReadingPlanDetail.tsx`):
- Add "Create Note" button next to each passage reference
- Button opens note creation dialog/modal
- Pre-populate note form with passage data:
  - Title: `Bible Reading: ${passageReference}`
  - Content: Passage text (fetch via `getBiblePassage()` if not already loaded)
  - Tags: "Bible Reading" (pre-selected, user can add more)
- Include reading plan context in note (optional field)

#### 2. Note Creation API Integration:
- Use existing `createNote()` API function
- Ensure tag support is implemented (dependency)
- Pass `tags: ["Bible Reading"]` in note creation request

#### 3. Context Preservation:
- Store reading plan metadata in note content or custom fields:
  - Reading plan title
  - Day number
  - Date read
- Consider adding `context` JSON field to notes for extended metadata

#### 4. UI/UX Considerations:
- "Create Note" button visibility: Always visible or only after passage is viewed?
- Progress integration: Optionally mark day as complete when note is created
- Note editor: Use rich text editor for user's comments
- Success feedback: Show confirmation, optionally navigate to new note

#### 5. Tag Management:
- Ensure "Bible Reading" tag exists (create on first use if needed)
- Consider making it a system tag (not user-deletable) or user-managed
- Tag search: Users can filter notes by "Bible Reading" tag

#### 6. Testing:
- Test note creation from reading plan passages
- Verify "Bible Reading" tag is automatically added
- Test with passage text fetching (success and error cases)
- Verify note search by "Bible Reading" tag works
- Test mobile and desktop interactions

### Acceptance Criteria
- [ ] "Create Note" button available in reading plan view
- [ ] Note creation form pre-populated with passage reference
- [ ] "Bible Reading" tag automatically added to note
- [ ] Passage text can be included in note content (if available)
- [ ] Created notes appear in user's notes list
- [ ] Notes searchable/filterable by "Bible Reading" tag
- [ ] User can edit note after creation (add personal comments)
- [ ] No regression in existing reading plan or note functionality
- [ ] Mobile-responsive design

## Implementation Status
### Completed Work
- ✅ Codebase exploration completed
- ✅ Integration design planned

### Blockers
- Depends on tagging system implementation (FEATURES-20260131-163408-ODS)

## Notes
- **Tag creation**: "Bible Reading" tag should be created automatically if not exists
- **Content format**: Consider storing passage text as markdown or structured JSON
- **Multiple passages**: Some reading plan days have multiple passages (e.g., "Genesis 1-3, Psalm 1")
- **Future enhancements**:
  - Template notes with reflection prompts
  - Daily reading journal view
  - Reading streak tracking with notes
  - Share reading notes with groups
