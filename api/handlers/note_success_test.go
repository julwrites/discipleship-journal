package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"discipleship_journal_api/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNoteHandler_CreateNote_Success(t *testing.T) {
	mockService := new(MockNoteService)
	handler := NewNoteHandler(nil, mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	reqBody := `{"title":"Test Note", "content":{"text":"hello"}, "tags":["tag1"]}`
	req := httptest.NewRequest("POST", "/api/notes", strings.NewReader(reqBody))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	expectedNote := &services.Note{
		ID:        "note-123",
		UserID:    userID.String(),
		Title:     "Test Note",
		Content:   json.RawMessage(`{"text":"hello"}`),
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tags: []services.Tag{
			{ID: "tag-1", Name: "tag1"},
		},
	}

	// Expect CreateNote with correct args.
	// Note: content is json.RawMessage, need to match it. tags is []string.
	// The variadic arg 'status' is matched by the last mock.Anything.
	mockService.On("CreateNote",
		mock.Anything,
		userID.String(),
		"Test Note",
		mock.MatchedBy(func(c json.RawMessage) bool {
			return string(c) == `{"text":"hello"}`
		}),
		[]string{"tag1"},
		mock.Anything).Return(expectedNote, nil)

	handler.CreateNote(w, req)

	// Note: Handler currently returns 200 OK, not 201 Created
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "note-123", resp["id"])

	mockService.AssertExpectations(t)
}
