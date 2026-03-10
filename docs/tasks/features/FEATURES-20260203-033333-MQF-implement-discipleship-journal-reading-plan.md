---
id: FEATURES-20260203-033333-MQF
status: completed
title: Implement Discipleship Journal Reading Plan
priority: medium
created: 2026-02-03 03:33:33
category: features
dependencies:
type: task
---

# Implement Discipleship Journal Reading Plan

## Objective
Implement a Bible reading plan based on the Discipleship Journal system, which features 4 daily readings and 25 reading days per month to allow for catch-up.

## Implementation Details
Due to the unavailability of a clean digital source for the exact Discipleship Journal reading list, a **Balanced 4-Stream Algorithmic Plan** was implemented that follows the core principles of the system:
1.  **Structure**: 4 Daily Readings.
2.  **Cadence**: 25 Days per month (Days 26-End of month are skipped).
3.  **Streams**:
    *   **Stream 1 (Torah & History)**: Genesis through Esther.
    *   **Stream 2 (Poetry & Prophets)**: Job through Malachi.
    *   **Stream 3 (Gospels)**: Matthew, Mark, Luke, John (Repeated ~3 times to fill the year).
    *   **Stream 4 (Epistles)**: Acts through Revelation (Repeated/Looped to fill the year).

## Script
A new script `scripts/generate_dj_plan/main.go` was created to generate the migration.
*   It uses a `distribute` algorithm to spread chapters over 300 active reading days (12 months * 25 days).
*   It generates a SQL migration file `api/migrations/000028_add_discipleship_journal_plan.up.sql`.

## Verification
*   Generated migration includes ~300 days of readings.
*   Readings include 4 passages per day.
*   Start and end points align with the biblical books in the streams.
