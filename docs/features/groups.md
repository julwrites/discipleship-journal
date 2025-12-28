# Social Groups & Sharing

## Overview
The Social Groups feature allows users to form discipleship groups, manage membership, and share journal notes within the group for accountability and encouragement.

## User Stories
- As a user, I want to create a group so I can invite others to join me.
- As a user, I want to join a group to participate in community.
- As a user, I want to see which groups I belong to.
- As a group admin, I want to add or remove members.
- As a user, I want to share a specific journal note with my group so they can read it.
- As a group member, I want to see notes shared by others in the group.

## API Endpoints

### Group Management
- `POST /api/groups`: Create a new group. The creator becomes the admin.
- `GET /api/groups`: List groups the current user is a member of.
- `GET /api/groups/search?q={name}`: Search for groups by name.
- `POST /api/groups/{id}/join`: Join a group (currently open join).
- `DELETE /api/groups/{id}/leave`: Leave a group.
- `GET /api/groups/{id}/members`: List members of a group (Members only).
- `POST /api/groups/{id}/members`: Add a member to the group (Admin only).
- `DELETE /api/groups/{id}/members/{userId}`: Remove a member from the group (Admin only).

### Group Sharing
- `POST /api/groups/{id}/shares`: Share a note with the group. Requires `note_id` and optional `comment`.
- `GET /api/groups/{id}/shares`: List shared notes in the group.
- `GET /api/groups/{id}/shares/{shareId}`: Get details of a shared note (including the note content and comment).

## Data Model

### Group
- `ID`: UUID
- `Name`: String
- `Description`: String
- `CreatedBy`: UUID (User ID)

### GroupMember
- `GroupID`: UUID
- `UserID`: UUID
- `Role`: Enum (admin, member)
- `JoinedAt`: Timestamp

### GroupShare
- `ID`: UUID
- `GroupID`: UUID
- `NoteID`: UUID
- `UserID`: UUID (Sharer)
- `Comment`: String
- `SharedAt`: Timestamp

## Frontend Implementation
- **GroupsPage**: Lists user's groups and allows searching/joining/creating groups.
- **GroupDetail**: Shows members and shared notes feed.
- **NoteEditor**: Includes a "Share" action to share the current note to a group.

## Future Enhancements
- Invite-only groups (links or email invites).
- Multiple admins.
- Group chat/discussion board independent of notes.
- Comments on shared notes.
