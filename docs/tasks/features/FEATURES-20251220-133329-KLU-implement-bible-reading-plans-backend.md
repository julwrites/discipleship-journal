---
id: FEATURES-20251220-133329-KLU
status: completed
title: Implement Bible Reading Plans Backend
priority: high
created: 2025-12-20 13:33:29
category: features
dependencies: FOUNDATION-004
type: task
parent: FEATURES-20251213-092127-YFS
---

# Implement Bible Reading Plans Backend

Implement the backend support for Bible Reading Plans, including database schema, service layer, and API endpoints.

## Todo List
- [x] Create database migrations for `reading_plans`, `reading_plan_days`, `user_reading_plans`, `user_reading_plan_progress`
- [x] Implement `ReadingPlanService`
- [x] Implement `ReadingPlanHandler`
- [x] Register API routes
- [x] Add unit tests for Service and Handler

## Detailed Requirements

### Database Schema
*   `reading_plans`:
    *   `id` (UUID, PK)
    *   `title` (TEXT)
    *   `description` (TEXT)
    *   `days` (INT) - Total number of days
    *   `created_at`, `updated_at`
*   `reading_plan_days`:
    *   `id` (UUID, PK)
    *   `reading_plan_id` (UUID, FK)
    *   `day_number` (INT)
    *   `passage` (TEXT) - e.g. "Genesis 1-3"
    *   `created_at`
*   `user_reading_plans`:
    *   `id` (UUID, PK)
    *   `user_id` (UUID, FK)
    *   `reading_plan_id` (UUID, FK)
    *   `start_date` (TIMESTAMP)
    *   `status` (VARCHAR) - active, completed
    *   `created_at`, `updated_at`
*   `user_reading_plan_progress`:
    *   `id` (UUID, PK)
    *   `user_reading_plan_id` (UUID, FK)
    *   `day_number` (INT)
    *   `completed_at` (TIMESTAMP)
    *   `created_at`

### API Endpoints
*   `GET /api/reading-plans` - List available plans
*   `GET /api/reading-plans/{id}` - Get plan details (with days)
*   `POST /api/reading-plans/{id}/subscribe` - Subscribe to a plan
*   `GET /api/my-reading-plans` - List user's plans
*   `POST /api/my-reading-plans/{id}/progress` - Mark day as complete
