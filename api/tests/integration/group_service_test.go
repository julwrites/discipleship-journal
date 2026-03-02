package integration

import (
	"context"
	"testing"
	"time"

	"discipleship_journal_api/models"
	"discipleship_journal_api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupService_Integration(t *testing.T) {
	pool, cleanup := SetupIntegrationDB(t)
	defer cleanup()

	notificationService := services.NewMockNotificationService()
	service := services.NewGroupService(pool, notificationService)
	ctx := context.Background()

	// Helper to seed user
	createUser := func(firebaseUID string) string {
		id := uuid.New().String()
		// Ensure firebase_uid is unique
		// Explicitly cast parameter or use explicit string formatting to avoid type inference issues
		email := firebaseUID + "@test.com"
		_, err := pool.ExecContext(ctx, "INSERT INTO users (id, firebase_uid, email, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())", id, firebaseUID, email)
		require.NoError(t, err)
		return id
	}

	user1ID := createUser("user1_group_test")
	user2ID := createUser("user2_group_test")
	user3ID := createUser("user3_group_test")

	t.Run("CreateGroup", func(t *testing.T) {
		name := "Test Group"
		desc := "Description"
		group, err := service.CreateGroup(ctx, user1ID, name, &desc, "group")
		require.NoError(t, err)
		assert.NotEmpty(t, group.ID)
		assert.Equal(t, name, group.Name)
		assert.Equal(t, desc, *group.Description)
		assert.Equal(t, "admin", group.Role)
		assert.Equal(t, "group", group.Type)

		// Verify membership
		members, err := service.GetGroupMembers(ctx, group.ID, user1ID)
		require.NoError(t, err)
		assert.Len(t, members, 1)
		assert.Equal(t, user1ID, members[0].UserID)
		assert.Equal(t, "admin", members[0].Role)
	})

	t.Run("ListUserGroups", func(t *testing.T) {
		// Clean slate for user2
		user2Groups, err := service.ListUserGroups(ctx, user2ID)
		require.NoError(t, err)
		assert.Empty(t, user2Groups)

		// Create a group
		g1, err := service.CreateGroup(ctx, user2ID, "User 2 Group", nil, "group")
		require.NoError(t, err)

		groups, err := service.ListUserGroups(ctx, user2ID)
		require.NoError(t, err)
		assert.Len(t, groups, 1)
		assert.Equal(t, g1.ID, groups[0].ID)
	})

	t.Run("SearchGroups", func(t *testing.T) {
		// Create groups with specific names
		_, err := service.CreateGroup(ctx, user3ID, "Alpha Omega", nil, "group")
		require.NoError(t, err)
		_, err = service.CreateGroup(ctx, user3ID, "Beta Theta", nil, "group")
		require.NoError(t, err)

		// Create a direct group (should be excluded)
		// We need to simulate direct group creation via SQL or Service if possible
		// Service GetOrCreateDirectGroup requires connection, which might be tedious to set up just for this check
		// So we insert directly into DB to ensure it exists and has name containing "Alpha"
		_, err = pool.ExecContext(ctx, "INSERT INTO `groups` (name, type, created_by) VALUES (?, ?, ?)", "Direct Alpha", "direct", user3ID)
		require.NoError(t, err)

		// Search
		results, err := service.SearchGroups(ctx, "Alpha", user1ID)
		require.NoError(t, err)
		assert.NotEmpty(t, results)

		foundAlpha := false
		foundDirect := false
		for _, g := range results {
			if g.Name == "Alpha Omega" {
				foundAlpha = true
			}
			if g.Name == "Direct Alpha" {
				foundDirect = true
			}
		}
		assert.True(t, foundAlpha, "Should find 'Alpha Omega'")
		assert.False(t, foundDirect, "Should NOT find 'Direct Alpha'")
	})

	t.Run("Join and Leave Group", func(t *testing.T) {
		g, err := service.CreateGroup(ctx, user1ID, "Joinable Group", nil, "group")
		require.NoError(t, err)

		// User 2 joins
		err = service.JoinGroup(ctx, g.ID, user2ID)
		require.NoError(t, err)

		// Verify membership
		members, err := service.GetGroupMembers(ctx, g.ID, user1ID) // user1 is admin, can see members
		require.NoError(t, err)
		assert.Len(t, members, 2)

		// User 2 joins again (should fail)
		err = service.JoinGroup(ctx, g.ID, user2ID)
		assert.ErrorIs(t, err, models.ErrAlreadyExists)

		// User 2 leaves
		err = service.LeaveGroup(ctx, g.ID, user2ID)
		require.NoError(t, err)

		// Verify membership
		members, err = service.GetGroupMembers(ctx, g.ID, user1ID)
		require.NoError(t, err)
		assert.Len(t, members, 1)

		// User 2 leaves again (should fail)
		err = service.LeaveGroup(ctx, g.ID, user2ID)
		assert.ErrorIs(t, err, models.ErrNotFound)
	})

	t.Run("AddGroupMember (Admin Only & Connection Required & Notification)", func(t *testing.T) {
		g, err := service.CreateGroup(ctx, user1ID, "Admin Group", nil, "group")
		require.NoError(t, err)

		// 1. Try to add user2 (Not connected) -> Fail
		err = service.AddGroupMember(ctx, user1ID, g.ID, user2ID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user is not in your connections")

		// 2. Connect user1 and user2
		_, err = pool.ExecContext(ctx, "INSERT INTO connections (requester_id, receiver_id, status) VALUES (?, ?, 'accepted')", user1ID, user2ID)
		require.NoError(t, err)

		// Setup Notification Mock
		notifyChan := make(chan struct{}, 1)
		notificationService.SendNotificationFunc = func(ctx context.Context, userID, title, body string, data map[string]string) error {
			assert.Equal(t, user2ID, userID)
			assert.Equal(t, "Group Invitation", title)
			notifyChan <- struct{}{}
			return nil
		}

		// 3. Try to add user2 (Connected) -> Success
		err = service.AddGroupMember(ctx, user1ID, g.ID, user2ID)
		require.NoError(t, err)

		// Verify notification received
		select {
		case <-notifyChan:
			// Success
		case <-time.After(2 * time.Second):
			t.Fatal("Notification not sent within timeout")
		}

		// 4. Verify user2 is member
		members, err := service.GetGroupMembers(ctx, g.ID, user1ID)
		require.NoError(t, err)
		found := false
		for _, m := range members {
			if m.UserID == user2ID {
				found = true
				break
			}
		}
		assert.True(t, found)

		// 5. Try to add user3 by user2 (user2 is member, not admin) -> Fail
		err = service.AddGroupMember(ctx, user2ID, g.ID, user3ID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "admin rights required")
	})

	t.Run("GetOrCreateDirectGroup", func(t *testing.T) {
		// user2 and user3 are NOT connected
		_, _, err := service.GetOrCreateDirectGroup(ctx, user2ID, user3ID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user is not in your connections")

		// Connect user2 and user3
		_, err = pool.ExecContext(ctx, "INSERT INTO connections (requester_id, receiver_id, status) VALUES (?, ?, 'accepted')", user2ID, user3ID)
		require.NoError(t, err)

		// 1. Create first time
		group1, created1, err := service.GetOrCreateDirectGroup(ctx, user2ID, user3ID)
		require.NoError(t, err)
		assert.True(t, created1)
		assert.NotEmpty(t, group1.ID)

		// 2. Get second time (should return same group)
		group2, created2, err := service.GetOrCreateDirectGroup(ctx, user2ID, user3ID)
		require.NoError(t, err)
		assert.False(t, created2)
		assert.Equal(t, group1.ID, group2.ID)

		// 3. Verify members
		// Note: Direct groups implicitly make both admins in implementation
		members, err := service.GetGroupMembers(ctx, group1.ID, user2ID)
		require.NoError(t, err)
		assert.Len(t, members, 2)

		// 4. Try to create direct group with self
		_, _, err = service.GetOrCreateDirectGroup(ctx, user2ID, user2ID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot create direct group with yourself")
	})

	t.Run("RemoveGroupMember", func(t *testing.T) {
		g, err := service.CreateGroup(ctx, user1ID, "Remove Group", nil, "group")
		require.NoError(t, err)

		// Add user2 (assuming connection exists from previous test, but let's ensure it)
		// Since tests run in sequence and we are using same DB/users, connection between user1 and user2 exists.
		err = service.AddGroupMember(ctx, user1ID, g.ID, user2ID)
		require.NoError(t, err)

		// Remove user2
		err = service.RemoveGroupMember(ctx, user1ID, g.ID, user2ID)
		require.NoError(t, err)

		// Verify removed
		members, err := service.GetGroupMembers(ctx, g.ID, user1ID)
		require.NoError(t, err)
		assert.Len(t, members, 1) // only admin left
	})
}
