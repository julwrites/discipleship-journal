# Journaling (Notes)

## Overview
The Journaling feature allows users to create, edit, and manage rich text notes. Users can associate these notes with Bible passages, organize them using tags, and revisit them later. This feature is core to the "Discipleship Journal" concept, encouraging reflection and study.

## User Stories
- As a user, I want to create a note with a title and rich content (bold, italics, lists, images) to record my thoughts.
- As a user, I want to tag my notes (e.g., "Sermon", "Quiet Time") so I can organize and find them easily.
- As a user, I want to associate a note with a Bible reference so I can recall the context of my reflection.
- As a user, I want to edit existing notes to update or correct information.
- As a user, I want to delete notes I no longer need (soft delete) so they don't clutter my view.
- As a user, I want to view a list of all my notes, sorted by date or title.
- As a user, I want to search my notes by keyword or tag.

## API Endpoints

### Notes
- `GET /api/notes`: Fetch notes. Supports pagination, search (`q`), filtering (`tag`, `startDate`, `endDate`), and sorting (`sortBy`, `sortOrder`).
  - Search (`q`) performs a partial match on title and content. Title searches are optimized using a trigram index (`pg_trgm`).
- `POST /api/notes`: Create a new note. Payload: `{ title, content (JSON), tags: [] }`.
- `GET /api/notes/{id}`: Fetch a specific note by ID.
- `PUT /api/notes/{id}`: Update a note. Payload: `{ title, content (JSON), tags: [] }`.
- `DELETE /api/notes/{id}`: Soft delete a note.

### Tags
- `GET /api/tags`: Fetch all tags created by the user.
- `POST /api/tags`: Create a new tag. Payload: `{ name }`.
- `DELETE /api/tags/{id}`: Delete a tag.

## Data Model

### Note
- `ID`: UUID (Primary Key)
- `UserID`: UUID (Foreign Key to Users)
- `Title`: String
- `Content`: JSON (TipTap/ProseMirror JSON format)
- `CreatedAt`: Timestamp
- `UpdatedAt`: Timestamp
- `DeletedAt`: Timestamp (Nullable, for soft delete)
- `Tags`: Array of Tag objects (derived from many-to-many relationship)

### Tag
- `ID`: UUID (Primary Key)
- `UserID`: UUID (Foreign Key to Users)
- `Name`: String
- `CreatedAt`: Timestamp

### NoteTag
- `NoteID`: UUID
- `TagID`: UUID

## Frontend Implementation
- **NotesPage**: Displays a list of notes with a search bar and filter controls.
- **NoteEditor**: The core component for creating and editing notes.
  - **Rich Text**: Uses [TipTap](https://tiptap.dev/) for rich text editing (Markdown support, bold, italic, lists).
  - **Images**: Supports adding images via URL, including support for Alt Text to improve accessibility.
  - **Tags**: Allows selecting existing tags or creating new ones on the fly.
  - **Bible Context**: Can be pre-populated with a Bible reference (e.g., when navigating from a Reading Plan).
- **TagsPage**: A dedicated page (`/tags`) for managing tag definitions (create, delete).

## Soft Delete
Notes are "soft deleted" by setting the `deleted_at` timestamp.
- **Behavior**: Deleted notes are excluded from default API queries (`GET /api/notes`).
- **Restoration**: Currently, there is no UI to restore deleted notes, but the data is preserved in the database.

## Future Enhancements
- Note templates.
- "Trash" view to restore soft-deleted notes.
- Public sharing of notes.
- Collaborative editing.
