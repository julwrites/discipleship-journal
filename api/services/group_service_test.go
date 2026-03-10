package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"discipleship_journal_api/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

func TestCreateGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID
	name := "Test Group"
	desc := "A test group"
	groupType := "group"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `groups`").
			WithArgs(sqlmock.AnyArg(), name, &desc, userID, groupType).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO group_members").
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		group, err := service.CreateGroup(ctx, userID, name, &desc, groupType)

		assert.NoError(t, err)
		assert.NotNil(t, group)
		assert.NotEmpty(t, group.ID)
		assert.Equal(t, name, group.Name)
		assert.Equal(t, &desc, group.Description)
		assert.Equal(t, userID, group.CreatedBy)
		assert.Equal(t, "admin", group.Role)
	})

	t.Run("db error on insert group", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `groups`").
			WithArgs(sqlmock.AnyArg(), name, &desc, userID, groupType).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		group, err := service.CreateGroup(ctx, userID, name, &desc, groupType)

		assert.Error(t, err)
		assert.Nil(t, group)
	})

	t.Run("db error on insert member", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `groups`").
			WithArgs(sqlmock.AnyArg(), name, &desc, userID, groupType).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO group_members").
			WithArgs(sqlmock.AnyArg(), userID).
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

func TestCreateGroup_DefaultType(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID
	name := "Default Type Group"
	desc := "A test group"
	// Empty group type should default to "group"

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `groups`").
		WithArgs(sqlmock.AnyArg(), name, &desc, userID, "group"). // Expect "group" default
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO group_members").
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	group, err := service.CreateGroup(ctx, userID, name, &desc, "")

	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "group", group.Type)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUserGroups(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	group1ID := uuid.New().String()
	group2ID := uuid.New().String()
	user2ID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT g.id, g.name, g.description, g.created_by, g.type, gm.role").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_by", "type", "role"}).
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
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID
	user2ID := uuid.New().String()
	query := "test"

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT g.id, g.name, g.description, g.created_by, g.type").
			WithArgs(userID, "%"+query+"%").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_by", "type", "role"}).
				AddRow(groupID, "Test Group", strPtr("Desc"), user2ID, strPtr("group"), ""))

		groups, err := service.SearchGroups(ctx, query, userID)

		assert.NoError(t, err)
		if assert.Len(t, groups, 1) {
			assert.Equal(t, groupID, groups[0].ID)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT g.id, g.name, g.description, g.created_by, g.type").
			WithArgs(userID, "%"+query+"%").
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
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("INSERT INTO group_members").
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := service.JoinGroup(ctx, groupID, userID)
		assert.NoError(t, err)
	})

	t.Run("already member", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		err := service.JoinGroup(ctx, groupID, userID)
		assert.Error(t, err)
		assert.Equal(t, models.ErrAlreadyExists, err)
	})
}

func TestLeaveGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM group_members").
			WithArgs(groupID, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := service.LeaveGroup(ctx, groupID, userID)
		assert.NoError(t, err)
	})

	t.Run("not member", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM group_members").
			WithArgs(groupID, userID).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := service.LeaveGroup(ctx, groupID, userID)
		assert.Error(t, err)
		assert.Equal(t, models.ErrNotFound, err)
	})
}

func TestGetGroupMembers(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(groupID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectQuery("SELECT gm.user_id, COALESCE").
			WithArgs(groupID).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "display_name", "email", "role", "joined_at"}).
				AddRow(userID, "User 1", "user1@example.com", "member", now))

		members, err := service.GetGroupMembers(ctx, groupID, userID)
		assert.NoError(t, err)
		assert.Len(t, members, 1)
		assert.Equal(t, userID, members[0].UserID)
	})
}

func TestAddGroupMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	adminID := uuid.New().String()
	targetID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID

	t.Run("success", func(t *testing.T) {
		// Admin Check
		mock.ExpectQuery("SELECT role FROM group_members").
			WithArgs(groupID, adminID).
			WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("admin"))

		// Connection check
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(adminID, targetID, targetID, adminID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Insert
		mock.ExpectExec("INSERT INTO group_members").
			WithArgs(groupID, targetID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Fetch details for notification
		mock.ExpectQuery("SELECT name FROM `groups`").
			WithArgs(groupID).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Test Group"))

		err := service.AddGroupMember(ctx, adminID, groupID, targetID)
		assert.NoError(t, err)
	})

	t.Run("not admin", func(t *testing.T) {
		mock.ExpectQuery("SELECT role FROM group_members").
			WithArgs(groupID, adminID).
			WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("member"))

		err := service.AddGroupMember(ctx, adminID, groupID, targetID)
		assert.Error(t, err)
		assert.EqualError(t, err, "admin rights required")
	})

	t.Run("not connected", func(t *testing.T) {
		// Admin Check
		mock.ExpectQuery("SELECT role FROM group_members").
			WithArgs(groupID, adminID).
			WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("admin"))

		// Connection check
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(adminID, targetID, targetID, adminID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		err := service.AddGroupMember(ctx, adminID, groupID, targetID)
		assert.Error(t, err)
		assert.EqualError(t, err, "user is not in your connections")
	})
}

func TestAddGroupMember_NameQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	adminID := uuid.New().String()
	targetID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID

	// Admin Check
	mock.ExpectQuery("SELECT role FROM group_members").
		WithArgs(groupID, adminID).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("admin"))

	// Connection check
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(adminID, targetID, targetID, adminID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Insert
	mock.ExpectExec("INSERT INTO group_members").
		WithArgs(groupID, targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Fetch name error (should not fail function)
	mock.ExpectQuery("SELECT name FROM `groups`").
		WithArgs(groupID).
		WillReturnError(errors.New("db error"))

	err = service.AddGroupMember(ctx, adminID, groupID, targetID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveGroupMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	adminID := uuid.New().String()
	targetID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID

	t.Run("success", func(t *testing.T) {
		// Admin Check
		mock.ExpectQuery("SELECT role FROM group_members").
			WithArgs(groupID, adminID).
			WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("admin"))

		// Delete
		mock.ExpectExec("DELETE FROM group_members").
			WithArgs(groupID, targetID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := service.RemoveGroupMember(ctx, adminID, groupID, targetID)
		assert.NoError(t, err)
	})
}

func TestGetOrCreateDirectGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	partnerID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID

	t.Run("success new group", func(t *testing.T) {
		// Check connection
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(userID, partnerID, partnerID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Check existing group
		mock.ExpectQuery("SELECT g.id FROM `groups` g").
			WithArgs(userID, partnerID).
			WillReturnError(sql.ErrNoRows)

		// Create
		mock.ExpectQuery("SELECT COALESCE").WithArgs(partnerID).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User 2"))
		mock.ExpectQuery("SELECT COALESCE").WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User 1"))

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `groups`").
			WithArgs(sqlmock.AnyArg(), "Direct: User 1 & User 2", userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO group_members").WithArgs(sqlmock.AnyArg(), userID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO group_members").WithArgs(sqlmock.AnyArg(), partnerID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		group, created, err := service.GetOrCreateDirectGroup(ctx, userID, partnerID)
		assert.NoError(t, err)
		assert.True(t, created)
		assert.NotEmpty(t, group.ID)
	})
}

func TestGetOrCreateDirectGroup_ConnectionCheckError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	partnerID := uuid.New().String()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(userID, partnerID, partnerID, userID).
		WillReturnError(errors.New("db error"))

	_, _, err = service.GetOrCreateDirectGroup(ctx, userID, partnerID)
	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
}

func TestGetOrCreateDirectGroup_InsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	partnerID := uuid.New().String()

	// Check connection
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(userID, partnerID, partnerID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Check existing group
	mock.ExpectQuery("SELECT g.id FROM `groups` g").
		WithArgs(userID, partnerID).
		WillReturnError(sql.ErrNoRows)

	// Fetch partner name
	mock.ExpectQuery("SELECT COALESCE").WithArgs(partnerID).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User 2"))

	// Fetch my name
	mock.ExpectQuery("SELECT COALESCE").WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User 1"))

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `groups`").
		WithArgs(sqlmock.AnyArg(), "Direct: User 1 & User 2", userID).
		WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()

	_, _, err = service.GetOrCreateDirectGroup(ctx, userID, partnerID)
	assert.Error(t, err)
}

func TestGetOrCreateDirectGroup_MemberInsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	partnerID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID

	// Check connection
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(userID, partnerID, partnerID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Check existing group
	mock.ExpectQuery("SELECT g.id FROM `groups` g").
		WithArgs(userID, partnerID).
		WillReturnError(sql.ErrNoRows)

	// Fetch partner name
	mock.ExpectQuery("SELECT COALESCE").WithArgs(partnerID).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User 2"))

	// Fetch my name
	mock.ExpectQuery("SELECT COALESCE").WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User 1"))

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `groups`").
		WithArgs(sqlmock.AnyArg(), "Direct: User 1 & User 2", userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO group_members").WithArgs(sqlmock.AnyArg(), userID).WillReturnError(errors.New("member insert error"))
	mock.ExpectRollback()

	_, _, err = service.GetOrCreateDirectGroup(ctx, userID, partnerID)
	assert.Error(t, err)
}

func TestGroupService_GetOrCreateDirectGroup_SelfError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()

	group, created, err := service.GetOrCreateDirectGroup(ctx, userID, userID)
	assert.Error(t, err)
	assert.Nil(t, group)
	assert.False(t, created)
	assert.EqualError(t, err, "cannot create direct group with yourself")
}

func TestGroupService_GetOrCreateDirectGroup_NotConnected(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	partnerID := uuid.New().String()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(userID, partnerID, partnerID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	group, created, err := service.GetOrCreateDirectGroup(ctx, userID, partnerID)
	assert.Error(t, err)
	assert.Nil(t, group)
	assert.False(t, created)
	assert.EqualError(t, err, "user is not in your connections")
}

func TestGroupService_GetGroupMembers_NotMember(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(groupID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	members, err := service.GetGroupMembers(ctx, groupID, userID)
	assert.Error(t, err)
	assert.Nil(t, members)
	assert.EqualError(t, err, "access denied")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGroupService_AddGroupMember_AdminCheckError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	adminID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID
	targetID := uuid.New().String()

	mock.ExpectQuery("SELECT role FROM group_members").
		WithArgs(groupID, adminID).
		WillReturnError(errors.New("db error"))

	err = service.AddGroupMember(ctx, adminID, groupID, targetID)
	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
}

func TestGroupService_RemoveGroupMember_AdminCheckError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	adminID := uuid.New().String()
	groupID := uuid.New().String()
	_ = groupID
	targetID := uuid.New().String()

	mock.ExpectQuery("SELECT role FROM group_members").
		WithArgs(groupID, adminID).
		WillReturnError(sql.ErrNoRows) // User not in group

	err = service.RemoveGroupMember(ctx, adminID, groupID, targetID)
	assert.Error(t, err)
	assert.EqualError(t, err, "access denied")
}

func TestGroupService_GetOrCreateDirectGroup_ExistingGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	partnerID := uuid.New().String()
	existingGroupID := uuid.New().String()

	// Check connection
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(userID, partnerID, partnerID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Check existing group - FOUND
	mock.ExpectQuery("SELECT g.id FROM `groups` g").
		WithArgs(userID, partnerID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(existingGroupID))

	group, created, err := service.GetOrCreateDirectGroup(ctx, userID, partnerID)
	assert.NoError(t, err)
	assert.False(t, created)
	assert.Equal(t, existingGroupID, group.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGroupService_GetOrCreateDirectGroup_PartnerNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	partnerID := uuid.New().String()

	// Check connection
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(userID, partnerID, partnerID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Check existing group - Not Found
	mock.ExpectQuery("SELECT g.id FROM `groups` g").
		WithArgs(userID, partnerID).
		WillReturnError(sql.ErrNoRows)

	// Fetch partner name - Not Found
	mock.ExpectQuery("SELECT COALESCE").WithArgs(partnerID).
		WillReturnError(sql.ErrNoRows)

	group, created, err := service.GetOrCreateDirectGroup(ctx, userID, partnerID)
	assert.Error(t, err)
	assert.Equal(t, models.ErrNotFound, err)
	assert.Nil(t, group)
	assert.False(t, created)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGroupService_GetOrCreateDirectGroup_TransactionError(t *testing.T) {
	db, mock, err := sqlmock.New()
	_ = mock
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockNotif := NewMockNotificationService()
	service := NewGroupService(db, mockNotif)

	ctx := context.Background()
	userID := uuid.New().String()
	partnerID := uuid.New().String()

	// Check connection
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(userID, partnerID, partnerID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Check existing group - Not Found
	mock.ExpectQuery("SELECT g.id FROM `groups` g").
		WithArgs(userID, partnerID).
		WillReturnError(sql.ErrNoRows)

	// Fetch partner name
	mock.ExpectQuery("SELECT COALESCE").WithArgs(partnerID).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User 2"))

	// Fetch my name
	mock.ExpectQuery("SELECT COALESCE").WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("User 1"))

	// Begin Transaction - Error
	mock.ExpectBegin().WillReturnError(errors.New("tx error"))

	group, created, err := service.GetOrCreateDirectGroup(ctx, userID, partnerID)
	assert.Error(t, err)
	assert.EqualError(t, err, "tx error")
	assert.Nil(t, group)
	assert.False(t, created)
	assert.NoError(t, mock.ExpectationsWereMet())
}
