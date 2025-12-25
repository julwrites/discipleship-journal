package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"discipleship_journal_api/database"
)

type GroupResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
	Role        string `json:"role,omitempty"`
}

type GroupMemberResponse struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	JoinedAt    string `json:"joined_at"`
}

type SharedNoteResponse struct {
	ID       string                 `json:"id"` // This is the share ID
	GroupID  string                 `json:"group_id"`
	NoteID   string                 `json:"note_id"`
	Title    string                 `json:"title"`
	Content  map[string]interface{} `json:"content,omitempty"`
	SharedBy string                 `json:"shared_by"`
	SharedAt string                 `json:"shared_at"`
	Comment  string                 `json:"comment"`
}

var (
	ErrGroupNotFound    = errors.New("group not found")
	ErrAccessDenied     = errors.New("access denied")
	ErrAlreadyMember    = errors.New("already a member")
	ErrNotMember        = errors.New("not a member")
	ErrNoteNotFound     = errors.New("note not found")
	ErrNoteNotOwned     = errors.New("you can only share your own notes")
	ErrShareNotFound    = errors.New("shared note not found")
	ErrInsufficientPerms = errors.New("admin rights required")
)

type GroupService interface {
	CreateGroup(ctx context.Context, name, description, userID string) (string, error)
	ListMyGroups(ctx context.Context, userID string) ([]GroupResponse, error)
	SearchGroups(ctx context.Context, query, userID string) ([]GroupResponse, error)
	JoinGroup(ctx context.Context, groupID, userID string) error
	LeaveGroup(ctx context.Context, groupID, userID string) error
	GetGroupMembers(ctx context.Context, groupID, userID string) ([]GroupMemberResponse, error)
	AddGroupMember(ctx context.Context, groupID, targetUserID, actorUserID string) error
	RemoveGroupMember(ctx context.Context, groupID, targetUserID, actorUserID string) error

	ShareNoteToGroup(ctx context.Context, groupID, noteID, userID, comment string) error
	ListGroupShares(ctx context.Context, groupID, userID string) ([]SharedNoteResponse, error)
	GetSharedNoteDetails(ctx context.Context, shareID, groupID, userID string) (*SharedNoteResponse, error)
}

type groupService struct {
	db                  database.DBInterface
	notificationService NotificationService
}

func NewGroupService(db database.DBInterface, notificationService NotificationService) GroupService {
	return &groupService{db: db, notificationService: notificationService}
}

func (s *groupService) CreateGroup(ctx context.Context, name, description, userID string) (string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("Failed to rollback transaction", "error", err)
		}
	}()

	var groupID string
	err = tx.QueryRow(ctx,
		"INSERT INTO groups (name, description, created_by) VALUES ($1, $2, $3) RETURNING id",
		name, description, userID).Scan(&groupID)
	if err != nil {
		return "", err
	}

	// Add creator as admin
	_, err = tx.Exec(ctx,
		"INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'admin')",
		groupID, userID)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return groupID, nil
}

func (s *groupService) ListMyGroups(ctx context.Context, userID string) ([]GroupResponse, error) {
	rows, err := s.db.Query(ctx,
		`SELECT g.id, g.name, g.description, g.created_by, gm.role
		 FROM groups g
		 JOIN group_members gm ON g.id = gm.group_id
		 WHERE gm.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []GroupResponse
	for rows.Next() {
		var g GroupResponse
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.CreatedBy, &g.Role); err != nil {
			continue
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (s *groupService) SearchGroups(ctx context.Context, query, userID string) ([]GroupResponse, error) {
	rows, err := s.db.Query(ctx,
		`SELECT g.id, g.name, g.description, g.created_by,
		 COALESCE((SELECT role FROM group_members WHERE group_id = g.id AND user_id = $2), '') as role
		 FROM groups g
		 WHERE g.name ILIKE $1 LIMIT 20`, "%"+query+"%", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []GroupResponse
	for rows.Next() {
		var g GroupResponse
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.CreatedBy, &g.Role); err != nil {
			continue
		}
		groups = append(groups, g)
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
		return ErrAlreadyMember
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
		return ErrNotMember
	}
	return nil
}

func (s *groupService) GetGroupMembers(ctx context.Context, groupID, userID string) ([]GroupMemberResponse, error) {
	// Check if user is a member
	var isMember bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userID).Scan(&isMember)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, ErrAccessDenied
	}

	rows, err := s.db.Query(ctx,
		`SELECT gm.user_id, u.display_name, u.email, gm.role, gm.joined_at
		 FROM group_members gm
		 JOIN users u ON gm.user_id = u.id
		 WHERE gm.group_id = $1`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []GroupMemberResponse
	for rows.Next() {
		var m GroupMemberResponse
		var joinedAt time.Time
		if err := rows.Scan(&m.UserID, &m.DisplayName, &m.Email, &m.Role, &joinedAt); err != nil {
			continue
		}
		m.JoinedAt = joinedAt.Format(time.RFC3339)
		members = append(members, m)
	}
	return members, nil
}

func (s *groupService) AddGroupMember(ctx context.Context, groupID, targetUserID, actorUserID string) error {
	// Verify admin role
	var role string
	err := s.db.QueryRow(ctx,
		"SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, actorUserID).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccessDenied
		}
		return err
	}

	if role != "admin" {
		return ErrInsufficientPerms
	}

	// Insert member
	_, err = s.db.Exec(ctx,
		"INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'member') ON CONFLICT DO NOTHING",
		groupID, targetUserID)
	if err != nil {
		return err
	}

	// Get group name for notification
	var groupName string
	if err := s.db.QueryRow(ctx, "SELECT name FROM groups WHERE id = $1", groupID).Scan(&groupName); err != nil {
		groupName = "a group"
	}

	// Send notification
	go func() {
		ctx := context.Background()
		err := s.notificationService.SendNotification(ctx, targetUserID, "Group Invitation", "You have been added to "+groupName, map[string]string{
			"type": "group_invite",
			"id":   groupID,
		})
		if err != nil {
			slog.Error("Failed to send notification", "error", err)
		}
	}()

	return nil
}

func (s *groupService) RemoveGroupMember(ctx context.Context, groupID, targetUserID, actorUserID string) error {
	// Verify admin role
	var role string
	err := s.db.QueryRow(ctx,
		"SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, actorUserID).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccessDenied
		}
		return err
	}

	if role != "admin" {
		return ErrInsufficientPerms
	}

	_, err = s.db.Exec(ctx,
		"DELETE FROM group_members WHERE group_id = $1 AND user_id = $2",
		groupID, targetUserID)
	return err
}

func (s *groupService) ShareNoteToGroup(ctx context.Context, groupID, noteID, userID, comment string) error {
	// 1. Verify membership
	var isMember bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userID).Scan(&isMember)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrAccessDenied
	}

	// 2. Verify ownership
	var ownerID string
	err = s.db.QueryRow(ctx,
		"SELECT user_id FROM notes WHERE id = $1 AND deleted_at IS NULL",
		noteID).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoteNotFound
		}
		return err
	}
	if ownerID != userID {
		return ErrNoteNotOwned
	}

	// 3. Create share
	_, err = s.db.Exec(ctx,
		"INSERT INTO group_shares (group_id, note_id, shared_by, comment) VALUES ($1, $2, $3, $4) ON CONFLICT (group_id, note_id) DO UPDATE SET shared_at = NOW(), comment = $4",
		groupID, noteID, userID, comment)
	if err != nil {
		return err
	}

	// Notification Logic
	go func() {
		// Use a fresh context for async operation
		ctx := context.Background()

		var groupName string
		if err := s.db.QueryRow(ctx, "SELECT name FROM groups WHERE id = $1", groupID).Scan(&groupName); err != nil {
			groupName = "Group"
		}

		var sharerName string
		if err := s.db.QueryRow(ctx, "SELECT display_name FROM users WHERE id = $1", userID).Scan(&sharerName); err != nil {
			sharerName = "Someone"
		}

		rows, err := s.db.Query(ctx, "SELECT user_id FROM group_members WHERE group_id = $1 AND user_id != $2", groupID, userID)
		if err != nil {
			slog.Error("Failed to get members for notification", "error", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var memberID string
			if err := rows.Scan(&memberID); err == nil {
				err := s.notificationService.SendNotification(ctx, memberID, "New Shared Note", sharerName+" shared a note in "+groupName, map[string]string{
					"type":     "note_share",
					"group_id": groupID,
					"note_id":  noteID,
				})
				if err != nil {
					slog.Error("Failed to send notification", "user_id", memberID, "error", err)
				}
			}
		}
	}()

	return nil
}

func (s *groupService) ListGroupShares(ctx context.Context, groupID, userID string) ([]SharedNoteResponse, error) {
	// Verify membership
	var isMember bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userID).Scan(&isMember)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrAccessDenied
	}

	rows, err := s.db.Query(ctx,
		`SELECT gs.id, gs.group_id, gs.note_id, n.title, u.display_name, gs.shared_at, gs.comment
		 FROM group_shares gs
		 JOIN notes n ON gs.note_id = n.id
		 JOIN users u ON gs.shared_by = u.id
		 WHERE gs.group_id = $1 AND n.deleted_at IS NULL
		 ORDER BY gs.shared_at DESC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []SharedNoteResponse
	for rows.Next() {
		var s SharedNoteResponse
		var sharedAt time.Time
		if err := rows.Scan(&s.ID, &s.GroupID, &s.NoteID, &s.Title, &s.SharedBy, &sharedAt, &s.Comment); err != nil {
			continue
		}
		s.SharedAt = sharedAt.Format(time.RFC3339)
		shares = append(shares, s)
	}
	return shares, nil
}

func (s *groupService) GetSharedNoteDetails(ctx context.Context, shareID, groupID, userID string) (*SharedNoteResponse, error) {
	// Verify membership
	var isMember bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)",
		groupID, userID).Scan(&isMember)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrAccessDenied
	}

	var resp SharedNoteResponse
	var sharedAt time.Time
	err = s.db.QueryRow(ctx,
		`SELECT gs.id, gs.group_id, gs.note_id, n.title, n.content, u.display_name, gs.shared_at, gs.comment
		 FROM group_shares gs
		 JOIN notes n ON gs.note_id = n.id
		 JOIN users u ON gs.shared_by = u.id
		 WHERE gs.id = $1 AND gs.group_id = $2 AND n.deleted_at IS NULL`, shareID, groupID).Scan(
		&resp.ID, &resp.GroupID, &resp.NoteID, &resp.Title, &resp.Content, &resp.SharedBy, &sharedAt, &resp.Comment)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrShareNotFound
		}
		return nil, err
	}
	resp.SharedAt = sharedAt.Format(time.RFC3339)
	return &resp, nil
}
