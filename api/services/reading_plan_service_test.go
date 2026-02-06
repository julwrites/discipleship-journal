package services

import (
	"context"
	"testing"
	"time"

	"discipleship_journal_api/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPlans(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)

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

func TestGetPlanProgress(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
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
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
	userID := uuid.New()
	planID := uuid.New()

	// Mock find active plan - Not found
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnError(pgx.ErrNoRows)

	_, err = service.GetPlanProgress(context.Background(), userID, planID)
	assert.Error(t, err)
	assert.Equal(t, models.ErrNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkDayComplete_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
	userID := uuid.New()
	planID := uuid.New()
	dayNumber := 1

	// Mock find active plan - Not Found
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnError(pgx.ErrNoRows)

	err = service.MarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.Error(t, err)
	assert.Equal(t, models.ErrNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlan(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
	id := uuid.New()

	rows := mock.NewRows([]string{"id", "title", "description", "days", "created_at", "updated_at"}).
		AddRow(id, "Plan 1", "Desc 1", 30, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT id, title, description, days, created_at, updated_at\s+FROM reading_plans\s+WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(rows)

	plan, err := service.GetPlan(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, plan.ID)
	assert.Equal(t, "Plan 1", plan.Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlanDays(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
	planID := uuid.New()

	rows := mock.NewRows([]string{"id", "reading_plan_id", "day_number", "passage", "created_at"}).
		AddRow(uuid.New(), planID, 1, "Gen 1", time.Now()).
		AddRow(uuid.New(), planID, 2, "Gen 2", time.Now())

	mock.ExpectQuery(`SELECT id, reading_plan_id, day_number, passage, created_at\s+FROM reading_plan_days\s+WHERE reading_plan_id = \$1`).
		WithArgs(planID).
		WillReturnRows(rows)

	days, err := service.GetPlanDays(context.Background(), planID)
	assert.NoError(t, err)
	assert.Len(t, days, 2)
	assert.Equal(t, "Gen 1", days[0].Passage)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscribe(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
	userID := uuid.New()
	planID := uuid.New()

	// Mock check if exists
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnError(pgx.ErrNoRows) // Not found, so proceed

	// Mock Insert
	rows := mock.NewRows([]string{"id", "user_id", "reading_plan_id", "start_date", "status", "created_at", "updated_at"}).
		AddRow(uuid.New(), userID, planID, time.Now(), "active", time.Now(), time.Now())

	mock.ExpectQuery(`INSERT INTO user_reading_plans`).
		WithArgs(userID, planID, pgxmock.AnyArg()).
		WillReturnRows(rows)

	userPlan, err := service.Subscribe(context.Background(), userID, planID)
	assert.NoError(t, err)
	assert.Equal(t, userID, userPlan.UserID)
	assert.Equal(t, "active", userPlan.Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscribe_AlreadyExists(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
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
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
	userID := uuid.New()
	planID := uuid.New()
	userPlanID := uuid.New()
	dayNumber := 1

	// Mock find active plan
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(userPlanID))

	// Mock Insert Progress
	mock.ExpectExec(`INSERT INTO user_reading_plan_progress`).
		WithArgs(userPlanID, dayNumber, pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

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
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
	userID := uuid.New()
	planID := uuid.New()

	// Mock Delete
	mock.ExpectExec(`DELETE FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = service.Unsubscribe(context.Background(), userID, planID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnsubscribe_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
	userID := uuid.New()
	planID := uuid.New()

	// Mock Delete - No rows affected
	mock.ExpectExec(`DELETE FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = service.Unsubscribe(context.Background(), userID, planID)
	assert.Error(t, err)
	assert.Equal(t, models.ErrNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserPlans(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
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

	mock.ExpectQuery(`SELECT u.id, u.user_id, u.reading_plan_id, u.start_date, u.status, u.created_at, u.updated_at,\s+p.id, p.title, p.description, p.days, p.created_at, p.updated_at\s+FROM user_reading_plans u\s+JOIN reading_plans p ON u.reading_plan_id = p.id\s+WHERE u.user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	plans, err := service.GetUserPlans(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, plans, 1)
	assert.Equal(t, "Plan Title", plans[0].Plan.Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnmarkDayComplete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
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
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = service.UnmarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnmarkDayComplete_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	service := NewReadingPlanService(mock)
	userID := uuid.New()
	planID := uuid.New()
	dayNumber := 1

	// Mock find active plan - Not Found
	mock.ExpectQuery(`SELECT id FROM user_reading_plans`).
		WithArgs(userID, planID).
		WillReturnError(pgx.ErrNoRows)

	err = service.UnmarkDayComplete(context.Background(), userID, planID, dayNumber)
	assert.Error(t, err)
	assert.Equal(t, models.ErrNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
