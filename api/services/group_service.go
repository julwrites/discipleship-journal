package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"discipleship_journal_api/models"
	"github.com/jackc/pgx/v5"
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

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	slog.Info("Creating group", "name", name, "description", description, "user", userID, "type", groupType)
	var groupID string
	err = tx.QueryRow(ctx,
		"INSERT INTO groups (name, description, created_by, type) VALUES ($1, $2, $3, $4) RETURNING id",
		name, description, userID, groupType).Scan(&groupID)
	if err != nil {
		slog.Error("Failed to create group", "error", err)
		return nil, err
	}

	// Add creator as admin
	_, err = tx.Exec(ctx,
		"INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'admin')",
		groupID, userID)
	if err != nil {
		slog.Error("Failed to add member", "error", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
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
	rows, err := s.db.Query(ctx,
		`SELECT g.id, g.name, g.description, g.created_by, g.type, gm.role
		 FROM groups g
		 JOIN group_members gm ON g.id = gm.group_id
		 WHERE gm.user_id = $1`, userID)
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
	rows, err := s.db.Query(ctx,
		`SELECT g.id, g.name, g.description, g.created_by, g.type,
		 COALESCE((SELECT role FROM group_members WHERE group_id = g.id AND user_id = $2), '') as role
		 FROM groups g
		 WHERE g.name ILIKE $1 AND (g.type = 'group' OR g.type IS NULL) LIMIT 20`, "%"+query+"%", userID)
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
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return models.ErrAlreadyExists
	}

	_, err = s.db.Exec(ctx,
		"INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'member')",
		groupID, userID)
	return err
}

func (s *groupService) LeaveGroup(ctx context.Context, groupID, userID string) error {
	result, err := s.db.Exec(ctx,
		"DELETE FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *groupService) GetGroupMembers(ctx context.Context, groupID, userID string) ([]GroupMember, error) {
	var isMember bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userID).Scan(&isMember)
	if err != nil {
		return nil, err
	}

	if !isMember {
		// Use a specific error for forbidden access?
		// For now returning generic error or we can define ErrForbidden
		return nil, fmt.Errorf("access denied")
	}

	rows, err := s.db.Query(ctx,
		`SELECT gm.user_id, COALESCE(u.username, u.email), u.email, gm.role, gm.joined_at
		 FROM group_members gm
		 JOIN users u ON gm.user_id = u.id
		 WHERE gm.group_id = $1`, groupID)
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
	err := s.db.QueryRow(ctx,
		"SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, adminID).Scan(&role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("access denied")
		}
		return err
	}

	if role != "admin" {
		return fmt.Errorf("admin rights required")
	}

	// Verify connection exists between requester and target user
	var isConnected bool
	err = s.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM connections
			WHERE ((requester_id = $1 AND receiver_id = $2) OR (requester_id = $2 AND receiver_id = $1))
			AND status = 'accepted'
		)`, adminID, targetUserID).Scan(&isConnected)

	if err != nil {
		return err
	}

	if !isConnected {
		return fmt.Errorf("user is not in your connections")
	}

	_, err = s.db.Exec(ctx,
		"INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'member') ON CONFLICT DO NOTHING",
		groupID, targetUserID)
	if err != nil {
		return err
	}

	// Get group name synchronously
	var groupName string
	if err := s.db.QueryRow(ctx, "SELECT name FROM groups WHERE id = $1", groupID).Scan(&groupName); err != nil {
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
	err := s.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM connections
			WHERE ((requester_id = $1 AND receiver_id = $2) OR (requester_id = $2 AND receiver_id = $1))
			AND status = 'accepted'
		)`, userID, partnerID).Scan(&isConnected)
	if err != nil {
		return nil, false, err
	}
	if !isConnected {
		return nil, false, fmt.Errorf("user is not in your connections")
	}

	// 2. Check if Direct group exists
	var groupID string
	err = s.db.QueryRow(ctx, `
		SELECT g.id FROM groups g
		JOIN group_members gm1 ON g.id = gm1.group_id AND gm1.user_id = $1
		JOIN group_members gm2 ON g.id = gm2.group_id AND gm2.user_id = $2
		WHERE g.type = 'direct' LIMIT 1
	`, userID, partnerID).Scan(&groupID)

	if err == nil {
		return &Group{ID: groupID}, false, nil
	}

	if err != pgx.ErrNoRows {
		return nil, false, err
	}

	// 3. Create Group
	var partnerName string
	err = s.db.QueryRow(ctx, "SELECT COALESCE(username, email) FROM users WHERE id=$1", partnerID).Scan(&partnerName)
	if err != nil {
		return nil, false, models.ErrNotFound // Partner not found
	}

	var myName string
	err = s.db.QueryRow(ctx, "SELECT COALESCE(username, email) FROM users WHERE id=$1", userID).Scan(&myName)
	if err != nil {
		return nil, false, err
	}

	groupName := "Direct: " + myName + " & " + partnerName

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := tx.QueryRow(ctx,
		"INSERT INTO groups (name, type, created_by) VALUES ($1, 'direct', $2) RETURNING id",
		groupName, userID).Scan(&groupID); err != nil {
		return nil, false, err
	}

	// Add members
	if _, err := tx.Exec(ctx, "INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'admin')", groupID, userID); err != nil {
		return nil, false, err
	}
	if _, err := tx.Exec(ctx, "INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'admin')", groupID, partnerID); err != nil {
		return nil, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, false, err
	}

	return &Group{ID: groupID}, true, nil
}

func (s *groupService) RemoveGroupMember(ctx context.Context, adminID, groupID, targetUserID string) error {
	// Verify admin role
	var role string
	err := s.db.QueryRow(ctx,
		"SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, adminID).Scan(&role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("access denied")
		}
		return err
	}

	if role != "admin" {
		return fmt.Errorf("admin rights required")
	}

	_, err = s.db.Exec(ctx,
		"DELETE FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, targetUserID)
	return err
}
