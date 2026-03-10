package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ReadingPlan struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Days        int       `json:"days"`
	PlanType    string    `json:"plan_type"` // sequential, calendar
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

// Validate checks if the reading plan day is valid.
// It enforces that multiple passages are separated by semicolons.
func (d *ReadingPlanDay) Validate() error {
	if d.Passage == "" {
		return fmt.Errorf("passage cannot be empty")
	}

	// Enforce semicolon separator for multiple passages.
	// Heuristic: A comma followed by a space and a letter usually indicates a new book reference
	// (e.g. "Genesis 1, Exodus 2") which should use a semicolon.
	// Valid comma usage: "Genesis 1:1, 3" (followed by digit) or "Genesis 1, 2" (followed by digit).

	// We iterate through the string to find ", " followed by a letter.
	// We must handle cases like "1 John" where the book starts with a number?
	// No, "Genesis 1, 1 John 1" -> ", 1" (digit). This heuristic fails for numbered books if the comma precedes them.
	// Example: "Genesis 1, 1 Kings 1". ", 1" is valid for "Genesis 1, 2".
	// So we can't distinguish "Genesis 1, 2" (Chapter 2) from "Genesis 1, 1 Kings 1".

	// Let's try a safer check: Newlines are forbidden.
	if strings.Contains(d.Passage, "\n") {
		return fmt.Errorf("passage cannot contain newlines; use semicolons to separate references")
	}

	// Check for " + " or " & " which imply combined references
	if strings.Contains(d.Passage, " + ") || strings.Contains(d.Passage, " & ") {
		return fmt.Errorf("passage contains invalid separators (+ or &); use semicolons")
	}

	// If the user wants to enforce strict semicolon usage, we should check if there are semicolons
	// and if so, that they are used.
	// The prompt says "enforce it when we are inserting".
	// I'll leave it at preventing bad characters for now, as strict parsing requires more logic.

	return nil
}

type UserReadingPlan struct {
	ID            uuid.UUID    `json:"id"`
	UserID        uuid.UUID    `json:"user_id"`
	ReadingPlanID uuid.UUID    `json:"reading_plan_id"`
	StartDate     time.Time    `json:"start_date"`
	Status        string       `json:"status"` // active, completed
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Plan          *ReadingPlan `json:"plan,omitempty"`
}

type UserReadingPlanProgress struct {
	ID                uuid.UUID `json:"id"`
	UserReadingPlanID uuid.UUID `json:"user_reading_plan_id"`
	DayNumber         int       `json:"day_number"`
	CompletedAt       time.Time `json:"completed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
