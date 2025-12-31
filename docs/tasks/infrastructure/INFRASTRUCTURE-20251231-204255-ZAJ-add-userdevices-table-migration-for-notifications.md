---
id: INFRASTRUCTURE-20251231-204255-ZAJ
status: deferred
title: Add user_devices table migration for notifications
priority: low
created: 2025-12-31 20:42:55
category: infrastructure
dependencies:
type: task
---

# Add user_devices table migration for notifications

## Status: Deferred
This task has been deferred to prioritize core functionality fixes. Notifications are not critical for MVP.

## Problem
The notification service references a `user_devices` table that doesn't exist in migrations, making push notifications non-functional.

## Requirements
1. Create migration for `user_devices` table
2. Update notification service to handle missing table gracefully
3. Test notification registration and sending

## Steps
1. Create migration file for `user_devices` table
2. Run migration in development
3. Test notification registration
4. Test notification sending
5. Deploy to production

## Acceptance Criteria
- [ ] `user_devices` table migration created
- [ ] Notification service works with real table
- [ ] Notifications can be registered and sent
- [ ] Graceful handling of missing FCM tokens
