package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"discipleship_journal_api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGroupHandler_CreateGroup_Success(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	reqBody := `{"name":"Test Group", "description":"Test Desc", "type":"group"}`
	req := httptest.NewRequest("POST", "/api/groups", strings.NewReader(reqBody))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	expectedGroup := &services.Group{
		ID:          "group-123",
		Name:        "Test Group",
		Description: func() *string { s := "Test Desc"; return &s }(),
		Type:        "group",
	}

	mockService.On("CreateGroup", mock.Anything, userID.String(), "Test Group", mock.Anything, "group").Return(expectedGroup, nil)

	handler.CreateGroup(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "group-123", resp["id"])

	mockService.AssertExpectations(t)
}
