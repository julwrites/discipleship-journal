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

func (m *MockNoteService) CreateNote(ctx context.Context, userID, title string, content json.RawMessage) (*services.Note, error) {
	args := m.Called(ctx, userID, title, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Note), args.Error(1)
}

func (m *MockNoteService) DeleteNote(ctx context.Context, userID, noteID string) error {
	args := m.Called(ctx, userID, noteID)
	return args.Error(0)
}

func (m *MockNoteService) UpdateNote(ctx context.Context, userID, noteID, title string, content json.RawMessage) error {
	args := m.Called(ctx, userID, noteID, title, content)
	return args.Error(0)
}

func (m *MockNoteService) GetNote(ctx context.Context, userID, noteID string) (*services.Note, error) {
	args := m.Called(ctx, userID, noteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Note), args.Error(1)
}

func (m *MockNoteService) GetNotes(ctx context.Context, userID string, page, limit int, searchQuery string) ([]services.Note, int, error) {
	args := m.Called(ctx, userID, page, limit, searchQuery)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]services.Note), args.Int(1), args.Error(2)
}
