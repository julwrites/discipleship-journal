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

// TestNoteAPI_Contract verifies that the backend response matches the frontend expectation.
// Reference: web/src/services/api.ts
func TestNoteAPI_Contract(t *testing.T) {
	r, pool, cleanup := SetupContractTest(t)
	defer cleanup()

	// Seed User
	ctx := context.Background()
	testUUID := "00000000-0000-0000-0000-000000000001"

	// We only need to seed if we are NOT skipping.
	if !testing.Short() {
		// Insert user to satisfy FK
		// We use ON CONFLICT DO NOTHING in case it runs multiple times or in parallel (though parallel integration tests usually get their own DB/Container)
		_, err := pool.Exec(ctx, "INSERT INTO users (id, firebase_uid, email, created_at, updated_at) VALUES ($1, 'test-uid', 'test@example.com', NOW(), NOW()) ON CONFLICT (id) DO NOTHING", testUUID)
		require.NoError(t, err)
	}

	t.Run("CreateNote Contract", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping contract test in short mode")
		}

		// Payload
		payload := map[string]interface{}{
			"title": "Contract Test Note",
			"content": map[string]interface{}{
				"type": "doc",
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Hello World",
							},
						},
					},
				},
			},
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.NotEmpty(t, resp["id"])
	})

	t.Run("GetNotes Contract", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping contract test in short mode")
		}

		req := httptest.NewRequest("GET", "/api/notes?limit=10&page=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		// Check 'data'
		_, ok := resp["data"].([]interface{})
		assert.True(t, ok, "'data' field should be an array")

		// Check 'meta'
		meta, ok := resp["meta"].(map[string]interface{})
		assert.True(t, ok, "'meta' field should be an object")
		assert.NotNil(t, meta["total"])
	})
}
