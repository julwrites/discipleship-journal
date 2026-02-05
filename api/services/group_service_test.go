package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"discipleship_journal_api/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

func TestCreateGroup(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(mock, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	name := "Test Group"
	desc := "A test group"
	groupType := "group"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO groups").
			WithArgs(name, &desc, userID, groupType).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(groupID))

		mock.ExpectExec("INSERT INTO group_members").
			WithArgs(groupID, userID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mock.ExpectCommit()

		group, err := service.CreateGroup(ctx, userID, name, &desc, groupType)

		assert.NoError(t, err)
		assert.NotNil(t, group)
		assert.Equal(t, groupID, group.ID)
		assert.Equal(t, name, group.Name)
		assert.Equal(t, &desc, group.Description)
		assert.Equal(t, userID, group.CreatedBy)
		assert.Equal(t, "admin", group.Role)
	})

	t.Run("db error on insert group", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO groups").
			WithArgs(name, &desc, userID, groupType).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		group, err := service.CreateGroup(ctx, userID, name, &desc, groupType)

		assert.Error(t, err)
		assert.Nil(t, group)
	})

	t.Run("db error on insert member", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO groups").
			WithArgs(name, &desc, userID, groupType).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(groupID))

		mock.ExpectExec("INSERT INTO group_members").
			WithArgs(groupID, userID).
			WillReturnError(errors.New("db error member"))
		mock.ExpectRollback()

		group, err := service.CreateGroup(ctx, userID, name, &desc, groupType)

		assert.Error(t, err)
		assert.Nil(t, group)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListUserGroups(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(mock, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	group1ID := uuid.New().String()
	group2ID := uuid.New().String()
	user2ID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT g.id, g.name, g.description, g.created_by, g.type, gm.role").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_by", "type", "role"}).
				AddRow(group1ID, "Group 1", strPtr("Desc 1"), user2ID, strPtr("group"), "member").
				AddRow(group2ID, "Group 2", nil, userID, strPtr("group"), "admin"))

		groups, err := service.ListUserGroups(ctx, userID)

		assert.NoError(t, err)
		if assert.Len(t, groups, 2) {
			assert.Equal(t, group1ID, groups[0].ID)
			assert.Equal(t, "member", groups[0].Role)
			assert.Equal(t, group2ID, groups[1].ID)
			assert.Equal(t, "admin", groups[1].Role)
			assert.Nil(t, groups[1].Description)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT g.id, g.name, g.description, g.created_by, g.type, gm.role").
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		groups, err := service.ListUserGroups(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, groups)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSearchGroups(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(mock, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	user2ID := uuid.New().String()
	query := "test"

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT g.id, g.name, g.description, g.created_by, g.type").
			WithArgs("%"+query+"%", userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "created_by", "type", "role"}).
				AddRow(groupID, "Test Group", strPtr("Desc"), user2ID, strPtr("group"), ""))

		groups, err := service.SearchGroups(ctx, query, userID)

		assert.NoError(t, err)
		if assert.Len(t, groups, 1) {
			assert.Equal(t, groupID, groups[0].ID)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT g.id, g.name, g.description, g.created_by, g.type").
			WithArgs("%"+query+"%", userID).
			WillReturnError(errors.New("db error"))

		groups, err := service.SearchGroups(ctx, query, userID)

		assert.Error(t, err)
		assert.Nil(t, groups)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestJoinGroup(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(mock, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("INSERT INTO group_members").
			WithArgs(groupID, userID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := service.JoinGroup(ctx, groupID, userID)
		assert.NoError(t, err)
	})

	t.Run("already member", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		err := service.JoinGroup(ctx, groupID, userID)
		assert.Error(t, err)
		assert.Equal(t, models.ErrAlreadyExists, err)
	})
}

func TestLeaveGroup(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(mock, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM group_members").
			WithArgs(groupID, userID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := service.LeaveGroup(ctx, groupID, userID)
		assert.NoError(t, err)
	})

	t.Run("not member", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM group_members").
			WithArgs(groupID, userID).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := service.LeaveGroup(ctx, groupID, userID)
		assert.Error(t, err)
		assert.Equal(t, models.ErrNotFound, err)
	})
}

func TestGetGroupMembers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(mock, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectQuery("SELECT gm.user_id, COALESCE").
			WithArgs(groupID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "display_name", "email", "role", "joined_at"}).
				AddRow(userID, "User 1", "user1@example.com", "member", now))

		members, err := service.GetGroupMembers(ctx, groupID, userID)
		assert.NoError(t, err)
		assert.Len(t, members, 1)
		assert.Equal(t, userID, members[0].UserID)
	})

	t.Run("access denied", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		members, err := service.GetGroupMembers(ctx, groupID, userID)
		assert.Error(t, err)
		assert.Nil(t, members)
		assert.Equal(t, "access denied", err.Error())
	})
}

func TestAddGroupMember(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(mock, mockNotif)

	ctx := context.Background()
	adminID := uuid.New().String()
	targetID := uuid.New().String()
	groupID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT role FROM group_members").
			WithArgs(groupID, adminID).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("admin"))

		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(adminID, targetID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectExec("INSERT INTO group_members").
			WithArgs(groupID, targetID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mock.ExpectQuery("SELECT name FROM groups").
			WithArgs(groupID).
			WillReturnRows(pgxmock.NewRows([]string{"name"}).AddRow("Group 1"))

		done := make(chan bool)
		mockNotif.SendNotificationFunc = func(ctx context.Context, userID, title, body string, data map[string]string) error {
			assert.Equal(t, targetID, userID)
			close(done)
			return nil
		}

		err := service.AddGroupMember(ctx, adminID, groupID, targetID)
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("Notification not sent")
		}
	})

	t.Run("not admin", func(t *testing.T) {
		mock.ExpectQuery("SELECT role FROM group_members").
			WithArgs(groupID, adminID).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("member"))

		err := service.AddGroupMember(ctx, adminID, groupID, targetID)
		assert.Error(t, err)
		assert.Equal(t, "admin rights required", err.Error())
	})
}

func TestRemoveGroupMember(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(mock, mockNotif)

	ctx := context.Background()
	adminID := uuid.New().String()
	targetID := uuid.New().String()
	groupID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT role FROM group_members").
			WithArgs(groupID, adminID).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("admin"))

		mock.ExpectExec("DELETE FROM group_members").
			WithArgs(groupID, targetID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := service.RemoveGroupMember(ctx, adminID, groupID, targetID)
		assert.NoError(t, err)
	})

	t.Run("not admin", func(t *testing.T) {
		mock.ExpectQuery("SELECT role FROM group_members").
			WithArgs(groupID, adminID).
			WillReturnRows(pgxmock.NewRows([]string{"role"}).AddRow("member"))

		err := service.RemoveGroupMember(ctx, adminID, groupID, targetID)
		assert.Error(t, err)
		assert.Equal(t, "admin rights required", err.Error())
	})
}

func TestGetOrCreateDirectGroup(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(mock, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	partnerID := uuid.New().String()
	groupID := uuid.New().String()
	existingGroupID := uuid.New().String()

	t.Run("success new group", func(t *testing.T) {
		// 1. Check connection
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(userID, partnerID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		// 2. Check existing
		mock.ExpectQuery("SELECT g.id FROM groups g").
			WithArgs(userID, partnerID).
			WillReturnError(pgx.ErrNoRows)

		// 3. Create
		// Fetch names
		mock.ExpectQuery("SELECT COALESCE").
			WithArgs(partnerID).
			WillReturnRows(pgxmock.NewRows([]string{"name"}).AddRow("User 2"))
		mock.ExpectQuery("SELECT COALESCE").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"name"}).AddRow("User 1"))

		// Tx
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO groups").
			WithArgs("Direct: User 1 & User 2", userID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(groupID))

		mock.ExpectExec("INSERT INTO group_members").
			WithArgs(groupID, userID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
		mock.ExpectExec("INSERT INTO group_members").
			WithArgs(groupID, partnerID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mock.ExpectCommit()

		group, isNew, err := service.GetOrCreateDirectGroup(ctx, userID, partnerID)
		assert.NoError(t, err)
		assert.True(t, isNew)
		assert.Equal(t, groupID, group.ID)
	})

	t.Run("success existing group", func(t *testing.T) {
		// 1. Check connection
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(userID, partnerID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		// 2. Check existing
		mock.ExpectQuery("SELECT g.id FROM groups g").
			WithArgs(userID, partnerID).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(existingGroupID))

		group, isNew, err := service.GetOrCreateDirectGroup(ctx, userID, partnerID)
		assert.NoError(t, err)
		assert.False(t, isNew)
		assert.Equal(t, existingGroupID, group.ID)
	})
}
