package services

import (
	"context"
	"database/sql"
	"discipleship_journal_api/models"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedBy   string  `json:"created_by"`
	Type        string  `json:"type"`
	Role        string  `json:"role,omitempty"` // Derived field
}

type GroupMember struct {
	UserID      string    `json:"user_id"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	JoinedAt    time.Time `json:"joined_at"`
}

type GroupServiceInterface interface {
	CreateGroup(ctx context.Context, userID, name string, description *string, groupType string) (*Group, error)
	ListUserGroups(ctx context.Context, userID string) ([]Group, error)
	SearchGroups(ctx context.Context, query, userID string) ([]Group, error)
	JoinGroup(ctx context.Context, groupID, userID string) error
	LeaveGroup(ctx context.Context, groupID, userID string) error
	GetGroupMembers(ctx context.Context, groupID, userID string) ([]GroupMember, error)
	AddGroupMember(ctx context.Context, adminID, groupID, targetUserID string) error
	GetOrCreateDirectGroup(ctx context.Context, userID, partnerID string) (*Group, bool, error)
	RemoveGroupMember(ctx context.Context, adminID, groupID, targetUserID string) error
}

type groupService struct {
	db                  DBInterface
	notificationService NotificationService
}

func NewGroupService(db DBInterface, notificationService NotificationService) GroupServiceInterface {
	return &groupService{
		db:                  db,
		notificationService: notificationService,
	}
}

func (s *groupService) CreateGroup(ctx context.Context, userID, name string, description *string, groupType string) (*Group, error) {
	if groupType == "" {
		groupType = "group"
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	slog.Info("Creating group", "name", name, "description", description, "user", userID, "type", groupType)
	groupID := uuid.New().String()
	_, err = tx.ExecContext(ctx,
		"INSERT INTO `groups` (id, name, description, created_by, type) VALUES (?, ?, ?, ?, ?)",
		groupID, name, description, userID, groupType)
	if err != nil {
		slog.Error("Failed to create group", "error", err)
		return nil, err
	}

	// Add creator as admin
	_, err = tx.ExecContext(ctx,
		"INSERT INTO group_members (group_id, user_id, role) VALUES (?, ?, 'admin')",
		groupID, userID)
	if err != nil {
		slog.Error("Failed to add member", "error", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Group{
		ID:          groupID,
		Name:        name,
		Description: description,
		CreatedBy:   userID,
		Type:        groupType,
		Role:        "admin",
	}, nil
}

func (s *groupService) ListUserGroups(ctx context.Context, userID string) ([]Group, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT g.id, g.name, g.description, g.created_by, g.type, gm.role
		 FROM `+"`groups`"+` g
		 JOIN group_members gm ON g.id = gm.group_id
		 WHERE gm.user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		var description *string
		var groupType *string
		if err := rows.Scan(&g.ID, &g.Name, &description, &g.CreatedBy, &groupType, &g.Role); err != nil {
			slog.Warn("Failed to scan group row", "error", err)
			continue
		}
		g.Description = description
		if groupType != nil {
			g.Type = *groupType
		} else {
			g.Type = "group"
		}
		groups = append(groups, g)
	}

	if groups == nil {
		groups = []Group{}
	}
	return groups, nil
}

func (s *groupService) SearchGroups(ctx context.Context, query, userID string) ([]Group, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT g.id, g.name, g.description, g.created_by, g.type,
		 COALESCE((SELECT role FROM group_members WHERE group_id = g.id AND user_id = ?), '') as role
		 FROM `+"`groups`"+` g
		 WHERE g.name LIKE ? AND (g.type = 'group' OR g.type IS NULL) LIMIT 20`, userID, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		var description *string
		var groupType *string
		if err := rows.Scan(&g.ID, &g.Name, &description, &g.CreatedBy, &groupType, &g.Role); err != nil {
			slog.Warn("Failed to scan group row", "error", err)
			continue
		}
		g.Description = description
		if groupType != nil {
			g.Type = *groupType
		} else {
			g.Type = "group"
		}
		groups = append(groups, g)
	}

	if groups == nil {
		groups = []Group{}
	}
	return groups, nil
}

func (s *groupService) JoinGroup(ctx context.Context, groupID, userID string) error {
	var exists bool
	err := s.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return models.ErrAlreadyExists
	}

	_, err = s.db.ExecContext(ctx,
		"INSERT INTO group_members (group_id, user_id, role) VALUES (?, ?, 'member')",
		groupID, userID)
	return err
}

func (s *groupService) LeaveGroup(ctx context.Context, groupID, userID string) error {
	result, err := s.db.ExecContext(ctx,
		"DELETE FROM group_members WHERE group_id = ? AND user_id = ?",
		groupID, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *groupService) GetGroupMembers(ctx context.Context, groupID, userID string) ([]GroupMember, error) {
	var isMember bool
	err := s.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupID, userID).Scan(&isMember)
	if err != nil {
		return nil, err
	}

	if !isMember {
		// Use a specific error for forbidden access?
		// For now returning generic error or we can define ErrForbidden
		return nil, fmt.Errorf("access denied")
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT gm.user_id, COALESCE(u.username, u.email), u.email, gm.role, gm.joined_at
		 FROM group_members gm
		 JOIN users u ON gm.user_id = u.id
		 WHERE gm.group_id = ?`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []GroupMember
	for rows.Next() {
		var m GroupMember
		var joinedAt time.Time
		if err := rows.Scan(&m.UserID, &m.DisplayName, &m.Email, &m.Role, &joinedAt); err != nil {
			continue
		}
		m.JoinedAt = joinedAt
		members = append(members, m)
	}

	if members == nil {
		members = []GroupMember{}
	}
	return members, nil
}

func (s *groupService) AddGroupMember(ctx context.Context, adminID, groupID, targetUserID string) error {
	// Verify admin role
	var role string
	err := s.db.QueryRowContext(ctx,
		"SELECT role FROM group_members WHERE group_id = ? AND user_id = ?",
		groupID, adminID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("access denied")
		}
		return err
	}

	if role != "admin" {
		return fmt.Errorf("admin rights required")
	}

	// Verify connection exists between requester and target user
	var isConnected bool
	err = s.db.QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM connections
			WHERE ((requester_id = ? AND receiver_id = ?) OR (requester_id = ? AND receiver_id = ?))
			AND status = 'accepted'
		)`, adminID, targetUserID, targetUserID, adminID).Scan(&isConnected)

	if err != nil {
		return err
	}

	if !isConnected {
		return fmt.Errorf("user is not in your connections")
	}

	_, err = s.db.ExecContext(ctx,
		"INSERT INTO group_members (group_id, user_id, role) VALUES (?, ?, 'member') ON DUPLICATE KEY UPDATE group_id=group_id",
		groupID, targetUserID)
	if err != nil {
		return err
	}

	// Get group name synchronously
	var groupName string
	if err := s.db.QueryRowContext(ctx, "SELECT name FROM `groups` WHERE id = ?", groupID).Scan(&groupName); err != nil {
		groupName = "a group"
	}

	// Send notification in background? Or synchronous?
	// Handler did it in background `go func()`. Service typically does it synchronously or spawns goroutine.
	// Spawning goroutine in service is okay, but cleaner if caller handles it or if service uses a worker.
	// For now, mirroring handler behavior but being careful with context.
	// Using a detached context for notification is safer if the request context is cancelled.

	go func() {
		notifyCtx := context.Background()
		err := s.notificationService.SendNotification(notifyCtx, targetUserID, "Group Invitation", "You have been added to "+groupName, map[string]string{
			"type": "group_invite",
			"id":   groupID,
		})
		if err != nil {
			slog.Error("Failed to send notification", "error", err)
		}
	}()

	return nil
}

func (s *groupService) GetOrCreateDirectGroup(ctx context.Context, userID, partnerID string) (*Group, bool, error) {
	if partnerID == userID {
		return nil, false, fmt.Errorf("cannot create direct group with yourself")
	}

	// 1. Check if connected
	var isConnected bool
	err := s.db.QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM connections
			WHERE ((requester_id = ? AND receiver_id = ?) OR (requester_id = ? AND receiver_id = ?))
			AND status = 'accepted'
		)`, userID, partnerID, partnerID, userID).Scan(&isConnected)
	if err != nil {
		return nil, false, err
	}
	if !isConnected {
		return nil, false, fmt.Errorf("user is not in your connections")
	}

	// 2. Check if Direct group exists
	var groupID string
	err = s.db.QueryRowContext(ctx, `
		SELECT g.id FROM `+"`groups`"+` g
		JOIN group_members gm1 ON g.id = gm1.group_id AND gm1.user_id = ?
		JOIN group_members gm2 ON g.id = gm2.group_id AND gm2.user_id = ?
		WHERE g.type = 'direct' LIMIT 1
	`, userID, partnerID).Scan(&groupID)

	if err == nil {
		return &Group{ID: groupID}, false, nil
	}

	if err != sql.ErrNoRows {
		return nil, false, err
	}

	// 3. Create Group
	var partnerName string
	err = s.db.QueryRowContext(ctx, "SELECT COALESCE(username, email) FROM users WHERE id=?", partnerID).Scan(&partnerName)
	if err != nil {
		return nil, false, models.ErrNotFound // Partner not found
	}

	var myName string
	err = s.db.QueryRowContext(ctx, "SELECT COALESCE(username, email) FROM users WHERE id=?", userID).Scan(&myName)
	if err != nil {
		return nil, false, err
	}

	groupName := "Direct: " + myName + " & " + partnerName

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = tx.Rollback() }()

	groupID = uuid.New().String()
	if _, err := tx.ExecContext(ctx,
		"INSERT INTO `groups` (id, name, type, created_by) VALUES (?, ?, 'direct', ?)",
		groupID, groupName, userID); err != nil {
		return nil, false, err
	}

	// Add members
	if _, err := tx.ExecContext(ctx, "INSERT INTO group_members (group_id, user_id, role) VALUES (?, ?, 'admin')", groupID, userID); err != nil {
		return nil, false, err
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO group_members (group_id, user_id, role) VALUES (?, ?, 'admin')", groupID, partnerID); err != nil {
		return nil, false, err
	}

	if err := tx.Commit(); err != nil {
		return nil, false, err
	}

	return &Group{ID: groupID}, true, nil
}

func (s *groupService) RemoveGroupMember(ctx context.Context, adminID, groupID, targetUserID string) error {
	// Verify admin role
	var role string
	err := s.db.QueryRowContext(ctx,
		"SELECT role FROM group_members WHERE group_id = ? AND user_id = ?",
		groupID, adminID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("access denied")
		}
		return err
	}

	if role != "admin" {
		return fmt.Errorf("admin rights required")
	}

	_, err = s.db.ExecContext(ctx,
		"DELETE FROM group_members WHERE group_id = ? AND user_id = ?",
		groupID, targetUserID)
	return err
}
