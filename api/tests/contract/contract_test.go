package contract

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"discipleship_journal_api/handlers"
	"discipleship_journal_api/services"
	"discipleship_journal_api/tests/integration"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock Auth Middleware to bypass Firebase
type MockAuthMiddleware struct{}

func (m *MockAuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Inject a fake user into context
		// We use the same key as defined in helpers.go (usually 'user_id' or similar)
		// But since we can't easily import internal keys if they are not exported,
		// we might rely on the fact that existing handlers use `GetUserUUID`.
		// Ideally, we'd reuse the middleware logic but with a mock verify.

		// However, `GetUserUUID` looks up the user by Firebase UID in the DB.
		// So we need to ensure the user exists in the DB first.
		next.ServeHTTP(w, r)
	})
}

// SetupRouter creates a router similar to main.go but with real DB and Mock Auth
func SetupRouter(t *testing.T) (*chi.Mux, func(), string) {
	// 1. Setup Real DB
	pool, cleanup := integration.SetupIntegrationDB(t)

	// 2. Set global DB (handlers often rely on it if not injected, but our main injects them)
	// However, some helpers might rely on global. To be safe, we can set it, but better to rely on injection.
	// Looking at main.go: `noteService := services.NewNoteService(database.DB)`
	// So we just need to pass `pool` to services.

	// 3. Create a Test User in the DB
	// We need a user to perform actions.
	testUserUID := "test-user-uid"
	var testUserID string

	// Manual insert since we have the pool
	row := pool.QueryRowContext(context.Background(),
		"INSERT INTO users (username, email, firebase_uid, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW()) RETURNING id",
		"testuser", "test@example.com", testUserUID)
	err := row.Scan(&testUserID)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// 4. Setup Services & Handlers
	// Mocks for external dependencies
	mockNotification := services.NewMockNotificationService()

	// Real Services using the integration DB pool
	noteService := services.NewNoteService(pool)
	groupService := services.NewGroupService(pool, mockNotification)

	noteHandler := handlers.NewNoteHandler(pool, noteService)
	groupHandler := handlers.NewGroupHandler(groupService)
	connectionHandler := handlers.NewConnectionHandler(pool, mockNotification)

	r := chi.NewRouter()

	// Inject Mock Auth Middleware that sets the TestUserKey
	// We need to match how `GetUserUUID` works.
	// It checks `r.Context().Value(handlers.TestUserKey)`.
	// So we just need to set that.
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), handlers.TestUserKey, testUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Post("/api/notes", noteHandler.CreateNote)
	r.Get("/api/notes", noteHandler.GetNotes)
	r.Put("/api/notes/{id}", noteHandler.UpdateNote)

	r.Post("/api/groups", groupHandler.CreateGroup)
	r.Post("/api/connections/request", connectionHandler.SendConnectionRequest)

	return r, cleanup, testUserID
}

func TestContract_CreateNote(t *testing.T) {
	r, cleanup, _ := SetupRouter(t)
	defer cleanup()

	// Payload matching web/src/services/api.ts: createNote
	// body: JSON.stringify({ title, content })
	// content is passed as `{ markdown: "..." }` in NoteEditor.tsx
	payload := `{"title": "Contract Test Note", "content": {"markdown": "# Hello World"}}`

	req := httptest.NewRequest("POST", "/api/notes", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Note: API currently returns 200 OK, not 201 Created for Notes.
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.NotEmpty(t, resp["id"], "Response should contain ID")
	// API only returns {id}
}

func TestContract_GetNotes_Filter(t *testing.T) {
	r, cleanup, _ := SetupRouter(t)
	defer cleanup()

	// 1. Create a note first
	payload := `{"title": "Searchable Note", "content": {"markdown": "Secret Content"}}`
	req := httptest.NewRequest("POST", "/api/notes", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// 2. Search for it using api.ts style params
	// params.append("q", filter.search);
	req = httptest.NewRequest("GET", "/api/notes?q=Searchable", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// API returns { "data": [...], "meta": {...} }
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	data, ok := resp["data"].([]interface{})
	require.True(t, ok, "data field should be an array")
	assert.NotEmpty(t, data)

	note := data[0].(map[string]interface{})
	assert.Equal(t, "Searchable Note", note["title"])
}

func TestContract_CreateGroup(t *testing.T) {
	r, cleanup, _ := SetupRouter(t)
	defer cleanup()

	// api.ts: createGroup(data: { name: string; description: string })
	payload := `{"name": "Contract Group", "description": "Testing Groups"}`

	req := httptest.NewRequest("POST", "/api/groups", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.NotEmpty(t, resp["id"])
	// API only returns {id} for groups too
}
