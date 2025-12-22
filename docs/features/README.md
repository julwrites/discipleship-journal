# Feature Matrix

| Feature | Status | Description |
| :--- | :--- | :--- |
| **Authentication** | Implemented | User sign-up and sign-in using Firebase Auth (Google, Email, Passwordless). |
| **User Profile** | Implemented | Create and update user profile information (`/api/users/me`). |
| **Journaling** | Implemented | Create, read, and update notes (`/api/notes`). |
| **Bible Integration** | Implemented | Fetch Bible passages (`/api/bible/passage`). |
| **AI Chat** | Implemented | Chat with an AI assistant (`/api/chat`). |
| **AI Ask** | Implemented | Ask specific questions to AI (`/api/ai/ask`). |
| **Connections** | Implemented | Connect with other users (`/api/connections`). |
| **Social Groups** | Implemented | Create and join groups for discipleship (`/api/groups`). |
| **Note Sharing** | Implemented | Share notes with groups (`/api/groups/{id}/shares`). |
| **PWA Support** | Implemented | Installable app with offline capabilities. |
| **Delete Notes** | Implemented | Delete notes via API and UI (Hard Delete currently). |
| **Rich Text Editing** | Implemented | Markdown-based rich text editor (`Tiptap`). |
| **Auto-Save** | Implemented | Debounced auto-save for notes. |
| **Push Notifications** | Implemented | FCM notifications for groups and messages. |
| **Bible Reading Plans** | Implemented | Subscribe to reading plans and track progress ([Details](reading_plans.md)). |

## Future Features
*   Search/Filter notes (Debounced search implemented, advanced filtering pending).
