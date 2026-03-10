package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGroupHandler_CreateGroup_Unauthorized(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	req := httptest.NewRequest("POST", "/api/groups", strings.NewReader(`{"name":"Test Group"}`))
	// No user context
	w := httptest.NewRecorder()

	handler.CreateGroup(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGroupHandler_ListMyGroups_Unauthorized(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	req := httptest.NewRequest("GET", "/api/groups", nil)
	// No user context
	w := httptest.NewRecorder()

	handler.ListMyGroups(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGroupHandler_SearchGroups_Unauthorized(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	req := httptest.NewRequest("GET", "/api/groups/search?q=query", nil)
	// No user context
	w := httptest.NewRecorder()

	handler.SearchGroups(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code) // GetUserUUIDFromContext returns error which returns 500 in SearchGroups
}

func TestGroupHandler_SearchGroups_ShortQuery(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	req := httptest.NewRequest("GET", "/api/groups/search?q=ab", nil)
	w := httptest.NewRecorder()

	handler.SearchGroups(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGroupHandler_JoinGroup_Unauthorized(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	groupID := "group-id"
	req := httptest.NewRequest("POST", "/api/groups/"+groupID+"/join", nil)
	// No user context
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/groups/{id}/join", handler.JoinGroup)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code) // GetUserUUIDFromContext returns error which returns 500 in JoinGroup
}

func TestGroupHandler_LeaveGroup_Unauthorized(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	groupID := "group-id"
	req := httptest.NewRequest("POST", "/api/groups/"+groupID+"/leave", nil)
	// No user context
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/groups/{id}/leave", handler.LeaveGroup)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code) // GetUserUUIDFromContext returns error which returns 500 in LeaveGroup
}

func TestGroupHandler_GetGroupMembers_Unauthorized(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	groupID := "group-id"
	req := httptest.NewRequest("GET", "/api/groups/"+groupID+"/members", nil)
	// No user context
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/api/groups/{id}/members", handler.GetGroupMembers)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code) // GetUserUUIDFromContext returns error which returns 500 in GetGroupMembers
}

func TestGroupHandler_AddGroupMember_Unauthorized(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	groupID := "group-id"
	req := httptest.NewRequest("POST", "/api/groups/"+groupID+"/members", strings.NewReader(`{"user_id":"u2"}`))
	// No user context
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/groups/{id}/members", handler.AddGroupMember)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code) // GetUserUUIDFromContext returns error which returns 500 in AddGroupMember
}

func TestGroupHandler_AddGroupMember_InvalidBody(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	groupID := "group-id"
	req := httptest.NewRequest("POST", "/api/groups/"+groupID+"/members", strings.NewReader(`{invalid}`))
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/groups/{id}/members", handler.AddGroupMember)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGroupHandler_GetOrCreateDirectGroup_Unauthorized(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	req := httptest.NewRequest("POST", "/api/groups/direct", strings.NewReader(`{"partner_id":"p1"}`))
	// No user context
	w := httptest.NewRecorder()

	handler.GetOrCreateDirectGroup(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code) // GetUserUUIDFromContext returns error which returns 500 in GetOrCreateDirectGroup
}

func TestGroupHandler_GetOrCreateDirectGroup_InvalidBody(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	req := httptest.NewRequest("POST", "/api/groups/direct", strings.NewReader(`{invalid}`))
	w := httptest.NewRecorder()

	handler.GetOrCreateDirectGroup(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGroupHandler_RemoveGroupMember_Unauthorized(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	groupID := "group-id"
	targetUserID := "u2"
	req := httptest.NewRequest("DELETE", "/api/groups/"+groupID+"/members/"+targetUserID, nil)
	// No user context
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/api/groups/{id}/members/{userId}", handler.RemoveGroupMember)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code) // GetUserUUIDFromContext returns error which returns 500 in RemoveGroupMember
}

func TestGroupHandler_CreateGroup_ServiceError(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), TestUserKey, userID.String())

	req := httptest.NewRequest("POST", "/api/groups", strings.NewReader(`{"name":"Test Group"}`))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.On("CreateGroup", mock.Anything, userID.String(), "Test Group", mock.Anything, "").Return(nil, assert.AnError)

	handler.CreateGroup(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGroupHandler_CreateGroup_InvalidBody(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	req := httptest.NewRequest("POST", "/api/groups", strings.NewReader(`{invalid}`))
	w := httptest.NewRecorder()

	handler.CreateGroup(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGroupHandler_CreateGroup_ValidationErr(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	req := httptest.NewRequest("POST", "/api/groups", strings.NewReader(`{"name":""}`)) // Empty name
	w := httptest.NewRecorder()

	handler.CreateGroup(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
