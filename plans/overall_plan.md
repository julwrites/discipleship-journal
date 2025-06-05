# PWA Web App Development Plan: Discipleship Journal

## Objective
To develop a Progressive Web Application (PWA) that enables users to log in, connect with other users, form groups, and share personal or group notes, with robust integration for embedding, highlighting, commenting on, and exploring translations of Bible passages.

## Core Features Summary
*   User Login & Authentication
*   User Connections (friending system)
*   User Groups with Membership Management
*   Note Sharing:
    *   1-to-1 Private Notes
    *   Group Notes
*   Bible Passage Integration:
    *   Embedding passages in notes
    *   Highlighting text within passages
    *   Commenting on verses/highlights
    *   Viewing Hebrew/Greek translations

## Target Audience (Initial)
*   Individuals engaged in personal Bible study.
*   Members of small Bible study groups.
*   Students of theology or ministry.
*   Anyone looking for a collaborative note-taking tool focused on biblical texts.

## Technology Stack (Proposed - subject to final decision)
*   **Frontend (PWA):** React or Vue.js, State Management (Redux/Zustand or Vuex/Pinia), PWA libraries.
*   **Backend:** Node.js (Express/NestJS) or Firebase/Supabase (BaaS).
*   **Database:** PostgreSQL or Firestore/MongoDB (depending on backend choice).
*   **Bible API:** To be selected (e.g., ESV API, NET Bible API, Digital Bible Platform).
*   **Authentication:** Firebase Auth, Auth0, or custom JWT implementation.

## Development Plan

### Module 1: User Authentication & Profile

**US001: New User Signup**
*   **As a** new user,
*   **I want to** sign up with my email address and a secure password,
*   **So that** I can create an account and access the application's features.
*   **Acceptance Criteria:**
    1.  User can enter an email and password (with confirmation).
    2.  Password meets minimum strength requirements (e.g., 8 chars, mix of types).
    3.  Email uniqueness is validated.
    4.  Upon successful signup, the user is authenticated and redirected to the dashboard/home page.
    5.  Appropriate error messages are shown for invalid input or existing email.
*   **Test Scenarios:**
    1.  **TS001.1 (Happy Path):** Valid email, strong matching passwords -> Signup successful, user logged in, redirected.
    2.  **TS001.2 (Invalid Email):** Invalid email format -> Error message shown, no account created.
    3.  **TS001.3 (Weak Password):** Password too short or simple -> Error message shown, no account created.
    4.  **TS001.4 (Password Mismatch):** Password and confirmation don't match -> Error message shown, no account created.
    5.  **TS001.5 (Existing Email):** Email already registered -> Error message shown, no account created.

**Task Breakdown for US001:**
*   **Task-001.1: Design User Schema/Model**
    *   **Dependencies:** None
    *   **Implementation:** Define database schema for users (ID, email, password_hash, name (optional), created_at, updated_at).
    *   **Verification:** Schema accommodates all required user information.
    *   **Test Details:** Schema review, tested implicitly via API tests.
    *   **Status:** Not Started
*   **Task-001.2: Develop Signup UI Form**
    *   **Dependencies:** None
    *   **Implementation:** Create frontend form with fields for email, password, confirm password. Basic client-side validation for format and non-empty.
    *   **Verification:** Form renders correctly, inputs are captured.
    *   **Test Details:** Component rendering tests, visual inspection.
    *   **Status:** Not Started
*   **Task-001.3: Implement Backend Signup API Endpoint**
    *   **Dependencies:** Task-001.1
    *   **Implementation:** Create `POST /api/auth/signup` endpoint. Validate input, check email uniqueness, hash password, store user in DB. Return session token/cookie.
    *   **Verification:** Endpoint correctly processes valid/invalid requests, creates user, returns appropriate responses.
    *   **Test Details:** API unit/integration tests for success and error cases (Postman/automated).
    *   **Status:** Not Started
*   **Task-001.4: Integrate Signup UI with Backend API**
    *   **Dependencies:** Task-001.2, Task-001.3
    *   **Implementation:** Frontend form submission calls the signup API. Handle API responses, display errors, or redirect on success.
    *   **Verification:** End-to-end signup flow works as per ACs.
    *   **Test Details:** E2E test for the entire signup process.
    *   **Status:** Not Started

**US002: User Login**
*   **As a** returning user,
*   **I want to** log in with my registered email and password,
*   **So that** I can access my account and continue using the application.
*   **Acceptance Criteria:**
    1.  User can enter their email and password.
    2.  Upon successful authentication, the user is redirected to the dashboard/home page.
    3.  Appropriate error messages are shown for incorrect credentials or non-existent accounts.
    4.  Session/token is established.
*   **Test Scenarios:**
    1.  **TS002.1 (Happy Path):** Valid credentials -> Login successful, user logged in, redirected.
    2.  **TS002.2 (Invalid Password):** Correct email, incorrect password -> Error message shown.
    3.  **TS002.3 (Non-existent Email):** Email not registered -> Error message shown.
    4.  **TS002.4 (Empty Fields):** Submitting with empty fields -> Error message shown.

**Task Breakdown for US002:**
*   **Task-002.1: Develop Login UI Form**
    *   **Dependencies:** None
    *   **Implementation:** Create frontend form with fields for email and password.
    *   **Verification:** Form renders correctly, inputs are captured.
    *   **Test Details:** Component rendering tests, visual inspection.
    *   **Status:** Not Started
*   **Task-002.2: Implement Backend Login API Endpoint**
    *   **Dependencies:** Task-001.1 (User Schema)
    *   **Implementation:** Create `POST /api/auth/login` endpoint. Validate input, find user by email, compare hashed password. Return session token/cookie.
    *   **Verification:** Endpoint correctly processes valid/invalid credentials, returns appropriate responses.
    *   **Test Details:** API unit/integration tests for success and error cases.
    *   **Status:** Not Started
*   **Task-002.3: Integrate Login UI with Backend API**
    *   **Dependencies:** Task-002.1, Task-002.2
    *   **Implementation:** Frontend form submission calls the login API. Handle API responses, display errors, or redirect on success. Manage session/token on client.
    *   **Verification:** End-to-end login flow works as per ACs.
    *   **Test Details:** E2E test for the entire login process.
    *   **Status:** Not Started

**US003: User Logout**
*   **As an** authenticated user,
*   **I want to** log out of the application,
*   **So that** my session is terminated and my account is secured.
*   **Acceptance Criteria:**
    1.  A logout option is available.
    2.  Clicking logout invalidates the current session/token.
    3.  User is redirected to the login page or public home page.
    4.  Previously protected routes are no longer accessible.
*   **Test Scenarios:**
    1.  **TS003.1 (Happy Path):** User clicks logout -> Session invalidated, redirected to login.
    2.  **TS003.2 (Access Protected Route After Logout):** Attempt to access dashboard -> Redirected to login.

**Task Breakdown for US003:**
*   **Task-003.1: Develop Logout UI Element**
    *   **Dependencies:** None
    *   **Implementation:** Add a logout button/link in the application's navigation/profile menu.
    *   **Verification:** Logout element is visible and clickable.
    *   **Test Details:** Visual inspection.
    *   **Status:** Not Started
*   **Task-003.2: Implement Backend Logout API Endpoint (Optional, if server-side session invalidation needed)**
    *   **Dependencies:** None
    *   **Implementation:** Create `POST /api/auth/logout` endpoint to invalidate server-side session or token blacklist.
    *   **Verification:** Session is correctly invalidated on the server.
    *   **Test Details:** API unit test.
    *   **Status:** Not Started
*   **Task-003.3: Implement Client-Side Logout Logic**
    *   **Dependencies:** Task-003.1
    *   **Implementation:** On logout click, clear client-side session/token, call logout API (if any), redirect user.
    *   **Verification:** Client-side token removed, user redirected.
    *   **Test Details:** E2E test for logout and subsequent access attempt.
    *   **Status:** Not Started

---

### Module 2: User Connections

**US004: Search for Other Users**
*   **As a** user,
*   **I want to** search for other registered users by their username or email,
*   **So that** I can find people I know and initiate a connection.
*   **Acceptance Criteria:**
    1.  A search interface is available.
    2.  Users can input a search term (username/email).
    3.  Search results display matching users (e.g., username, profile picture if available).
    4.  Search results do not display the searching user.
    5.  An option to send a connection request is available next to each search result (if not already connected or request pending).
*   **Test Scenarios:**
    1.  **TS004.1 (Happy Path - Username):** Search existing username -> User found.
    2.  **TS004.2 (Happy Path - Email):** Search existing email -> User found.
    3.  **TS004.3 (Partial Match):** Search partial username -> Matching users found.
    4.  **TS004.4 (No Match):** Search non-existent username -> "No users found" message.
    5.  **TS004.5 (Search Self):** Search for own username -> Self not listed.

**Task Breakdown for US004:**
*   **Task-004.1: Design Connection/Friendship Schema/Model**
    *   **Dependencies:** Task-001.1
    *   **Implementation:** Define DB schema for connections (e.g., `user_connections` table: `id`, `requester_id`, `addressee_id`, `status` (pending, accepted, rejected, blocked), `created_at`, `updated_at`).
    *   **Verification:** Schema supports different states of a connection.
    *   **Test Details:** Schema review.
    *   **Status:** Not Started
*   **Task-004.2: Develop User Search UI**
    *   **Dependencies:** None
    *   **Implementation:** Create frontend search bar and results display area.
    *   **Verification:** UI elements render correctly.
    *   **Test Details:** Component rendering tests.
    *   **Status:** Not Started
*   **Task-004.3: Implement Backend User Search API Endpoint**
    *   **Dependencies:** Task-001.1
    *   **Implementation:** Create `GET /api/users/search?q=<term>` endpoint. Query users table by username/email, excluding the requester. Paginate results.
    *   **Verification:** API returns correct users based on search term, handles no results.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-004.4: Integrate Search UI with Backend API**
    *   **Dependencies:** Task-004.2, Task-004.3
    *   **Implementation:** Frontend calls search API on input, displays results. Includes "Send Request" button logic.
    *   **Verification:** Search results are displayed dynamically and correctly.
    *   **Test Details:** E2E test for user search.
    *   **Status:** Not Started

**US005: Send Connection Request**
*   **As a** user,
*   **I want to** send a connection request to another user,
*   **So that** we can become connected if they accept.
*   **Acceptance Criteria:**
    1.  A "Send Request" option is available for users not yet connected and with no pending request.
    2.  Clicking "Send Request" creates a pending connection record.
    3.  The sender sees an indication that the request is pending (e.g., button changes to "Request Sent").
    4.  The recipient receives a notification (in-app initially).
*   **Test Scenarios:**
    1.  **TS005.1 (Happy Path):** User A sends request to User B -> Request status pending, User A sees "Request Sent", User B gets notification.
    2.  **TS005.2 (Already Connected):** Attempt to send request to already connected user -> Option disabled or error.
    3.  **TS005.3 (Request Already Sent):** Attempt to send request when one is pending -> Option disabled or error.

**Task Breakdown for US005:**
*   **Task-005.1: Implement Backend Send Connection Request API Endpoint**
    *   **Dependencies:** Task-004.1
    *   **Implementation:** Create `POST /api/connections/request` endpoint. Takes target user ID. Creates a connection record with 'pending' status. Prevents duplicate requests or requests to self.
    *   **Verification:** API creates pending connection, handles duplicates/self-requests.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-005.2: Integrate "Send Request" button with API**
    *   **Dependencies:** Task-004.4 (UI for button), Task-005.1
    *   **Implementation:** Frontend "Send Request" button calls the API. Updates UI to show "Request Sent" or similar.
    *   **Verification:** UI updates correctly after sending a request.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started
*   **Task-005.3: Implement Basic In-App Notification for Recipient**
    *   **Dependencies:** Task-005.1
    *   **Implementation:** When a connection request is created, the recipient user gets an in-app indicator (e.g., a badge on a "Requests" icon).
    *   **Verification:** Recipient sees a notification for a new request.
    *   **Test Details:** Manual test or E2E if notification system is testable.
    *   **Status:** Not Started

**US006: Manage Connection Requests**
*   **As a** user,
*   **I want to** view incoming connection requests and choose to accept or reject them,
*   **So that** I can control who I connect with.
*   **Acceptance Criteria:**
    1.  A dedicated section/page lists incoming connection requests.
    2.  For each request, the sender's username is displayed.
    3.  Options to "Accept" or "Reject" are available.
    4.  Accepting updates the connection status to "accepted"; sender and recipient are now connected.
    5.  Rejecting updates the connection status to "rejected" (or deletes the request).
    6.  The request is removed from the pending list after action.
*   **Test Scenarios:**
    1.  **TS006.1 (Accept Request):** User B accepts User A's request -> Connection status "accepted", request removed from list, both users see each other as connected.
    2.  **TS006.2 (Reject Request):** User B rejects User A's request -> Connection status "rejected"/deleted, request removed from list.
    3.  **TS006.3 (View List):** User with pending requests sees them listed.

**Task Breakdown for US006:**
*   **Task-006.1: Develop UI for Managing Connection Requests**
    *   **Dependencies:** None
    *   **Implementation:** Create frontend page/section to list pending requests with Accept/Reject buttons.
    *   **Verification:** UI displays requests correctly.
    *   **Test Details:** Component rendering tests.
    *   **Status:** Not Started
*   **Task-006.2: Implement Backend API Endpoints for Listing and Responding to Requests**
    *   **Dependencies:** Task-004.1
    *   **Implementation:**
        *   `GET /api/connections/requests/pending` (lists incoming requests).
        *   `POST /api/connections/requests/{request_id}/accept`.
        *   `POST /api/connections/requests/{request_id}/reject`.
        Update connection status in DB.
    *   **Verification:** APIs correctly list requests and update status.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-006.3: Integrate Connection Request Management UI with Backend APIs**
    *   **Dependencies:** Task-006.1, Task-006.2
    *   **Implementation:** Frontend lists requests, calls appropriate API on accept/reject, updates UI.
    *   **Verification:** Managing requests works end-to-end.
    *   **Test Details:** E2E tests.
    *   **Status:** Not Started

**US007: View Connected Users**
*   **As a** user,
*   **I want to** see a list of my established connections,
*   **So that** I can easily identify and interact with them (e.g., for sharing notes).
*   **Acceptance Criteria:**
    1.  A dedicated section/page lists all users with whom an "accepted" connection exists.
    2.  The list displays usernames.
    3.  (Optional) An option to remove a connection.
*   **Test Scenarios:**
    1.  **TS007.1 (View List):** User with connections sees them listed.
    2.  **TS007.2 (No Connections):** User with no connections sees an appropriate message.

**Task Breakdown for US007:**
*   **Task-007.1: Develop UI for Displaying Connected Users**
    *   **Dependencies:** None
    *   **Implementation:** Create frontend page/section to list connected users.
    *   **Verification:** UI renders correctly.
    *   **Test Details:** Component rendering tests.
    *   **Status:** Not Started
*   **Task-007.2: Implement Backend API Endpoint for Listing Connected Users**
    *   **Dependencies:** Task-004.1
    *   **Implementation:** Create `GET /api/connections` endpoint. Returns users where connection status is 'accepted'.
    *   **Verification:** API returns correct list of connected users.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-007.3: Integrate Connected Users UI with Backend API**
    *   **Dependencies:** Task-007.1, Task-007.2
    *   **Implementation:** Frontend calls API and displays the list.
    *   **Verification:** Connected users list is displayed correctly.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started

---

### Module 3: Groups

**US008: Create a Group**
*   **As a** user,
*   **I want to** create a new group with a name and an optional description,
*   **So that** I can establish a shared space for discussions and notes with specific people.
*   **Acceptance Criteria:**
    1.  User can access a "Create Group" form.
    2.  User must provide a group name. Description is optional.
    3.  Upon successful creation, the group is listed in the user's groups, and the creator becomes an admin.
    4.  Appropriate error messages for invalid input (e.g., empty name).
*   **Test Scenarios:**
    1.  **TS008.1 (Happy Path):** Valid name and description -> Group created, user is admin.
    2.  **TS008.2 (Name Only):** Valid name, no description -> Group created.
    3.  **TS008.3 (Empty Name):** No name provided -> Error message, group not created.

**Task Breakdown for US008:**
*   **Task-008.1: Design Group and Group Membership Schemas/Models**
    *   **Dependencies:** Task-001.1
    *   **Implementation:**
        *   `groups` table: `id`, `name`, `description`, `creator_id`, `created_at`, `updated_at`.
        *   `group_members` table: `id`, `group_id`, `user_id`, `role` (admin, member), `joined_at`.
    *   **Verification:** Schemas support group creation and membership with roles.
    *   **Test Details:** Schema review.
    *   **Status:** Not Started
*   **Task-008.2: Develop Create Group UI Form**
    *   **Dependencies:** None
    *   **Implementation:** Frontend form for group name and description.
    *   **Verification:** Form renders correctly.
    *   **Test Details:** Component rendering tests.
    *   **Status:** Not Started
*   **Task-008.3: Implement Backend Create Group API Endpoint**
    *   **Dependencies:** Task-008.1
    *   **Implementation:** `POST /api/groups` endpoint. Validates input. Creates group record. Adds creator to `group_members` as admin.
    *   **Verification:** API creates group and assigns admin role.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-008.4: Integrate Create Group UI with Backend API**
    *   **Dependencies:** Task-008.2, Task-008.3
    *   **Implementation:** Frontend form calls API, handles response.
    *   **Verification:** Group creation works end-to-end.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started

**US009: Invite Users to a Group**
*   **As a** group admin,
*   **I want to** invite my connected users to join a group I manage,
*   **So that** they can become members and participate.
*   **Acceptance Criteria:**
    1.  Group admin can access an "Invite Users" feature within a group they manage.
    2.  Admin can see a list of their connections who are not already members or invited.
    3.  Admin can select users and send invitations.
    4.  Invited users receive a notification (in-app).
    5.  A record of the invitation (pending status) is created.
*   **Test Scenarios:**
    1.  **TS009.1 (Happy Path):** Admin invites connected user -> Invitation sent, user gets notification.
    2.  **TS009.2 (Invite Non-Connected User):** Logic prevents inviting non-connected users (or feature not shown).
    3.  **TS009.3 (Invite Existing Member):** Logic prevents inviting existing members.

**Task Breakdown for US009:**
*   **Task-009.1: Design Group Invitation Schema/Model (if different from membership or if detailed tracking needed)**
    *   **Dependencies:** Task-008.1
    *   **Implementation:** Extend `group_members` with a `status` like 'pending_invitation', 'joined', or have a separate `group_invitations` table.
    *   **Verification:** Schema supports tracking invitations.
    *   **Test Details:** Schema review.
    *   **Status:** Not Started
*   **Task-009.2: Develop Invite Users UI**
    *   **Dependencies:** US007 (View Connected Users - for list)
    *   **Implementation:** Frontend interface within group settings to list connections and select for invitation.
    *   **Verification:** UI displays eligible users and allows selection.
    *   **Test Details:** Component tests.
    *   **Status:** Not Started
*   **Task-009.3: Implement Backend Invite to Group API Endpoint**
    *   **Dependencies:** Task-008.1, Task-009.1
    *   **Implementation:** `POST /api/groups/{group_id}/invitations` endpoint. Takes user IDs. Creates `group_members` records with 'pending_invitation' status (or records in `group_invitations`). Validates admin rights and user eligibility.
    *   **Verification:** API creates pending invitations.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-009.4: Integrate Invite UI with Backend API & Notifications**
    *   **Dependencies:** Task-009.2, Task-009.3
    *   **Implementation:** Frontend calls API, handles response. Trigger in-app notification for invited users.
    *   **Verification:** Invitations sent and notifications received.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started

**US010: Manage Group Invitations (Accept/Reject)**
*   **As a** user who has been invited to a group,
*   **I want to** view my group invitations and choose to accept or reject them,
*   **So that** I can join groups I'm interested in.
*   **Acceptance Criteria:**
    1.  Users can see a list of their pending group invitations.
    2.  For each invitation, group name and inviter are shown.
    3.  Options to "Accept" or "Reject" are available.
    4.  Accepting adds the user to the group as a "member" and removes the pending invitation.
    5.  Rejecting removes the pending invitation.
*   **Test Scenarios:**
    1.  **TS010.1 (Accept Invite):** User accepts -> Joins group as member, invitation removed.
    2.  **TS010.2 (Reject Invite):** User rejects -> Invitation removed, not a member.

**Task Breakdown for US010:**
*   **Task-010.1: Develop UI for Managing Group Invitations**
    *   **Dependencies:** None
    *   **Implementation:** Frontend page/section to list pending group invitations with Accept/Reject options.
    *   **Verification:** UI displays invitations correctly.
    *   **Test Details:** Component tests.
    *   **Status:** Not Started
*   **Task-010.2: Implement Backend API Endpoints for Listing and Responding to Group Invitations**
    *   **Dependencies:** Task-008.1, Task-009.1
    *   **Implementation:**
        *   `GET /api/me/group-invitations/pending`.
        *   `POST /api/me/group-invitations/{invitation_id}/accept` (updates `group_members` status or creates member record).
        *   `POST /api/me/group-invitations/{invitation_id}/reject` (deletes invitation or updates status).
    *   **Verification:** APIs list and process invitations correctly.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-010.3: Integrate Group Invitation Management UI with Backend APIs**
    *   **Dependencies:** Task-010.1, Task-010.2
    *   **Implementation:** Frontend lists invitations, calls appropriate API, updates UI.
    *   **Verification:** Managing group invitations works end-to-end.
    *   **Test Details:** E2E tests.
    *   **Status:** Not Started

**US011: View My Groups**
*   **As a** user,
*   **I want to** see a list of all groups I am a member of,
*   **So that** I can easily navigate to them.
*   **Acceptance Criteria:**
    1.  A dedicated section/page lists all groups the user is a member of.
    2.  The list displays group names.
    3.  Clicking a group name navigates to the group's detail page.
*   **Test Scenarios:**
    1.  **TS011.1 (View List):** User who is a member of groups sees them listed.
    2.  **TS011.2 (No Groups):** User not in any groups sees an appropriate message.

**Task Breakdown for US011:**
*   **Task-011.1: Develop UI for Displaying User's Groups**
    *   **Dependencies:** None
    *   **Implementation:** Frontend page/section to list user's groups.
    *   **Verification:** UI renders correctly.
    *   **Test Details:** Component tests.
    *   **Status:** Not Started
*   **Task-011.2: Implement Backend API Endpoint for Listing User's Groups**
    *   **Dependencies:** Task-008.1
    *   **Implementation:** `GET /api/me/groups` endpoint. Returns groups where user is a member.
    *   **Verification:** API returns correct list of groups.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-011.3: Integrate User's Groups UI with Backend API**
    *   **Dependencies:** Task-011.1, Task-011.2
    *   **Implementation:** Frontend calls API and displays the list of groups.
    *   **Verification:** User's groups list is displayed correctly.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started

---

### Module 4: Notes (Private & Group)

**US012: Create a Private Note**
*   **As a** user,
*   **I want to** create a new private note with a title and content,
*   **So that** I can capture my personal thoughts, reflections, and study notes.
*   **Acceptance Criteria:**
    1.  User can access a "Create Note" interface.
    2.  User can input a title and content for the note.
    3.  The note is saved as private to the user by default.
    4.  The created note appears in the user's list of private notes.
*   **Test Scenarios:**
    1.  **TS012.1 (Happy Path):** User enters title and content, saves -> Note created and listed.
    2.  **TS012.2 (Content Only):** User enters content, no title (system assigns default or prompts) -> Note created.
    3.  **TS012.3 (Empty Note):** Attempt to save empty note -> Error or disabled save.

**Task Breakdown for US012:**
*   **Task-012.1: Design Note Schema/Model**
    *   **Dependencies:** Task-001.1, Task-008.1 (for group association later)
    *   **Implementation:** `notes` table: `id`, `user_id` (creator), `group_id` (nullable, for group notes), `title`, `content` (rich text/markdown), `is_private` (boolean, or inferred if `group_id` is null), `created_at`, `updated_at`.
        *   Also `note_shares` table: `id`, `note_id`, `shared_with_user_id`, `permissions` (view, edit).
    *   **Verification:** Schema supports private notes, group notes, sharing, and content.
    *   **Test Details:** Schema review.
    *   **Status:** Not Started
*   **Task-012.2: Develop Note Editor UI (Basic)**
    *   **Dependencies:** None
    *   **Implementation:** Frontend interface with fields for title and a text area for content (basic initially, can be enhanced to rich text later). Save button.
    *   **Verification:** Editor UI renders, inputs captured.
    *   **Test Details:** Component tests.
    *   **Status:** Not Started
*   **Task-012.3: Implement Backend Create Private Note API Endpoint**
    *   **Dependencies:** Task-012.1
    *   **Implementation:** `POST /api/notes` endpoint. Takes title, content. Sets `user_id` to current user, `is_private` to true (or `group_id` to null).
    *   **Verification:** API creates private note for the user.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-012.4: Integrate Note Editor UI with Backend API for Private Notes**
    *   **Dependencies:** Task-012.2, Task-012.3
    *   **Implementation:** Frontend editor calls API to save note. Handles response.
    *   **Verification:** Private note creation works end-to-end.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started

**US013: View My Private Notes**
*   **As a** user,
*   **I want to** see a list of all my private notes,
*   **So that** I can easily access and manage them.
*   **Acceptance Criteria:**
    1.  A dedicated section lists the user's private notes (and notes shared with them 1-1).
    2.  List displays note titles and a snippet/date.
    3.  Clicking a note opens it for viewing/editing.
*   **Test Scenarios:**
    1.  **TS013.1 (View List):** User with private notes sees them listed.
    2.  **TS013.2 (No Notes):** User with no private notes sees an appropriate message.

**Task Breakdown for US013:** (Tasks for viewing notes shared 1-1 will be covered under US014)
*   **Task-013.1: Develop UI for Listing Private Notes**
    *   **Dependencies:** None
    *   **Implementation:** Frontend page/section to list user's private notes.
    *   **Verification:** UI renders correctly.
    *   **Test Details:** Component tests.
    *   **Status:** Not Started
*   **Task-013.2: Implement Backend API Endpoint for Listing User's Private Notes**
    *   **Dependencies:** Task-012.1
    *   **Implementation:** `GET /api/notes/private` endpoint. Returns notes where `user_id` is current user and `group_id` is null (or `is_private` is true).
    *   **Verification:** API returns correct list of private notes.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-013.3: Integrate Private Notes List UI with Backend API**
    *   **Dependencies:** Task-013.1, Task-013.2
    *   **Implementation:** Frontend calls API and displays the list.
    *   **Verification:** Private notes list is displayed correctly.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started

**US014: Share a Private Note (1-1)**
*   **As a** user,
*   **I want to** share one of my private notes with a specific connected user,
*   **So that** they can view (and optionally edit, TBD) the note.
*   **Acceptance Criteria:**
    1.  From a private note view/edit screen, an option to "Share" is available.
    2.  User can select one of their connections to share with.
    3.  (MVP) Shared note is view-only for the recipient initially.
    4.  The recipient can see the shared note in their "Shared with Me" list or integrated notes list.
*   **Test Scenarios:**
    1.  **TS014.1 (Happy Path):** User A shares note with User B -> User B can view the note.
    2.  **TS014.2 (Share with Non-Connected):** UI prevents sharing with non-connected users.

**Task Breakdown for US014:**
*   **Task-014.1: Develop Share Note UI**
    *   **Dependencies:** US007 (View Connected Users - for list), Task-012.2 (Note Editor UI where share button lives)
    *   **Implementation:** Frontend interface within note view to select a connection and initiate sharing.
    *   **Verification:** UI allows selecting a connection for sharing.
    *   **Test Details:** Component tests.
    *   **Status:** Not Started
*   **Task-014.2: Implement Backend Share Note API Endpoint**
    *   **Dependencies:** Task-012.1 (Note Schema, Note Shares Schema)
    *   **Implementation:** `POST /api/notes/{note_id}/share` endpoint. Takes `shared_with_user_id` and `permissions`. Creates a record in `note_shares`. Validates note ownership.
    *   **Verification:** API creates share record correctly.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-014.3: Integrate Share Note UI with Backend API**
    *   **Dependencies:** Task-014.1, Task-014.2
    *   **Implementation:** Frontend calls share API, handles response.
    *   **Verification:** Sharing a note 1-1 works end-to-end.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started
*   **Task-014.4: Modify Notes Listing API (US013) to include notes shared with the user**
    *   **Dependencies:** Task-013.2, Task-012.1
    *   **Implementation:** Update `GET /api/notes/private` (or create new `GET /api/notes/shared-with-me`) to include notes from `note_shares` where `shared_with_user_id` is current user.
    *   **Verification:** API correctly returns notes shared with the user.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started

**US015: Create a Group Note**
*   **As a** group member,
*   **I want to** create a note within a specific group I belong to,
*   **So that** all members of that group can view and collaborate on it.
*   **Acceptance Criteria:**
    1.  When in a group's context, a "Create Note" option is available.
    2.  User can input title and content.
    3.  The note is automatically associated with the current group and visible to all its members.
    4.  The created note appears in the group's note list.
*   **Test Scenarios:**
    1.  **TS015.1 (Happy Path):** Member creates note in group -> Note visible to all group members.
    2.  **TS015.2 (Non-Member Attempt):** UI/Logic prevents non-members from creating notes in a group.

**Task Breakdown for US015:**
*   **Task-015.1: Implement Backend Create Group Note API Endpoint**
    *   **Dependencies:** Task-012.1, Task-008.1
    *   **Implementation:** `POST /api/groups/{group_id}/notes` endpoint. Takes title, content. Sets `user_id` to current user, `group_id` to the specified group. Validates group membership.
    *   **Verification:** API creates group note, associated with the group.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-015.2: Integrate Note Editor UI (Task-012.2) for Group Notes**
    *   **Dependencies:** Task-012.2, Task-015.1
    *   **Implementation:** Adapt note editor UI or context to call the group note creation API when creating a note within a group.
    *   **Verification:** Group note creation works end-to-end.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started

**US016: View Group Notes**
*   **As a** group member,
*   **I want to** see a list of all notes shared within a group I belong to,
*   **So that** I can access and read them.
*   **Acceptance Criteria:**
    1.  Within a group's detail page, a section lists all notes belonging to that group.
    2.  List displays note titles, creator, and date.
    3.  Clicking a note opens it for viewing (and potential editing based on group rules later).
*   **Test Scenarios:**
    1.  **TS016.1 (View List):** Group member sees all notes for that group.
    2.  **TS016.2 (No Notes in Group):** Appropriate message shown.
    3.  **TS016.3 (Non-Member Access):** Non-member cannot see group notes (unless group is public - TBD).

**Task Breakdown for US016:**
*   **Task-016.1: Develop UI for Listing Group Notes**
    *   **Dependencies:** None
    *   **Implementation:** Frontend section within group view to list notes.
    *   **Verification:** UI renders correctly.
    *   **Test Details:** Component tests.
    *   **Status:** Not Started
*   **Task-016.2: Implement Backend API Endpoint for Listing Group Notes**
    *   **Dependencies:** Task-012.1, Task-008.1
    *   **Implementation:** `GET /api/groups/{group_id}/notes` endpoint. Returns notes where `group_id` matches. Validates user is member of the group.
    *   **Verification:** API returns correct notes for the group.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-016.3: Integrate Group Notes List UI with Backend API**
    *   **Dependencies:** Task-016.1, Task-016.2
    *   **Implementation:** Frontend calls API and displays the list of group notes.
    *   **Verification:** Group notes list is displayed correctly.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started

---

### Module 5: Bible Passage Integration

**US017: Embed Bible Passage in a Note**
*   **As a** user,
*   **When** creating or editing any note (private or group),
*   **I want to** search for and embed a Bible passage (e.g., John 3:16-18) by reference,
*   **So that** the text of the passage appears within my note for context and study.
*   **Acceptance Criteria:**
    1.  An interface within the note editor allows searching for Bible passages (Book, Chapter, Verse(s)).
    2.  User selects a passage.
    3.  The selected passage text is fetched from an external Bible API and displayed/inserted into the note content.
    4.  The reference (e.g., "John 3:16-18") and possibly the fetched text are stored with the note.
*   **Test Scenarios:**
    1.  **TS017.1 (Happy Path - Single Verse):** Search "John 3:16", select -> John 3:16 text embedded.
    2.  **TS017.2 (Happy Path - Verse Range):** Search "Romans 8:1-4", select -> Romans 8:1-4 text embedded.
    3.  **TS017.3 (Invalid Reference):** Search "Genesis 99:1" -> Error message, no passage embedded.
    4.  **TS017.4 (API Unavailable):** Bible API call fails -> Graceful error message.

**Task Breakdown for US017:**
*   **Task-017.1: Research and Select Bible API**
    *   **Dependencies:** None
    *   **Implementation:** Evaluate options (ESV API, NET Bible, Bible-api.com, etc.) based on features, terms, rate limits, translations.
    *   **Verification:** Decision made, API key obtained (if needed).
    *   **Test Details:** N/A
    *   **Status:** Not Started
*   **Task-017.2: Design Bible Passage Search/Selector UI**
    *   **Dependencies:** Task-012.2 (Note Editor)
    *   **Implementation:** Frontend component (modal or inline) for selecting Book, Chapter, Start Verse, End Verse.
    *   **Verification:** UI allows easy passage selection.
    *   **Test Details:** Component tests.
    *   **Status:** Not Started
*   **Task-017.3: Implement Backend Proxy/Wrapper for Bible API (Recommended)**
    *   **Dependencies:** Task-017.1
    *   **Implementation:** Create a backend endpoint (e.g., `GET /api/bible/passage?ref=<reference_string>`) that calls the chosen external Bible API. This hides API keys and can handle caching.
    *   **Verification:** Backend endpoint successfully fetches passage text from external API.
    *   **Test Details:** API unit/integration tests (mocking external API).
    *   **Status:** Not Started
*   **Task-017.4: Integrate Passage Selector UI with Backend Proxy and Note Editor**
    *   **Dependencies:** Task-017.2, Task-017.3, Task-012.2
    *   **Implementation:** Selector UI calls backend proxy. Fetched text is inserted into the note editor content. Store reference and/or text with the note.
    *   **Verification:** Bible passage can be embedded into a note.
    *   **Test Details:** E2E test for embedding a passage.
    *   **Status:** Not Started
*   **Task-017.5: Update Note Schema/Model for Bible Passages**
    *   **Dependencies:** Task-012.1
    *   **Implementation:** Decide how to store embedded passages. E.g., as special markdown/HTML, or as structured data within the note's `content` JSON, or in a separate `note_bible_passages` table linking passage reference to a note.
    *   **Verification:** Schema supports storing and retrieving embedded passage information.
    *   **Test Details:** Schema review.
    *   **Status:** Not Started

**US018: Highlight Text within an Embedded Bible Passage**
*   **As a** user,
*   **When** viewing a note containing an embedded Bible passage,
*   **I want to** select and highlight portions of the passage text,
*   **So that** I can emphasize key phrases or words for my study.
*   **Acceptance Criteria:**
    1.  User can select text within an embedded Bible passage in a note.
    2.  An option to "Highlight" the selection appears (e.g., context menu, toolbar).
    3.  The selected text is visually highlighted (e.g., yellow background).
    4.  Highlight information (passage reference, start/end offsets or selected text, color) is saved with the note.
    5.  Highlights persist and are displayed when the note is reopened.
*   **Test Scenarios:**
    1.  **TS018.1 (Happy Path):** Select text in passage, click highlight -> Text highlighted, saved. Reopen note -> Highlight visible.
    2.  **TS018.2 (Multiple Highlights):** User can create multiple distinct highlights in one passage.
    3.  **TS018.3 (Remove Highlight - Future):** (Optional) User can remove a highlight.

**Task Breakdown for US018:**
*   **Task-018.1: Design Data Structure for Highlights**
    *   **Dependencies:** Task-017.5
    *   **Implementation:** Define how highlight data is stored, linked to the note and the specific passage/verse. E.g., within note content JSON: `{ "passage_ref": "John 3:16", "highlights": [{"startOffset": 5, "endOffset": 10, "color": "yellow"}] }`.
    *   **Verification:** Data structure can accurately represent highlights.
    *   **Test Details:** Schema/structure review.
    *   **Status:** Not Started
*   **Task-018.2: Implement Frontend Highlighting Logic**
    *   **Dependencies:** Task-017.4 (Displaying embedded passages)
    *   **Implementation:** JavaScript to detect text selection within passages, display a highlight option, and apply visual styling.
    *   **Verification:** User can select and visually highlight text.
    *   **Test Details:** Component/integration tests for highlighting UI.
    *   **Status:** Not Started
*   **Task-018.3: Implement Save/Load Highlight Data**
    *   **Dependencies:** Task-018.1, Backend note save/load APIs (from US012, US015)
    *   **Implementation:**
        *   Frontend sends highlight data to backend when note is saved.
        *   Backend saves this data with the note.
        *   Frontend receives highlight data when note is loaded and re-applies highlights.
    *   **Verification:** Highlights persist across sessions.
    *   **Test Details:** E2E test for creating, saving, and viewing highlights.
    *   **Status:** Not Started

**US019: Add Comments to an Embedded Bible Passage/Highlight**
*   **As a** user,
*   **When** viewing a note containing an embedded Bible passage,
*   **I want to** add a textual comment anchored to a specific verse, range of verses, or a highlight I've made,
*   **So that** I can record my insights or questions related to that specific part of the text.
*   **Acceptance Criteria:**
    1.  User can select a verse, range, or existing highlight.
    2.  An option to "Add Comment" appears.
    3.  A text input field allows the user to write their comment.
    4.  The comment is saved and associated with the note and the specific anchor point (verse/highlight).
    5.  Comments are displayed near their anchor points when the note is viewed.
*   **Test Scenarios:**
    1.  **TS019.1 (Comment on Verse):** Select verse, add comment -> Comment saved and displayed.
    2.  **TS019.2 (Comment on Highlight):** Select highlight, add comment -> Comment saved and displayed near highlight.
    3.  **TS019.3 (Multiple Comments):** User can add multiple comments to different parts of a passage.

**Task Breakdown for US019:**
*   **Task-019.1: Design Data Structure for Comments on Passages/Highlights**
    *   **Dependencies:** Task-017.5, Task-018.1
    *   **Implementation:** Define how comment data is stored, linked to the note, passage reference, and specific anchor (verse ID, highlight ID, or text offsets). E.g., within note content JSON or a separate `note_passage_comments` table.
    *   **Verification:** Data structure can accurately represent comments and their anchors.
    *   **Test Details:** Schema/structure review.
    *   **Status:** Not Started
*   **Task-019.2: Implement Frontend Commenting UI**
    *   **Dependencies:** Task-017.4, Task-018.2
    *   **Implementation:** JavaScript to allow selection of anchor, display "Add Comment" option, show input field. Logic to display existing comments.
    *   **Verification:** User can initiate commenting and see existing comments.
    *   **Test Details:** Component/integration tests for commenting UI.
    *   **Status:** Not Started
*   **Task-019.3: Implement Save/Load Comment Data**
    *   **Dependencies:** Task-019.1, Backend note save/load APIs
    *   **Implementation:**
        *   Frontend sends comment data (text, anchor info) to backend when note is saved.
        *   Backend saves this data with the note.
        *   Frontend receives comment data when note is loaded and displays comments.
    *   **Verification:** Comments persist across sessions and are displayed correctly.
    *   **Test Details:** E2E test for adding, saving, and viewing comments.
    *   **Status:** Not Started

**US020: View Hebrew/Greek Translations for Embedded Passages**
*   **As a** user,
*   **When** viewing an embedded Bible passage in a note,
*   **I want to** be able to see the original Hebrew/Greek text for a selected verse or word (if available from the API),
*   **So that** I can gain deeper insight into the original language meaning.
*   **Acceptance Criteria:**
    1.  An option (e.g., toggle, button, hover) is available on a verse or selected word within an embedded passage.
    2.  Activating this option fetches and displays the corresponding Hebrew/Greek text from the Bible API.
    3.  The original language text is displayed clearly (e.g., in a tooltip, sidebar, or inline).
    4.  Graceful handling if original language text is not available for the selection.
*   **Test Scenarios:**
    1.  **TS020.1 (View Greek for NT Verse):** Select NT verse, activate translation -> Greek text shown.
    2.  **TS020.2 (View Hebrew for OT Verse):** Select OT verse, activate translation -> Hebrew text shown.
    3.  **TS020.3 (Translation Unavailable):** Select verse/word with no original language data -> "Translation not available" message.

**Task Breakdown for US020:**
*   **Task-020.1: Verify Bible API Support for Original Languages**
    *   **Dependencies:** Task-017.1
    *   **Implementation:** Confirm chosen Bible API can provide Hebrew/Greek texts and how to query them.
    *   **Verification:** API capabilities understood.
    *   **Test Details:** API documentation review, test API calls.
    *   **Status:** Not Started
*   **Task-020.2: Extend Backend Bible Proxy for Original Language Texts**
    *   **Dependencies:** Task-017.3, Task-020.1
    *   **Implementation:** Modify or add new endpoint(s) in the backend proxy to fetch Hebrew/Greek for a given reference/word.
    *   **Verification:** Backend can fetch original language texts.
    *   **Test Details:** API unit/integration tests.
    *   **Status:** Not Started
*   **Task-020.3: Implement Frontend UI for Requesting/Displaying Translations**
    *   **Dependencies:** Task-017.4
    *   **Implementation:** UI elements (e.g., icon on hover, context menu option) to trigger translation lookup. Display logic for showing the fetched original language text.
    *   **Verification:** User can trigger and view original language translations.
    *   **Test Details:** Component/integration tests.
    *   **Status:** Not Started
*   **Task-020.4: Integrate Frontend Translation UI with Backend Proxy**
    *   **Dependencies:** Task-020.2, Task-020.3
    *   **Implementation:** Frontend UI calls backend proxy to get original language text and displays it.
    *   **Verification:** Hebrew/Greek translations can be viewed for passages.
    *   **Test Details:** E2E test.
    *   **Status:** Not Started

## Verification Criteria (Overall Project)
*   All specified core features are implemented and functional as per their user story acceptance criteria.
*   Users can successfully sign up, log in, and log out.
*   Users can find, connect with, and disconnect from other users.
*   Users can create, join, invite to, and manage groups.
*   Users can create, view, and share private notes.
*   Users can create and view group notes within their respective groups.
*   Bible passages can be embedded, highlighted, commented upon, and original language translations viewed within notes.
*   The application functions correctly as a PWA (installable, basic offline access for app shell).
*   UI is responsive and user-friendly across common device sizes (desktop, tablet, mobile).
*   No critical bugs or regressions in existing functionality upon completion of each module.
*   Data integrity is maintained for users, connections, groups, notes, and Bible integrations.

## Potential Risks and Mitigations

1.  **Bible API Limitations/Costs/Terms of Service:**
    *   **Risk:** Selected API may have restrictive usage limits, unexpected costs, or terms that conflict with app features.
    *   **Mitigation:** Thoroughly research and test multiple Bible APIs early (Task-017.1). Design a backend proxy (Task-017.3) to allow easier switching if needed. Have a fallback plan for limited functionality if primary API fails.
2.  **Complexity of Rich Text Editor for Notes & Bible Integration:**
    *   **Risk:** Implementing highlighting, commenting, and embedding within a custom or third-party rich text editor can be highly complex and time-consuming.
    *   **Mitigation:** Start with basic text input for notes. For Bible integration, initially embed as simple text. Incrementally add rich text features and complex interactions. Evaluate off-the-shelf editors (e.g., Tiptap, Quill.js) for their extensibility.
3.  **Scalability of Real-time Features (if added later, e.g., collaborative editing):**
    *   **Risk:** If real-time collaboration is added, ensuring performance and data consistency can be challenging.
    *   **Mitigation:** For MVP, focus on non-real-time sharing. If real-time becomes a requirement, allocate specific R&D time and consider BaaS solutions (Firebase) that simplify this, or technologies like WebSockets and CRDTs.
4.  **PWA Offline Complexity:**
    *   **Risk:** Implementing robust offline data sync for notes and user data can be complex.
    *   **Mitigation:** For MVP PWA, focus on caching the app shell and read-only content. Implement more advanced offline data storage and sync strategies (e.g., IndexedDB, Background Sync API) in later iterations.
5.  **Scope Creep:**
    *   **Risk:** Adding too many features beyond the defined MVP can delay launch and reduce focus.
    *   **Mitigation:** Strictly adhere to the prioritized user stories for the initial versions. Maintain a backlog for future features and enhancements, to be addressed post-MVP based on user feedback.

## Alternative Approaches

1.  **Phased Rollout by Feature Set:** Focus on delivering one complete module (e.g., Auth + Private Notes + Basic Bible Embed) as a very lean MVP, then iteratively add Connections, then Groups, etc.
    *   *Pros:* Faster initial feedback, reduced initial complexity.
    *   *Cons:* Core value proposition might not be fully realized early on.
2.  **BaaS-Heavy Approach:** Maximize use of Firebase/Supabase for Auth, Database, Storage, and potentially Functions to reduce backend development time significantly.
    *   *Pros:* Rapid development, built-in scalability for many features.
    *   *Cons:* Vendor lock-in, potential cost implications at scale, less flexibility than custom backend.
3.  **Focus on Desktop Web First, PWA Later:** Build as a standard responsive web app first, then add PWA features (service worker, manifest).
    *   *Pros:* Simplifies initial development by deferring PWA-specific complexities.
    *   *Cons:* Delays PWA benefits like installability and enhanced offline capabilities.

*(This plan provides a comprehensive starting point. Each task would be further detailed during sprint planning if following an Agile methodology.)*
