package services

import (
	"context"
	"testing"
	"time"

	"database/sql"
	"discipleship_journal_api/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPlans(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)

	rows := mock.NewRows([]string{"id", "title", "description", "days", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Plan 1", "Desc 1", 30, time.Now(), time.Now()).
		AddRow(uuid.New(), "Plan 2", "Desc 2", 60, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT id, title, description, days, created_at, updated_at\s+FROM reading_plans`).
		WillReturnRows(rows)

	plans, err := service.GetAllPlans(context.Background())
	assert.NoError(t, err)
	assert.Len(t, plans, 2)
	assert.Equal(t, "Plan 1", plans[0].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkDayComplete_CheckCompletionError_GetDays(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()
	userPlanID := uuid.New()
	dayNumber := 1

	// Mock find active plan
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(userPlanID))

	// Mock Insert Progress
	mock.ExpectExec(`INSERT IGNORE INTO user_reading_plan_progress`).
		WithArgs(userPlanID, dayNumber, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect check completion queries
	// 1. Get total days - Error
	mock.ExpectQuery(`SELECT days FROM reading_plans`).
		WithArgs(planID).
		WillReturnError(assert.AnError)

	err = service.MarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkDayComplete_CheckCompletionError_GetCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()
	userPlanID := uuid.New()
	dayNumber := 1

	// Mock find active plan
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(userPlanID))

	// Mock Insert Progress
	mock.ExpectExec(`INSERT IGNORE INTO user_reading_plan_progress`).
		WithArgs(userPlanID, dayNumber, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect check completion queries
	// 1. Get total days
	mock.ExpectQuery(`SELECT days FROM reading_plans`).
		WithArgs(planID).
		WillReturnRows(mock.NewRows([]string{"days"}).AddRow(30))

	// 2. Get completed count - Error
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM user_reading_plan_progress`).
		WithArgs(userPlanID).
		WillReturnError(assert.AnError)

	err = service.MarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkDayComplete_CheckCompletionError_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()
	userPlanID := uuid.New()
	dayNumber := 30

	// Mock find active plan
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(userPlanID))

	// Mock Insert Progress
	mock.ExpectExec(`INSERT IGNORE INTO user_reading_plan_progress`).
		WithArgs(userPlanID, dayNumber, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect check completion queries
	// 1. Get total days
	mock.ExpectQuery(`SELECT days FROM reading_plans`).
		WithArgs(planID).
		WillReturnRows(mock.NewRows([]string{"days"}).AddRow(30))

	// 2. Get completed count
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM user_reading_plan_progress`).
		WithArgs(userPlanID).
		WillReturnRows(mock.NewRows([]string{"count"}).AddRow(30))

	// 3. Update status - Error
	mock.ExpectExec(`UPDATE user_reading_plans`).
		WithArgs(userPlanID).
		WillReturnError(assert.AnError)

	err = service.MarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlan_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	id := uuid.New()

	mock.ExpectQuery(`SELECT id, title, description, days, created_at, updated_at\s+FROM reading_plans\s+WHERE id = \?`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	plan, err := service.GetPlan(context.Background(), id)
	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Equal(t, models.ErrNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlan_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	id := uuid.New()

	mock.ExpectQuery(`SELECT id, title, description, days, created_at, updated_at\s+FROM reading_plans\s+WHERE id = \?`).
		WithArgs(id).
		WillReturnError(assert.AnError)

	plan, err := service.GetPlan(context.Background(), id)
	assert.Error(t, err)
	assert.Nil(t, plan)
	assert.Equal(t, assert.AnError, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlanDays_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	planID := uuid.New()

	mock.ExpectQuery(`SELECT id, reading_plan_id, day_number, passage, created_at\s+FROM reading_plan_days\s+WHERE reading_plan_id = \?`).
		WithArgs(planID).
		WillReturnError(assert.AnError)

	days, err := service.GetPlanDays(context.Background(), planID)
	assert.Error(t, err)
	assert.Nil(t, days)
	assert.Equal(t, assert.AnError, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscribe_InsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()

	// Mock check if exists - Returns no rows (so proceed)
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnError(sql.ErrNoRows)

	// Mock Insert - Returns error
	mock.ExpectExec(`INSERT INTO user_reading_plans`).
		WithArgs(sqlmock.AnyArg(), userID, planID, sqlmock.AnyArg(), "active").
		WillReturnError(assert.AnError)

	userPlan, err := service.Subscribe(context.Background(), userID, planID)
	assert.Error(t, err)
	assert.Nil(t, userPlan)
	assert.Equal(t, assert.AnError, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserPlans_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT u.id, u.user_id, u.reading_plan_id, u.start_date, u.status, u.created_at, u.updated_at,\s+p.id, p.title, p.description, p.days, p.created_at, p.updated_at\s+FROM user_reading_plans u\s+JOIN reading_plans p ON u.reading_plan_id = p.id\s+WHERE u.user_id = \?`).
		WithArgs(userID).
		WillReturnError(assert.AnError)

	plans, err := service.GetUserPlans(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, plans)
	assert.Equal(t, assert.AnError, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkDayComplete_Completed(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()
	userPlanID := uuid.New()
	dayNumber := 30

	// Mock find active plan
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(userPlanID))

	// Mock Insert Progress
	mock.ExpectExec(`INSERT IGNORE INTO user_reading_plan_progress`).
		WithArgs(userPlanID, dayNumber, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect check completion queries
	// 1. Get total days
	mock.ExpectQuery(`SELECT days FROM reading_plans`).
		WithArgs(planID).
		WillReturnRows(mock.NewRows([]string{"days"}).AddRow(30))

	// 2. Get completed count
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM user_reading_plan_progress`).
		WithArgs(userPlanID).
		WillReturnRows(mock.NewRows([]string{"count"}).AddRow(30)) // 30 completed, total 30 -> complete

	// 3. Update status to completed
	mock.ExpectExec(`UPDATE user_reading_plans`).
		WithArgs(userPlanID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = service.MarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlanProgress(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()
	userPlanID := uuid.New()

	// Mock find active plan
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(userPlanID))

	// Mock get progress
	mock.ExpectQuery(`SELECT day_number FROM user_reading_plan_progress`).
		WithArgs(userPlanID).
		WillReturnRows(mock.NewRows([]string{"day_number"}).AddRow(1).AddRow(2).AddRow(5))

	days, err := service.GetPlanProgress(context.Background(), userID, planID)
	assert.NoError(t, err)
	assert.Len(t, days, 3)
	assert.Equal(t, []int{1, 2, 5}, days)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlanProgress_PlanNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()

	// Mock find active plan - Not found
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnError(sql.ErrNoRows)

	_, err = service.GetPlanProgress(context.Background(), userID, planID)
	assert.Error(t, err)
	assert.Equal(t, models.ErrNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkDayComplete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()
	dayNumber := 1

	// Mock find active plan - Not Found
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnError(sql.ErrNoRows)

	err = service.MarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.Error(t, err)
	assert.Equal(t, models.ErrNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	id := uuid.New()

	rows := mock.NewRows([]string{"id", "title", "description", "days", "created_at", "updated_at"}).
		AddRow(id, "Plan 1", "Desc 1", 30, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT id, title, description, days, created_at, updated_at\s+FROM reading_plans\s+WHERE id = \?`).
		WithArgs(id).
		WillReturnRows(rows)

	plan, err := service.GetPlan(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, plan.ID)
	assert.Equal(t, "Plan 1", plan.Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlanDays(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	planID := uuid.New()

	rows := mock.NewRows([]string{"id", "reading_plan_id", "day_number", "passage", "created_at"}).
		AddRow(uuid.New(), planID, 1, "Gen 1", time.Now()).
		AddRow(uuid.New(), planID, 2, "Gen 2", time.Now())

	mock.ExpectQuery(`SELECT id, reading_plan_id, day_number, passage, created_at\s+FROM reading_plan_days\s+WHERE reading_plan_id = \?`).
		WithArgs(planID).
		WillReturnRows(rows)

	days, err := service.GetPlanDays(context.Background(), planID)
	assert.NoError(t, err)
	assert.Len(t, days, 2)
	assert.Equal(t, "Gen 1", days[0].Passage)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscribe(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()

	// Mock check if exists
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnError(sql.ErrNoRows) // Not found, so proceed

	// Mock Insert
	mock.ExpectExec(`INSERT INTO user_reading_plans`).
		WithArgs(sqlmock.AnyArg(), userID, planID, sqlmock.AnyArg(), "active").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(`SELECT created_at, updated_at FROM user_reading_plans WHERE id = \?`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).AddRow(time.Now(), time.Now()))

	userPlan, err := service.Subscribe(context.Background(), userID, planID)
	assert.NoError(t, err)
	assert.Equal(t, userID, userPlan.UserID)
	assert.Equal(t, "active", userPlan.Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscribe_AlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()

	// Mock check if exists - Returns a row
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(uuid.New()))

	_, err = service.Subscribe(context.Background(), userID, planID)
	assert.Error(t, err)
	assert.Equal(t, models.ErrAlreadyExists, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkDayComplete(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()
	userPlanID := uuid.New()
	dayNumber := 1

	// Mock find active plan
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(userPlanID))

	// Mock Insert Progress
	mock.ExpectExec(`INSERT IGNORE INTO user_reading_plan_progress`).
		WithArgs(userPlanID, dayNumber, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect check completion queries
	// 1. Get total days
	mock.ExpectQuery(`SELECT days FROM reading_plans`).
		WithArgs(planID).
		WillReturnRows(mock.NewRows([]string{"days"}).AddRow(30))

	// 2. Get completed count
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM user_reading_plan_progress`).
		WithArgs(userPlanID).
		WillReturnRows(mock.NewRows([]string{"count"}).AddRow(1)) // 1 completed, total 30 -> not complete

	err = service.MarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnsubscribe(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()

	// Mock Delete
	mock.ExpectExec(`DELETE FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = service.Unsubscribe(context.Background(), userID, planID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnsubscribe_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()

	// Mock Delete - No rows affected
	mock.ExpectExec(`DELETE FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err = service.Unsubscribe(context.Background(), userID, planID)
	assert.Error(t, err)
	assert.Equal(t, models.ErrNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserPlans(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()

	rows := mock.NewRows([]string{
		"id", "user_id", "reading_plan_id", "start_date", "status", "created_at", "updated_at",
		"p_id", "title", "description", "days", "p_created_at", "p_updated_at",
	}).
		AddRow(
			uuid.New(), userID, planID, time.Now(), "active", time.Now(), time.Now(),
			planID, "Plan Title", "Plan Desc", 30, time.Now(), time.Now(),
		)

	mock.ExpectQuery(`SELECT u.id, u.user_id, u.reading_plan_id, u.start_date, u.status, u.created_at, u.updated_at,\s+p.id, p.title, p.description, p.days, p.created_at, p.updated_at\s+FROM user_reading_plans u\s+JOIN reading_plans p ON u.reading_plan_id = p.id\s+WHERE u.user_id = \?`).
		WithArgs(userID).
		WillReturnRows(rows)

	plans, err := service.GetUserPlans(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, plans, 1)
	assert.Equal(t, "Plan Title", plans[0].Plan.Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnmarkDayComplete(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()
	userPlanID := uuid.New()
	dayNumber := 1

	// Mock find active plan
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(userPlanID))

	// Mock Delete Progress
	mock.ExpectExec(`DELETE FROM user_reading_plan_progress`).
		WithArgs(userPlanID, dayNumber).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = service.UnmarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnmarkDayComplete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	assert.NoError(t, err)
	defer db.Close()

	service := NewReadingPlanService(db)
	userID := uuid.New()
	planID := uuid.New()
	dayNumber := 1

	// Mock find active plan - Not Found
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnError(sql.ErrNoRows)

	err = service.UnmarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.Error(t, err)
	assert.Equal(t, models.ErrNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
