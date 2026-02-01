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

func TestGroupAPI_Contract(t *testing.T) {
	r, pool, cleanup := SetupContractTest(t)
	defer cleanup()

	if testing.Short() {
		t.Skip("skipping contract test in short mode")
	}

	// Seed Users
	ctx := context.Background()
	user1ID := "00000000-0000-0000-0000-000000000001"
	user2ID := "00000000-0000-0000-0000-000000000002"

	// Upsert users
	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, firebase_uid, email, username, created_at, updated_at)
		VALUES ($1, 'test-uid-1', 'user1@example.com', 'user1', NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, user1ID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, firebase_uid, email, username, created_at, updated_at)
		VALUES ($1, 'test-uid-2', 'user2@example.com', 'user2', NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, user2ID)
	require.NoError(t, err)

	// Create connections table if needed? No, SetupIntegrationDB handles migrations.
	// But we might need to seed a connection for Direct Group test.

	t.Run("CreateGroup", func(t *testing.T) {
		payload := map[string]string{
			"name":        "Test Group",
			"description": "A test group description",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/groups", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-User-ID", user1ID)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp["id"])

		groupID := resp["id"]

		// Verify User 1 is admin
		// We can check via GetGroupMembers or DB

		// Subtest: GetMembers
		t.Run("GetMembers", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/groups/"+groupID+"/members", nil)
			req.Header.Set("X-Test-User-ID", user1ID)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var members []map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &members)
			require.NoError(t, err)
			assert.Len(t, members, 1)
			assert.Equal(t, user1ID, members[0]["user_id"])
			assert.Equal(t, "admin", members[0]["role"])
		})

		// Subtest: JoinGroup (User 2)
		t.Run("JoinGroup", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/groups/"+groupID+"/join", nil)
			req.Header.Set("X-Test-User-ID", user2ID)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			// Verify members count increased
			reqMembers := httptest.NewRequest("GET", "/api/groups/"+groupID+"/members", nil)
			reqMembers.Header.Set("X-Test-User-ID", user1ID) // User 1 checks
			wMembers := httptest.NewRecorder()
			r.ServeHTTP(wMembers, reqMembers)

			var members []map[string]interface{}
			err := json.Unmarshal(wMembers.Body.Bytes(), &members)
			require.NoError(t, err)
			assert.Len(t, members, 2)
		})

		// Subtest: LeaveGroup (User 2)
		t.Run("LeaveGroup", func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/groups/"+groupID+"/leave", nil)
			req.Header.Set("X-Test-User-ID", user2ID)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			// Verify members count decreased
			reqMembers := httptest.NewRequest("GET", "/api/groups/"+groupID+"/members", nil)
			reqMembers.Header.Set("X-Test-User-ID", user1ID)
			wMembers := httptest.NewRecorder()
			r.ServeHTTP(wMembers, reqMembers)

			var members []map[string]interface{}
			err := json.Unmarshal(wMembers.Body.Bytes(), &members)
			require.NoError(t, err)
			assert.Len(t, members, 1)
		})
	})

	t.Run("GetOrCreateDirectGroup", func(t *testing.T) {
		// Prerequisite: Users must be connected
		// We can inject connection directly into DB for this test setup
		_, err := pool.Exec(ctx, `
			INSERT INTO connections (requester_id, receiver_id, status, created_at, updated_at)
			VALUES ($1, $2, 'accepted', NOW(), NOW())
			ON CONFLICT DO NOTHING
		`, user1ID, user2ID)
		require.NoError(t, err)

		payload := map[string]string{
			"partner_id": user2ID,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/groups/direct", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-User-ID", user1ID)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code) // Created first time

		var resp map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp["id"])
		groupID := resp["id"]

		// Call again, should return same ID and 200 OK
		req2 := httptest.NewRequest("POST", "/api/groups/direct", bytes.NewReader(body))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("X-Test-User-ID", user1ID)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)
		var resp2 map[string]string
		err = json.Unmarshal(w2.Body.Bytes(), &resp2)
		require.NoError(t, err)
		assert.Equal(t, groupID, resp2["id"])
	})
}
