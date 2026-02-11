---
id: FEATURES-20260131-163357-YNQ
status: completed
title: Add Bible passage retrieval to reading plans: allow users to read passages when clicking references
priority: medium
created: 2026-01-31 16:33:57
category: features
dependencies:
type: task
---

# Add Bible passage retrieval to reading plans: allow users to read passages when clicking references

## Task Information
- **Dependencies**: None (can be implemented independently)
- **Related tasks**: FEATURES-20251220-133329-KLU (reading plans backend)

## Task Details
Currently, Bible reading plans only show passage references (e.g., "Genesis 1-3") as text. Users cannot read the actual Bible text without leaving the app. This task adds the ability to click on references and retrieve the full passage text using the existing Bible AI API integration.

### Current Implementation Analysis:
- **Bible API**: `GET /api/bible/passage?ref={reference}&version={version}` endpoint exists
- **Frontend API**: `getBiblePassage()` function in `services/api.ts`
- **Reference parsing**: `parseBibleReference()` utility normalizes references
- **Reading plan UI**: `ReadingPlanDetail.tsx` shows references as plain text
- **Version handling**: User's default Bible version available in settings

### Steps:

#### 1. Frontend Changes (`web/src/pages/ReadingPlanDetail.tsx`):
- Make passage references clickable (links or buttons)
- Add click handler to fetch passage text using `getBiblePassage()`
- Show loading state while fetching passage
- Display passage in modal, expandable section, or new page
- Handle errors (network issues, invalid references)

#### 2. Passage Display Component:
- Create reusable `BiblePassageModal` or `BiblePassageViewer` component
- Display passage text with formatting
- Include reference and version information
- Add close button and optional actions (share, copy, create note)
- Support user's default Bible version from settings

#### 3. Bible Version Integration:
- Use user's default Bible version from `user.settings.bible_version`
- Fallback to 'ESV' if no default set
- Optional: Add version selector in passage viewer

#### 4. Performance Considerations:
- Cache fetched passages locally (sessionStorage or React state)
- Avoid refetching same passage multiple times
- Implement debouncing for rapid clicks

#### 5. UI/UX Design:
- Visual indication that references are clickable (underline, color, cursor)
- Loading spinner or skeleton while fetching
- Readable passage formatting (line breaks, verse numbers)
- Mobile-responsive design for passage viewer
- Accessibility: keyboard navigation, screen reader support

#### 6. Testing:
- Test with various reference formats (single verse, range, multiple chapters)
- Test with different Bible versions
- Test error handling for invalid references
- Test mobile and desktop layouts
- Verify passage caching works

### Acceptance Criteria
- [ ] Passage references are clickable in reading plan view
- [ ] Clicking a reference fetches and displays the Bible passage
- [ ] Passage viewer shows text with proper formatting
- [ ] Uses user's default Bible version (fallback to ESV)
- [ ] Loading states shown during API fetch
- [ ] Error handling for failed fetches
- [ ] Mobile-responsive passage viewer
- [ ] No regression in existing reading plan functionality
- [ ] Accessible UI (keyboard navigation, screen readers)

## Implementation Status
### Completed Work
- ✅ Codebase exploration completed
- ✅ Integration design planned
- ✅ Frontend component `BiblePassageDialog` implemented
- ✅ Integration with `ReadingPlanDetail` (click to view passage)
- ✅ Bible Passage API is working
- ✅ Added version selector to `BiblePassageDialog`

### Blockers
None.

## Notes
- **Bible API limitations**: Some references may not be supported by the Bible AI API
- **Reference parsing**: Use existing `parseBibleReference()` utility
- **Performance**: Consider client-side caching of fetched passages
- **UI pattern**: Modal vs inline expansion vs new page decision needed
- **Future enhancements**: Could add highlighting, note-taking within passage viewer
- **Integration**: Could later connect with note creation feature
