package services

import (
	"context"
)

// MockGroupService is a mock implementation of GroupService for testing.
type MockGroupService struct {
	CreateGroupFunc          func(ctx context.Context, name, description, userID string) (string, error)
	ListMyGroupsFunc         func(ctx context.Context, userID string) ([]GroupResponse, error)
	SearchGroupsFunc         func(ctx context.Context, query, userID string) ([]GroupResponse, error)
	JoinGroupFunc            func(ctx context.Context, groupID, userID string) error
	LeaveGroupFunc           func(ctx context.Context, groupID, userID string) error
	GetGroupMembersFunc      func(ctx context.Context, groupID, userID string) ([]GroupMemberResponse, error)
	AddGroupMemberFunc       func(ctx context.Context, groupID, targetUserID, actorUserID string) error
	RemoveGroupMemberFunc    func(ctx context.Context, groupID, targetUserID, actorUserID string) error
	ShareNoteToGroupFunc     func(ctx context.Context, groupID, noteID, userID, comment string) error
	ListGroupSharesFunc      func(ctx context.Context, groupID, userID string) ([]SharedNoteResponse, error)
	GetSharedNoteDetailsFunc func(ctx context.Context, shareID, groupID, userID string) (*SharedNoteResponse, error)
}

func NewMockGroupService() *MockGroupService {
	return &MockGroupService{
		CreateGroupFunc: func(ctx context.Context, name, description, userID string) (string, error) {
			return "mock-group-id", nil
		},
		ListMyGroupsFunc: func(ctx context.Context, userID string) ([]GroupResponse, error) {
			return []GroupResponse{}, nil
		},
		SearchGroupsFunc: func(ctx context.Context, query, userID string) ([]GroupResponse, error) {
			return []GroupResponse{}, nil
		},
		JoinGroupFunc: func(ctx context.Context, groupID, userID string) error {
			return nil
		},
		LeaveGroupFunc: func(ctx context.Context, groupID, userID string) error {
			return nil
		},
		GetGroupMembersFunc: func(ctx context.Context, groupID, userID string) ([]GroupMemberResponse, error) {
			return []GroupMemberResponse{}, nil
		},
		AddGroupMemberFunc: func(ctx context.Context, groupID, targetUserID, actorUserID string) error {
			return nil
		},
		RemoveGroupMemberFunc: func(ctx context.Context, groupID, targetUserID, actorUserID string) error {
			return nil
		},
		ShareNoteToGroupFunc: func(ctx context.Context, groupID, noteID, userID, comment string) error {
			return nil
		},
		ListGroupSharesFunc: func(ctx context.Context, groupID, userID string) ([]SharedNoteResponse, error) {
			return []SharedNoteResponse{}, nil
		},
		GetSharedNoteDetailsFunc: func(ctx context.Context, shareID, groupID, userID string) (*SharedNoteResponse, error) {
			return &SharedNoteResponse{}, nil
		},
	}
}

func (m *MockGroupService) CreateGroup(ctx context.Context, name, description, userID string) (string, error) {
	if m.CreateGroupFunc != nil {
		return m.CreateGroupFunc(ctx, name, description, userID)
	}
	return "", nil
}

func (m *MockGroupService) ListMyGroups(ctx context.Context, userID string) ([]GroupResponse, error) {
	if m.ListMyGroupsFunc != nil {
		return m.ListMyGroupsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockGroupService) SearchGroups(ctx context.Context, query, userID string) ([]GroupResponse, error) {
	if m.SearchGroupsFunc != nil {
		return m.SearchGroupsFunc(ctx, query, userID)
	}
	return nil, nil
}

func (m *MockGroupService) JoinGroup(ctx context.Context, groupID, userID string) error {
	if m.JoinGroupFunc != nil {
		return m.JoinGroupFunc(ctx, groupID, userID)
	}
	return nil
}

func (m *MockGroupService) LeaveGroup(ctx context.Context, groupID, userID string) error {
	if m.LeaveGroupFunc != nil {
		return m.LeaveGroupFunc(ctx, groupID, userID)
	}
	return nil
}

func (m *MockGroupService) GetGroupMembers(ctx context.Context, groupID, userID string) ([]GroupMemberResponse, error) {
	if m.GetGroupMembersFunc != nil {
		return m.GetGroupMembersFunc(ctx, groupID, userID)
	}
	return nil, nil
}

func (m *MockGroupService) AddGroupMember(ctx context.Context, groupID, targetUserID, actorUserID string) error {
	if m.AddGroupMemberFunc != nil {
		return m.AddGroupMemberFunc(ctx, groupID, targetUserID, actorUserID)
	}
	return nil
}

func (m *MockGroupService) RemoveGroupMember(ctx context.Context, groupID, targetUserID, actorUserID string) error {
	if m.RemoveGroupMemberFunc != nil {
		return m.RemoveGroupMemberFunc(ctx, groupID, targetUserID, actorUserID)
	}
	return nil
}

func (m *MockGroupService) ShareNoteToGroup(ctx context.Context, groupID, noteID, userID, comment string) error {
	if m.ShareNoteToGroupFunc != nil {
		return m.ShareNoteToGroupFunc(ctx, groupID, noteID, userID, comment)
	}
	return nil
}

func (m *MockGroupService) ListGroupShares(ctx context.Context, groupID, userID string) ([]SharedNoteResponse, error) {
	if m.ListGroupSharesFunc != nil {
		return m.ListGroupSharesFunc(ctx, groupID, userID)
	}
	return nil, nil
}

func (m *MockGroupService) GetSharedNoteDetails(ctx context.Context, shareID, groupID, userID string) (*SharedNoteResponse, error) {
	if m.GetSharedNoteDetailsFunc != nil {
		return m.GetSharedNoteDetailsFunc(ctx, shareID, groupID, userID)
	}
	return nil, nil
}
