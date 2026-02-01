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

func TestConnectionAPI_Contract(t *testing.T) {
	r, pool, cleanup := SetupContractTest(t)
	defer cleanup()

	if testing.Short() {
		t.Skip("skipping contract test in short mode")
	}

	ctx := context.Background()
	userA_ID := "00000000-0000-0000-0000-000000000010"
	userB_ID := "00000000-0000-0000-0000-000000000011"

	// Seed Users
	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, firebase_uid, email, username, created_at, updated_at) VALUES
		($1, 'test-uid-A', 'userA@example.com', 'UserA', NOW(), NOW()),
		($2, 'test-uid-B', 'userB@example.com', 'UserB', NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, userA_ID, userB_ID)
	require.NoError(t, err)

	// Clean up connections for these users to ensure clean state
	_, err = pool.Exec(ctx, "DELETE FROM connections WHERE requester_id IN ($1, $2) OR receiver_id IN ($1, $2)", userA_ID, userB_ID)
	require.NoError(t, err)

	t.Run("ConnectionFlow", func(t *testing.T) {
		// 1. Send Request (User A -> User B)
		reqPayload := map[string]string{
			"receiver_id": userB_ID,
		}
		body, _ := json.Marshal(reqPayload)
		req := httptest.NewRequest("POST", "/api/connections/request", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-User-ID", userA_ID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp["id"])
		assert.Equal(t, "pending", resp["status"])

		connID := resp["id"].(string)

		// 2. List Connections (User B should see it)
		reqList := httptest.NewRequest("GET", "/api/connections", nil)
		reqList.Header.Set("X-Test-User-ID", userB_ID)
		wList := httptest.NewRecorder()
		r.ServeHTTP(wList, reqList)

		assert.Equal(t, http.StatusOK, wList.Code)
		var listResp []map[string]interface{}
		err = json.Unmarshal(wList.Body.Bytes(), &listResp)
		require.NoError(t, err)

		found := false
		for _, c := range listResp {
			if c["id"] == connID {
				found = true
				assert.Equal(t, "pending", c["status"])
				assert.Equal(t, userA_ID, c["requester_id"])
				break
			}
		}
		assert.True(t, found, "Connection request not found in receiver's list")

		// 3. Accept Request (User B)
		reqAccept := httptest.NewRequest("PUT", "/api/connections/"+connID, nil)
		reqAccept.Header.Set("X-Test-User-ID", userB_ID)
		wAccept := httptest.NewRecorder()
		r.ServeHTTP(wAccept, reqAccept)

		assert.Equal(t, http.StatusOK, wAccept.Code)

		// 4. Verify Status is Accepted (User A checks)
		reqListA := httptest.NewRequest("GET", "/api/connections", nil)
		reqListA.Header.Set("X-Test-User-ID", userA_ID)
		wListA := httptest.NewRecorder()
		r.ServeHTTP(wListA, reqListA)

		var listRespA []map[string]interface{}
		err = json.Unmarshal(wListA.Body.Bytes(), &listRespA)
		require.NoError(t, err)

		foundA := false
		for _, c := range listRespA {
			if c["id"] == connID {
				foundA = true
				assert.Equal(t, "accepted", c["status"])
				break
			}
		}
		assert.True(t, foundA, "Connection should be accepted")
	})
}
