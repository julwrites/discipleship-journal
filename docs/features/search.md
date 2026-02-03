# Advanced Search

## Overview
The Advanced Search feature allows users to filter and sort their journal notes to quickly find specific content. This goes beyond simple keyword matching by including date range filtering, tag filtering, and flexible sorting options.

## User Stories
- As a user, I want to search my notes by keywords so I can find specific entries.
- As a user, I want to filter notes by a date range (e.g., "Last 7 days", "Last 30 days") to review recent entries.
- As a user, I want to filter notes by a specific tag (e.g., "Sermon", "Prayer") to see all related notes.
- As a user, I want to sort my notes by "Last Updated" or "Created Date" to organize them as I prefer.
- As a user, I want to sort my notes in ascending or descending order.

## API Endpoints

### Notes Search
- `GET /api/notes`: Fetch notes with optional filters.
  - `q`: Search query string (partial match on title or content).
  - `tag`: Filter by tag name (exact match).
  - `startDate`: Filter notes created or updated after this date (RFC3339).
  - `endDate`: Filter notes created or updated before this date (RFC3339).
  - `sortBy`: Field to sort by (`updated_at`, `created_at`, `title`). Default: `updated_at`.
  - `sortOrder`: Direction to sort (`asc`, `desc`). Default: `desc`.

## Data Model

### NoteFilter
(Internal Service Struct)
- `SearchQuery`: String
- `Tag`: String
- `StartDate`: *Time
- `EndDate`: *Time
- `SortBy`: String
- `SortOrder`: String

## Frontend Implementation
- **Dashboard/NotesPage**: The main notes list provides a search bar and filter controls.
- **Filter Dialog**: A UI component that allows users to select:
  - Sort By (Updated, Created, Title)
  - Sort Order (Newest First, Oldest First)
  - Date Range (All Time, Today, Last 7 Days, Last 30 Days, Custom Range)
  - Tag (Select from available tags)
- **Service Layer**: `fetchNotes` in `web/src/services/api.ts` constructs the query string based on the `NoteFilter` object.

## Future Enhancements
- Save search presets.
- Highlight search terms in results.
