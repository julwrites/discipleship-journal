---
id: FEATURES-20260131-163408-ODS
status: review_requested
title: Implement tagging system for notes: support categories and search by tags
priority: medium
created: 2026-01-31 16:34:08
category: features
dependencies: 
type: task
---

# Implement tagging system for notes: support categories and search by tags

## Task Information
- **Dependencies**: None (foundational feature)
- **Related tasks**: FEATURES-20260131-163414-VGI (reading plan note integration depends on this)

## Task Details
Currently, notes do not support tagging or categorization, limiting organizational capabilities. This task implements a tagging system that allows users to assign tags to notes, filter/search by tags, and manage tag vocabulary.

### Current Implementation Analysis:
- **Notes schema**: `notes` table has `title`, `content` (JSONB), `status` fields
- **No tags**: No existing tag support in database or API
- **Search**: `GET /api/notes` supports search on `title` and `content::text`
- **Frontend**: Note creation/edit has no tag input fields

### Design Options:

#### Option A: Simple tags column (JSONB array)
- Add `tags` JSONB column to `notes` table, default `[]`
- Simple implementation, good for basic tagging
- Limited to simple string matching for search

#### Option B: Separate tags table with many-to-many relationship
- `tags` table: `id`, `name`, `user_id` (for user-specific tags)
- `note_tags` junction table: `note_id`, `tag_id`
- More flexible, supports tag management, counts, etc.
- More complex implementation

#### Option C: Hybrid approach
- `tags` JSONB column for simplicity
- Materialized view or computed column for search optimization

### Recommended: Option B (Separate tables)
Better long-term flexibility for tag management, analytics, and future features.

### Steps:

#### 1. Database Migration:
- Create `tags` table: `id` (UUID), `name` (VARCHAR), `user_id` (UUID, nullable for system tags), `created_at`
- Create `note_tags` junction table: `note_id` (UUID), `tag_id` (UUID), `created_at`
- Add unique constraints: `UNIQUE(user_id, name)` for tags, `UNIQUE(note_id, tag_id)` for note_tags
- Create indexes for performance

#### 2. Backend Models (`api/models/`):
- Add `Tag` and `NoteTag` model structs
- Update `Note` model to include `Tags []Tag` field

#### 3. Service Layer (`api/services/note_service.go`):
- Add tag management methods: `CreateTag`, `GetUserTags`, `DeleteTag`
- Update note CRUD to handle tags: associate/disassociate tags
- Add tag filtering to `GetNotes` method
- Add tag search functionality

#### 4. API Endpoints (`api/handlers/note.go`):
- `GET /api/tags` - List user's tags
- `POST /api/tags` - Create new tag
- `DELETE /api/tags/{id}` - Delete tag
- Extend note endpoints to accept `tags` array in request body
- Add `?tag=name` query parameter to `GET /api/notes`

#### 5. Frontend Changes:
- **Tag input component**: Create `TagInput` component with auto-complete
- **Note form**: Add tag field to note creation/edit
- **Tag management page**: List, create, delete tags
- **Note list filtering**: Add tag filter to dashboard/search
- **Tag display**: Show tags on note cards and detail view

#### 6. Search Integration:
- Update note search to filter by tags
- Add tag-based filtering in dashboard
- Consider tag cloud or popular tags display

#### 7. Testing:
- Test tag creation, assignment, removal
- Test note filtering by tags
- Test tag management (create, list, delete)
- Verify user isolation (users only see their own tags)
- Performance testing with many tags/notes

### Acceptance Criteria
- [ ] Database schema with tags and note_tags tables
- [ ] API endpoints for tag management
- [ ] Note CRUD operations support tags
- [ ] Frontend tag input component with auto-complete
- [ ] Note creation/edit includes tag selection
- [ ] Note filtering/search by tags
- [ ] Tags displayed on note cards and detail view
- [ ] User isolation: users only manage/see their own tags
- [ ] No regression in existing note functionality
- [ ] Performance: tag operations don't significantly impact note loading

## Implementation Status
### Completed Work
- ✅ Codebase exploration completed
- ✅ Design options evaluated

### Blockers
None yet.

## Notes
- **System tags**: Consider reserved system tags like 'Bible Reading' (could be user-creatable)
- **Tag normalization**: Case-insensitive, trim whitespace, prevent duplicates
- **Performance**: Indexes on `note_tags(note_id, tag_id)` and `tags(user_id, name)`
- **Migration**: Existing notes get empty tags array
- **UI/UX**: Tag input should support typing new tags and selecting existing
- **Future extensions**: Tag colors, tag hierarchies, tag analytics
