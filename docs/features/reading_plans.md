# Bible Reading Plans

## Overview
The Bible Reading Plans feature allows users to follow structured reading schedules to engage with Scripture systematically. Users can browse available plans, subscribe to them, and track their daily progress.

### Available Plans
- **Calendar-based**: Standard 365-day plans (e.g., M'Cheyne, Chronological).
- **Discipleship Journal**: A 4-stream approach to reading through the Bible in a year (25 days/month).
- **Custom/Algorithmic**: Plans generated algorithmically (e.g., Bible in 90 Days).

## User Stories
- As a user, I want to browse available reading plans so I can choose one that fits my goals.
- As a user, I want to subscribe to a reading plan so I can track my progress.
- As a user, I want to see my active reading plans to quickly access them.
- As a user, I want to mark a day's reading as complete to track my progress.
- As a user, I want to unmark a day if I marked it by mistake.
- As a user, I want to read the Bible passage directly within the app.
- As a user, I want to create a journal note linked to the day's reading to capture my reflections.
- As a user, I want to see which days I have completed.

## API Endpoints

### Public/Browse
- `GET /api/reading-plans`: List all available reading plans. Returns a list of plans where `days` is the total number of days (integer).
- `GET /api/reading-plans/{id}`: Get details of a specific reading plan. Response includes full plan details and `days` (array of `ReadingPlanDay` objects with `day_number` and `passage`).

### User/Subscription
- `POST /api/reading-plans/{id}/subscribe`: Subscribe the current user to a plan.
- `GET /api/my-reading-plans`: List reading plans the user is subscribed to, including status and start date.
- `GET /api/my-reading-plans/{id}/progress`: Get the list of completed days for a subscribed plan.
- `POST /api/my-reading-plans/{id}/progress`: Mark a specific day number as complete.
- `DELETE /api/my-reading-plans/{id}/progress/{day_number}`: Unmark a specific day number (undo completion).

## Data Model

### ReadingPlan
- `ID`: UUID
- `Title`: String
- `Description`: String
- `Days`: Integer (Total number of days) in list view, or Array of `ReadingPlanDay` in detail view.
- `PlanType`: String (e.g., 'calendar', 'algorithmic')

### ReadingPlanDay
- `ID`: UUID
- `ReadingPlanID`: UUID
- `DayNumber`: Integer
- `Passage`: String (Bible Reference)

### UserReadingPlan
- `ID`: UUID
- `UserID`: UUID
- `ReadingPlanID`: UUID
- `StartDate`: Timestamp
- `Status`: Enum (active, completed)

### UserReadingPlanProgress
- `ID`: UUID
- `UserReadingPlanID`: UUID
- `DayNumber`: Integer
- `CompletedAt`: Timestamp

## Frontend Implementation
- **ReadingPlansPage**: Main entry point. Contains tabs for "My Plans" and "Browse Plans".
- **ReadingPlanDetail**: Detailed view of a subscribed plan.
  - Lists all days with their passages.
  - **Mark/Unmark**: Users can toggle completion status of each day.
  - **Read Passage**: Clicking a reference opens the `BiblePassageDialog` to read the text in-app (supports version selection).
  - **Create Note**: A button on each day allows users to create a new note pre-populated with the passage reference and context tags.
- **Service Layer**: `web/src/services/api.ts` handles all API communication.

## Future Enhancements
- Reminder notifications.
- Ability to restart a plan.
- "Catch up" feature to adjust schedule.
- Group reading plans (shared progress).
