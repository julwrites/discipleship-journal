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

func (m *MockNoteService) CreateNote(ctx context.Context, userID, title string, content json.RawMessage, status ...string) (*services.Note, error) {
	// Variadic args are passed as a slice to Called if we pass them explicitly, or we can just pass them individually?
	// testify/mock handles variadic by collecting them.
	// But m.Called() needs to receive them.
	// Let's pass the status slice as a single argument to Called for simplicity in expectations?
	// Or we can expand them. Ideally we pass them as is.
	// Common pattern: pass variadic args to Called
	args := m.Called(ctx, userID, title, content, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Note), args.Error(1)
}

func (m *MockNoteService) DeleteNote(ctx context.Context, userID, noteID string) error {
	args := m.Called(ctx, userID, noteID)
	return args.Error(0)
}

func (m *MockNoteService) UpdateNote(ctx context.Context, userID, noteID, title string, content json.RawMessage, status ...string) error {
	args := m.Called(ctx, userID, noteID, title, content, status)
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
