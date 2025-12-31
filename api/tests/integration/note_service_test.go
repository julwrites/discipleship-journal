package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"discipleship_journal_api/models"
	"discipleship_journal_api/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoteService_Integration(t *testing.T) {
	pool, cleanup := SetupIntegrationDB(t)
	defer cleanup()

	service := services.NewNoteService(pool)
	ctx := context.Background()

	// Seed a user
	userID := "user-integration-test"
	// In a real scenario we might need to insert the user into the 'users' table if referential integrity is enforced.
	// Checking the schema...
	// Usually `users` table exists. Let's try to insert a user first to be safe.
	_, err := pool.Exec(ctx, "INSERT INTO users (id, firebase_uid, email, created_at, updated_at) VALUES (gen_random_uuid(), $1, 'test@example.com', NOW(), NOW()) ON CONFLICT (firebase_uid) DO NOTHING", userID)
	require.NoError(t, err)

	// We need the internal UUID for the user if the service uses it.
	// Looking at NoteService, it uses userID string directly in the query.
	// But wait, the notes table usually links to users.id (UUID).
	// Let's check the NoteService implementation again.
	// The `CreateNote` query uses `user_id` column.
	// If `user_id` in `notes` is a UUID, then passing "user-integration-test" (string) will fail if it's not a UUID.

	// Let's check the schema or assume based on `NoteService` using `string`.
	// In `api/services/note_service.go`, `userID` is string.
	// But typically in this project, `firebase_uid` is used for lookup, but internal FKs are UUIDs.
	// However, `NoteService` seems to take `userID` as is and insert it into `notes.user_id`.
	// If `notes.user_id` is UUID, we must provide a UUID.

	// Let's fetch the internal UUID for the seeded user.
	var internalUserID string
	err = pool.QueryRow(ctx, "SELECT id FROM users WHERE firebase_uid=$1", userID).Scan(&internalUserID)
	require.NoError(t, err)

	t.Run("CRUD Lifecycle", func(t *testing.T) {
		title := "Integration Test Note"
		content := json.RawMessage(`{"text": "integration content"}`)

		// 1. Create
		note, err := service.CreateNote(ctx, internalUserID, title, content)
		require.NoError(t, err)
		assert.NotEmpty(t, note.ID)
		assert.Equal(t, internalUserID, note.UserID)
		assert.Equal(t, title, note.Title)
		assert.JSONEq(t, string(content), string(note.Content))

		// 2. Get
		fetched, err := service.GetNote(ctx, internalUserID, note.ID)
		require.NoError(t, err)
		assert.Equal(t, note.ID, fetched.ID)
		assert.Equal(t, title, fetched.Title)

		// 3. Update
		newTitle := "Updated Title"
		newContent := json.RawMessage(`{"text": "updated content"}`)
		err = service.UpdateNote(ctx, internalUserID, note.ID, newTitle, newContent)
		require.NoError(t, err)

		updated, err := service.GetNote(ctx, internalUserID, note.ID)
		require.NoError(t, err)
		assert.Equal(t, newTitle, updated.Title)
		assert.JSONEq(t, string(newContent), string(updated.Content))

		// 4. Delete
		err = service.DeleteNote(ctx, internalUserID, note.ID)
		require.NoError(t, err)

		// 5. Verify Delete (Soft Delete)
		deleted, err := service.GetNote(ctx, internalUserID, note.ID)
		assert.Error(t, err)
		assert.Equal(t, models.ErrNotFound, err)
		assert.Nil(t, deleted)

		// Verify it's still in DB but with deleted_at
		var deletedAt *time.Time
		err = pool.QueryRow(ctx, "SELECT deleted_at FROM notes WHERE id=$1", note.ID).Scan(&deletedAt)
		require.NoError(t, err)
		assert.NotNil(t, deletedAt)
	})

	t.Run("Advanced Filtering", func(t *testing.T) {
		// Clean up existing notes for this user to have a clean slate (optional, but good for reliable counting)
		// We can just create new unique notes.

		baseTitle := "Filter Note"

		// Create 3 notes with different timestamps/titles
		// Note 1: "Alpha"
		n1, err := service.CreateNote(ctx, internalUserID, baseTitle+" Alpha", json.RawMessage(`{}`))
		require.NoError(t, err)

		// Note 2: "Beta"
		n2, err := service.CreateNote(ctx, internalUserID, baseTitle+" Beta", json.RawMessage(`{}`))
		require.NoError(t, err)

		// Note 3: "Gamma"
		n3, err := service.CreateNote(ctx, internalUserID, baseTitle+" Gamma", json.RawMessage(`{}`))
		require.NoError(t, err)

		// Manually update timestamps to test date filtering and sorting
		// N1: 2 days ago
		// N2: 1 day ago
		// N3: Just now (default)

		t1 := time.Now().Add(-48 * time.Hour)
		t2 := time.Now().Add(-24 * time.Hour)

		_, err = pool.Exec(ctx, "UPDATE notes SET updated_at=$1 WHERE id=$2", t1, n1.ID)
		require.NoError(t, err)
		_, err = pool.Exec(ctx, "UPDATE notes SET updated_at=$1 WHERE id=$2", t2, n2.ID)
		require.NoError(t, err)

		// Test Search
		filter := services.NoteFilter{SearchQuery: "Beta"}
		notes, total, err := service.GetNotes(ctx, internalUserID, 1, 10, filter)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, notes, 1)
		assert.Equal(t, n2.ID, notes[0].ID)

		// Test Sorting (Updated At ASC) -> Oldest first -> N1, N2, N3
		filter = services.NoteFilter{SortBy: "updated_at", SortOrder: "asc"}
		notes, total, err = service.GetNotes(ctx, internalUserID, 1, 10, filter)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, 3)
		// We need to filter the results to only include the ones we created,
		// because other tests might run in parallel or sequence on the same DB (if not cleaned).
		// But here we are using a fresh container per test suite typically, but within the suite we share it.
		// Let's check the IDs.

		var foundN1, foundN2, foundN3 bool
		for _, n := range notes {
			if n.ID == n1.ID { foundN1 = true }
			if n.ID == n2.ID { foundN2 = true }
			if n.ID == n3.ID { foundN3 = true }
		}

		assert.True(t, foundN1, "N1 should be found")
		assert.True(t, foundN2, "N2 should be found")
		assert.True(t, foundN3, "N3 should be found")

		// With pagination 10, we should see them if total is small.
		if len(notes) >= 3 {
             // Check relative order of our notes
             idx1, idx2, idx3 := -1, -1, -1
             for i, n := range notes {
                 if n.ID == n1.ID { idx1 = i }
                 if n.ID == n2.ID { idx2 = i }
                 if n.ID == n3.ID { idx3 = i }
             }
             if idx1 != -1 && idx2 != -1 && idx3 != -1 {
                 assert.True(t, idx1 < idx2, "N1 should be before N2")
                 assert.True(t, idx2 < idx3, "N2 should be before N3")
             }
		}

		// Test Date Filtering
		// Filter for notes updated in the last 30 hours (should include N2 and N3, exclude N1)
		startFilter := time.Now().Add(-30 * time.Hour)
		filter = services.NoteFilter{StartDate: &startFilter}
		notes, _, err = service.GetNotes(ctx, internalUserID, 1, 10, filter)
		require.NoError(t, err)

		foundN1 = false
		foundN2 = false
		foundN3 = false
		for _, n := range notes {
			if n.ID == n1.ID { foundN1 = true }
			if n.ID == n2.ID { foundN2 = true }
			if n.ID == n3.ID { foundN3 = true }
		}
		assert.False(t, foundN1, "N1 (48h old) should be excluded")
		assert.True(t, foundN2, "N2 (24h old) should be included")
		assert.True(t, foundN3, "N3 (new) should be included")
	})
}
