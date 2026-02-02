package handlers

import (
	"context"
	"discipleship_journal_api/services"
	"encoding/json"
	"github.com/stretchr/testify/mock"
)

type MockNoteService struct {
	mock.Mock
}

func (m *MockNoteService) CreateNote(ctx context.Context, userID, title string, content json.RawMessage, tags []string, status ...string) (*services.Note, error) {
	args := m.Called(ctx, userID, title, content, tags, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Note), args.Error(1)
}

func (m *MockNoteService) DeleteNote(ctx context.Context, userID, noteID string) error {
	args := m.Called(ctx, userID, noteID)
	return args.Error(0)
}

func (m *MockNoteService) UpdateNote(ctx context.Context, userID, noteID, title string, content json.RawMessage, tags []string, status ...string) error {
	args := m.Called(ctx, userID, noteID, title, content, tags, status)
	return args.Error(0)
}

func (m *MockNoteService) CreateTag(ctx context.Context, userID, name string) (*services.Tag, error) {
	args := m.Called(ctx, userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Tag), args.Error(1)
}

func (m *MockNoteService) GetUserTags(ctx context.Context, userID string) ([]services.Tag, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]services.Tag), args.Error(1)
}

func (m *MockNoteService) DeleteTag(ctx context.Context, userID, tagID string) error {
	args := m.Called(ctx, userID, tagID)
	return args.Error(0)
}

func (m *MockNoteService) GetNote(ctx context.Context, userID, noteID string) (*services.Note, error) {
	args := m.Called(ctx, userID, noteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Note), args.Error(1)
}

func (m *MockNoteService) GetNotes(ctx context.Context, userID string, page, limit int, filter services.NoteFilter) ([]services.Note, int, error) {
	args := m.Called(ctx, userID, page, limit, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]services.Note), args.Int(1), args.Error(2)
}

type MockGroupService struct {
	mock.Mock
}

func (m *MockGroupService) CreateGroup(ctx context.Context, userID, name string, description *string, groupType string) (*services.Group, error) {
	args := m.Called(ctx, userID, name, description, groupType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Group), args.Error(1)
}

func (m *MockGroupService) ListUserGroups(ctx context.Context, userID string) ([]services.Group, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]services.Group), args.Error(1)
}

func (m *MockGroupService) SearchGroups(ctx context.Context, query, userID string) ([]services.Group, error) {
	args := m.Called(ctx, query, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]services.Group), args.Error(1)
}

func (m *MockGroupService) JoinGroup(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}

func (m *MockGroupService) LeaveGroup(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}

func (m *MockGroupService) GetGroupMembers(ctx context.Context, groupID, userID string) ([]services.GroupMember, error) {
	args := m.Called(ctx, groupID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]services.GroupMember), args.Error(1)
}

func (m *MockGroupService) AddGroupMember(ctx context.Context, adminID, groupID, targetUserID string) error {
	args := m.Called(ctx, adminID, groupID, targetUserID)
	return args.Error(0)
}

func (m *MockGroupService) GetOrCreateDirectGroup(ctx context.Context, userID, partnerID string) (*services.Group, bool, error) {
	args := m.Called(ctx, userID, partnerID)
	if args.Get(0) == nil {
		return nil, false, args.Error(2)
	}
	return args.Get(0).(*services.Group), args.Bool(1), args.Error(2)
}

func (m *MockGroupService) RemoveGroupMember(ctx context.Context, adminID, groupID, targetUserID string) error {
	args := m.Called(ctx, adminID, groupID, targetUserID)
	return args.Error(0)
}
