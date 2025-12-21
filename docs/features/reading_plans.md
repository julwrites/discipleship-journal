# Bible Reading Plans

## Overview
The Bible Reading Plans feature allows users to follow structured reading schedules to engage with Scripture systematically. Users can browse available plans, subscribe to them, and track their daily progress.

## User Stories
- As a user, I want to browse available reading plans so I can choose one that fits my goals.
- As a user, I want to subscribe to a reading plan so I can track my progress.
- As a user, I want to see my active reading plans to quickly access them.
- As a user, I want to mark a day's reading as complete to track my progress.
- As a user, I want to see which days I have completed.

## API Endpoints

### Public/Browse
- `GET /api/reading-plans`: List all available reading plans.
- `GET /api/reading-plans/{id}`: Get details of a specific reading plan, including the schedule of readings.

### User/Subscription
- `POST /api/reading-plans/{id}/subscribe`: Subscribe the current user to a plan.
- `GET /api/my-reading-plans`: List reading plans the user is subscribed to, including status and start date.
- `GET /api/my-reading-plans/{id}/progress`: Get the list of completed days for a subscribed plan.
- `POST /api/my-reading-plans/{id}/progress`: Mark a specific day number as complete.

## Data Model

### ReadingPlan
- `ID`: UUID
- `Title`: String
- `Description`: String
- `Days`: Integer (Total number of days)

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
- **ReadingPlanDetail**: Detailed view of a subscribed plan. Lists all days with their passages and allows marking them as complete.
- **Service Layer**: `web/src/services/api.ts` handles all API communication.

## Future Enhancements
- Reminder notifications.
- Ability to restart a plan.
- "Catch up" feature to adjust schedule.
