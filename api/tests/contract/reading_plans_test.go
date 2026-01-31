package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadingPlanAPI_Contract(t *testing.T) {
	r, pool, cleanup := SetupContractTest(t)
	defer cleanup()

	if testing.Short() {
		t.Skip("skipping contract test in short mode")
	}

	ctx := context.Background()
	planID := "00000000-0000-0000-0000-000000000100"
	userID := "00000000-0000-0000-0000-000000000001" // Default test user

	// Seed Plan
	_, err := pool.Exec(ctx, `
		INSERT INTO reading_plans (id, title, description, days, plan_type, created_at, updated_at)
		VALUES ($1, 'Test Plan', 'A test plan', 365, 'calendar', NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, planID)
	require.NoError(t, err)

	// Seed Plan Days (at least Day 1)
	_, err = pool.Exec(ctx, `
		INSERT INTO reading_plan_days (reading_plan_id, day_number, passage, created_at)
		VALUES ($1, 1, 'Genesis 1', NOW())
		ON CONFLICT DO NOTHING
	`, planID)
	require.NoError(t, err)

	t.Run("GetReadingPlans", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/reading-plans", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		plans, ok := resp["data"].([]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, plans)

		found := false
		for _, p := range plans {
			plan := p.(map[string]interface{})
			if plan["id"] == planID {
				found = true
				break
			}
		}
		assert.True(t, found, "Seeded plan should be returned")
	})

	t.Run("SubscribeToPlan", func(t *testing.T) {
		// Ensure clean state
		_, err := pool.Exec(ctx, "DELETE FROM user_reading_plans WHERE user_id=$1 AND reading_plan_id=$2", userID, planID)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/reading-plans/"+planID+"/subscribe", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify in MyPlans
		reqMy := httptest.NewRequest("GET", "/api/my-reading-plans", nil)
		wMy := httptest.NewRecorder()
		r.ServeHTTP(wMy, reqMy)

		var resp map[string]interface{}
		json.Unmarshal(wMy.Body.Bytes(), &resp)
		myPlans := resp["data"].([]interface{})

		found := false
		for _, p := range myPlans {
			plan := p.(map[string]interface{})
			if plan["reading_plan_id"] == planID {
				found = true
				assert.Equal(t, "active", plan["status"])
				break
			}
		}
		assert.True(t, found, "Subscribed plan should appear in user plans")
	})

	t.Run("MarkPlanDayComplete", func(t *testing.T) {
		// Prerequisite: Must be subscribed (previous test runs first, or we ensure it)
		// We can just rely on previous test or idempotent subscription.

		payload := map[string]interface{}{
			"day_number": 1,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/my-reading-plans/"+planID+"/progress", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify Progress
		reqProg := httptest.NewRequest("GET", "/api/my-reading-plans/"+planID+"/progress", nil)
		wProg := httptest.NewRecorder()
		r.ServeHTTP(wProg, reqProg)

		assert.Equal(t, http.StatusOK, wProg.Code)
		var resp map[string]interface{}
		json.Unmarshal(wProg.Body.Bytes(), &resp)

		days := resp["completed_days"].([]interface{})
		assert.Contains(t, days, float64(1)) // JSON numbers are float64
	})
}
