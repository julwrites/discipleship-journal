package handlers

import (
	"context"
	"encoding/json"
	"discipleship_journal_api/services"
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
