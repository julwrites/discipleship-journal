---
id: FEATURES-20260131-163351-BCH
status: completed
title: Add unmark functionality to Bible reading plans: allow users to undo accidental day completions
priority: medium
created: 2026-01-31 16:33:51
category: features
dependencies: 
type: task
---

# Add unmark functionality to Bible reading plans: allow users to undo accidental day completions

## Task Information
- **Dependencies**: None (can be implemented independently)
- **Related tasks**: FEATURES-20251220-133329-KLU (reading plans backend)

## Task Details
Currently, Bible reading plans only allow users to mark days as complete, with no way to undo accidental completions. Once a day is marked, the UI shows a disabled checkmark button. This task adds the ability to unmark completed days.

### Current Implementation Analysis:
- **Database**: `user_reading_plan_progress` table stores `completed_at` timestamps
- **Backend**: `ReadingPlanService.MarkDayComplete()` only inserts records, no delete
- **API**: `POST /api/my-reading-plans/{id}/progress` marks days complete
- **Frontend**: `ReadingPlanDetail.tsx` shows disabled checkmark for completed days
- **Constraint**: `UNIQUE(user_reading_plan_id, day_number)` prevents duplicates

### Steps:

#### 1. Backend Changes (`api/services/reading_plan_service.go`):
- Add `UnmarkDayComplete(userReadingPlanID string, dayNumber int) error` method
- Implement SQL `DELETE FROM user_reading_plan_progress WHERE user_reading_plan_id = $1 AND day_number = $2`
- Update `ReadingPlanService` interface to include unmark method

#### 2. API Changes (`api/handlers/reading_plan.go`):
- Add `DELETE /api/my-reading-plans/{id}/progress/{day}` endpoint
- Validate day number range (1 to plan total days)
- Check user owns the reading plan subscription
- Call `UnmarkDayComplete` service method

#### 3. Frontend Changes (`web/src/pages/ReadingPlanDetail.tsx`):
- Update completed day UI: change from disabled checkmark to toggle button
- Add "Undo" or "Mark Incomplete" button for completed days
- Implement optimistic updates: immediate UI change, then API call
- Handle API errors and revert UI if unmark fails
- Update progress calculation after unmarking

#### 4. UI/UX Considerations:
- Visual design: Differentiate between "Mark Complete" and "Undo" states
- Confirmation: Optional confirmation dialog for unmarking
- Accessibility: Clear button labels, ARIA attributes
- Mobile: Ensure touch targets are adequate

#### 5. Testing:
- Test unmarking recently completed days
- Test unmarking days marked long ago
- Verify progress percentage updates correctly
- Test error scenarios (network errors, unauthorized access)
- Ensure no duplicate completion records after unmark/remark

### Acceptance Criteria
- [ ] `UnmarkDayComplete` method implemented in service
- [ ] `DELETE /api/my-reading-plans/{id}/progress/{day}` endpoint works
- [ ] Frontend shows toggle button for completed days (not just disabled checkmark)
- [ ] Users can successfully unmark accidentally completed days
- [ ] Progress percentage updates correctly after unmarking
- [ ] No regression in existing mark completion functionality
- [ ] Error handling: failed unmarks revert UI state
- [ ] Accessible UI with clear button labels

## Implementation Status
### Completed Work
- ✅ Codebase exploration completed
- ✅ Current limitations identified
- ✅ Backend: Implemented `UnmarkDayComplete` in Service and Handler
- ✅ Backend: Added API endpoint `DELETE /api/my-reading-plans/{id}/progress/{day}`
- ✅ Backend: Updated contract tests and mock tests
- ✅ Frontend: Added `unmarkPlanDayComplete` to API service
- ✅ Frontend: Updated `ReadingPlanDetail.tsx` to allow unmarking completed days

### Blockers
None.

## Notes
- **Data integrity**: Unmarking should delete the progress record, not just set `completed_at` to NULL
- **UI pattern**: Consider using toggle button that shows "Mark Complete" / "Mark Incomplete"
- **Performance**: Simple DELETE operation, minimal performance impact
- **Edge cases**: Handle concurrent mark/unmark operations
- **Backward compatibility**: Existing completed days remain in database
