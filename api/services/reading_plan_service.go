package services

import (
	"context"
	"time"

	"discipleship_journal_api/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ReadingPlanService interface {
	GetAllPlans(ctx context.Context) ([]*models.ReadingPlan, error)
	GetPlan(ctx context.Context, id uuid.UUID) (*models.ReadingPlan, error)
	GetPlanDays(ctx context.Context, planID uuid.UUID) ([]*models.ReadingPlanDay, error)
	Subscribe(ctx context.Context, userID, planID uuid.UUID) (*models.UserReadingPlan, error)
	GetUserPlans(ctx context.Context, userID uuid.UUID) ([]*models.UserReadingPlan, error)
	MarkDayComplete(ctx context.Context, userID, planID uuid.UUID, dayNumber int) error
	UnmarkDayComplete(ctx context.Context, userID, planID uuid.UUID, dayNumber int) error
	GetPlanProgress(ctx context.Context, userID, planID uuid.UUID) ([]int, error)
	Unsubscribe(ctx context.Context, userID, planID uuid.UUID) error
}

type readingPlanService struct {
	db DBInterfaceWithQuery
}

type DBInterfaceWithQuery interface {
	DBInterface
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func NewReadingPlanService(db DBInterfaceWithQuery) ReadingPlanService {
	return &readingPlanService{db: db}
}

func (s *readingPlanService) GetAllPlans(ctx context.Context) ([]*models.ReadingPlan, error) {
	query := `
		SELECT id, title, description, days, created_at, updated_at
		FROM reading_plans
		ORDER BY created_at DESC
	`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*models.ReadingPlan
	for rows.Next() {
		var p models.ReadingPlan
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Days, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, &p)
	}
	return plans, nil
}

func (s *readingPlanService) GetPlan(ctx context.Context, id uuid.UUID) (*models.ReadingPlan, error) {
	query := `
		SELECT id, title, description, days, created_at, updated_at
		FROM reading_plans
		WHERE id = $1
	`
	var p models.ReadingPlan
	err := s.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Title, &p.Description, &p.Days, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (s *readingPlanService) GetPlanDays(ctx context.Context, planID uuid.UUID) ([]*models.ReadingPlanDay, error) {
	query := `
		SELECT id, reading_plan_id, day_number, passage, created_at
		FROM reading_plan_days
		WHERE reading_plan_id = $1
		ORDER BY day_number ASC
	`
	rows, err := s.db.Query(ctx, query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []*models.ReadingPlanDay
	for rows.Next() {
		var d models.ReadingPlanDay
		if err := rows.Scan(&d.ID, &d.ReadingPlanID, &d.DayNumber, &d.Passage, &d.CreatedAt); err != nil {
			return nil, err
		}
		days = append(days, &d)
	}
	return days, nil
}

func (s *readingPlanService) Subscribe(ctx context.Context, userID, planID uuid.UUID) (*models.UserReadingPlan, error) {
	// Check if already subscribed
	checkQuery := `SELECT id FROM user_reading_plans WHERE user_id = $1 AND reading_plan_id = $2 AND status = 'active'`
	var existingID uuid.UUID
	err := s.db.QueryRow(ctx, checkQuery, userID, planID).Scan(&existingID)
	if err == nil {
		// Already exists
		return nil, models.ErrAlreadyExists
	}

	query := `
		INSERT INTO user_reading_plans (user_id, reading_plan_id, start_date, status)
		VALUES ($1, $2, $3, 'active')
		RETURNING id, user_id, reading_plan_id, start_date, status, created_at, updated_at
	`
	var p models.UserReadingPlan
	err = s.db.QueryRow(ctx, query, userID, planID, time.Now()).Scan(
		&p.ID, &p.UserID, &p.ReadingPlanID, &p.StartDate, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *readingPlanService) GetUserPlans(ctx context.Context, userID uuid.UUID) ([]*models.UserReadingPlan, error) {
	query := `
		SELECT u.id, u.user_id, u.reading_plan_id, u.start_date, u.status, u.created_at, u.updated_at,
		       p.id, p.title, p.description, p.days, p.created_at, p.updated_at
		FROM user_reading_plans u
		JOIN reading_plans p ON u.reading_plan_id = p.id
		WHERE u.user_id = $1
		ORDER BY u.start_date DESC
	`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*models.UserReadingPlan
	for rows.Next() {
		var u models.UserReadingPlan
		var p models.ReadingPlan
		if err := rows.Scan(
			&u.ID, &u.UserID, &u.ReadingPlanID, &u.StartDate, &u.Status, &u.CreatedAt, &u.UpdatedAt,
			&p.ID, &p.Title, &p.Description, &p.Days, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		u.Plan = &p
		plans = append(plans, &u)
	}
	return plans, nil
}

func (s *readingPlanService) MarkDayComplete(ctx context.Context, userID, planID uuid.UUID, dayNumber int) error {
	// First find the user's active plan
	var userPlanID uuid.UUID
	err := s.db.QueryRow(ctx, `
		SELECT id FROM user_reading_plans
		WHERE user_id = $1 AND reading_plan_id = $2 AND status = 'active'
	`, userID, planID).Scan(&userPlanID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.ErrNotFound
		}
		return err
	}

	// Insert progress
	query := `
		INSERT INTO user_reading_plan_progress (user_reading_plan_id, day_number, completed_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_reading_plan_id, day_number) DO NOTHING
	`
	_, err = s.db.Exec(ctx, query, userPlanID, dayNumber, time.Now())
	if err != nil {
		return err
	}

	// Check for plan completion
	return s.checkAndMarkPlanCompleted(ctx, userPlanID, planID)
}

func (s *readingPlanService) checkAndMarkPlanCompleted(ctx context.Context, userPlanID, planID uuid.UUID) error {
	// Get total days in plan
	var totalDays int
	err := s.db.QueryRow(ctx, "SELECT days FROM reading_plans WHERE id = $1", planID).Scan(&totalDays)
	if err != nil {
		return err
	}

	// Get completed days count
	var completedCount int
	err = s.db.QueryRow(ctx, "SELECT COUNT(*) FROM user_reading_plan_progress WHERE user_reading_plan_id = $1", userPlanID).Scan(&completedCount)
	if err != nil {
		return err
	}

	// If completed all days, update status
	if completedCount >= totalDays {
		_, err = s.db.Exec(ctx, `
			UPDATE user_reading_plans
			SET status = 'completed', updated_at = NOW()
			WHERE id = $1 AND status = 'active'
		`, userPlanID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *readingPlanService) UnmarkDayComplete(ctx context.Context, userID, planID uuid.UUID, dayNumber int) error {
	// First find the user's active plan
	var userPlanID uuid.UUID
	err := s.db.QueryRow(ctx, `
		SELECT id FROM user_reading_plans
		WHERE user_id = $1 AND reading_plan_id = $2 AND status = 'active'
	`, userID, planID).Scan(&userPlanID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.ErrNotFound
		}
		return err
	}

	// Delete progress
	query := `
		DELETE FROM user_reading_plan_progress
		WHERE user_reading_plan_id = $1 AND day_number = $2
	`
	_, err = s.db.Exec(ctx, query, userPlanID, dayNumber)
	return err
}

func (s *readingPlanService) GetPlanProgress(ctx context.Context, userID, planID uuid.UUID) ([]int, error) {
	var userPlanID uuid.UUID
	err := s.db.QueryRow(ctx, `
		SELECT id FROM user_reading_plans
		WHERE user_id = $1 AND reading_plan_id = $2
	`, userID, planID).Scan(&userPlanID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	query := `
		SELECT day_number
		FROM user_reading_plan_progress
		WHERE user_reading_plan_id = $1
		ORDER BY day_number ASC
	`
	rows, err := s.db.Query(ctx, query, userPlanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var completedDays []int
	for rows.Next() {
		var d int
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		completedDays = append(completedDays, d)
	}
	return completedDays, nil
}

func (s *readingPlanService) Unsubscribe(ctx context.Context, userID, planID uuid.UUID) error {
	// Delete the subscription. Cascade delete will handle progress.
	query := `
		DELETE FROM user_reading_plans
		WHERE user_id = $1 AND reading_plan_id = $2 AND status = 'active'
	`
	cmd, err := s.db.Exec(ctx, query, userID, planID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
