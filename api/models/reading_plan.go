package models

import (
	"time"

	"github.com/google/uuid"
)

type ReadingPlan struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Days        int       `json:"days"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ReadingPlanDay struct {
	ID            uuid.UUID `json:"id"`
	ReadingPlanID uuid.UUID `json:"reading_plan_id"`
	DayNumber     int       `json:"day_number"`
	Passage       string    `json:"passage"`
	CreatedAt     time.Time `json:"created_at"`
}

type UserReadingPlan struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	ReadingPlanID uuid.UUID `json:"reading_plan_id"`
	StartDate     time.Time `json:"start_date"`
	Status        string    `json:"status"` // active, completed
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Plan          *ReadingPlan `json:"plan,omitempty"`
}

type UserReadingPlanProgress struct {
	ID                uuid.UUID `json:"id"`
	UserReadingPlanID uuid.UUID `json:"user_reading_plan_id"`
	DayNumber         int       `json:"day_number"`
	CompletedAt       time.Time `json:"completed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
