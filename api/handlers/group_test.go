package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"discipleship_journal_api/middleware"
	"discipleship_journal_api/models"
	"discipleship_journal_api/services"

	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGroupHandler_CreateGroup(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    string
		setupMock      func(mockService *MockGroupService)
		expectedStatus int
		skipAuth       bool
	}{
		{
			name:        "Success",
			requestBody: `{"name": "Bible Study", "description": "Weekly study"}`,
			setupMock: func(mockService *MockGroupService) {
				desc := "Weekly study"
				mockService.On("CreateGroup", mock.Anything, mock.AnythingOfType("string"), "Bible Study", &desc, "").
					Return(&services.Group{ID: "00000000-0000-0000-0000-000000000123"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid Request",
			requestBody:    `{"name": ""}`, // too short
			setupMock:      func(_ *MockGroupService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Service Error",
			requestBody: `{"name": "Bible Study", "description": "Weekly study"}`,
			setupMock: func(mockService *MockGroupService) {
				desc := "Weekly study"
				mockService.On("CreateGroup", mock.Anything, mock.AnythingOfType("string"), "Bible Study", &desc, "").
					Return((*services.Group)(nil), errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Unauthorized",
			requestBody:    `{"name": "Bible Study"}`,
			setupMock:      func(_ *MockGroupService) {},
			expectedStatus: http.StatusUnauthorized,
			skipAuth:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockGroupService)
			h := NewGroupHandler(mockService)

			tc.setupMock(mockService)

			req := httptest.NewRequest("POST", "/groups", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			if !tc.skipAuth {
				testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
				ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
				dummyToken := &auth.Token{UID: "firebase-uid-123"}
				ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.CreateGroup(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_ListMyGroups(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		desc := "Desc 1"
		mockService.On("ListUserGroups", mock.Anything, mock.AnythingOfType("string")).
			Return([]services.Group{{ID: "1", Name: "Group 1", Description: &desc}}, nil)

		req := httptest.NewRequest("GET", "/groups", nil)
		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.ListMyGroups(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var groups []GroupResponse
		err := json.Unmarshal(w.Body.Bytes(), &groups)
		assert.NoError(t, err)
		assert.Len(t, groups, 1)
		assert.Equal(t, "Group 1", groups[0].Name)
		mockService.AssertExpectations(t)
	})

	t.Run("DBError", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("ListUserGroups", mock.Anything, mock.AnythingOfType("string")).
			Return([]services.Group(nil), errors.New("db error"))

		req := httptest.NewRequest("GET", "/groups", nil)
		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.ListMyGroups(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGroupHandler_SearchGroups(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("SearchGroups", mock.Anything, "Bible", mock.AnythingOfType("string")).
			Return([]services.Group{{ID: "1", Name: "Bible Study"}}, nil)

		req := httptest.NewRequest("GET", "/groups/search?q=Bible", nil)
		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.SearchGroups(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Query Too Short", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		req := httptest.NewRequest("GET", "/groups/search?q=Bi", nil)
		w := httptest.NewRecorder()
		h.SearchGroups(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("SearchGroups", mock.Anything, "Bible", mock.AnythingOfType("string")).
			Return([]services.Group(nil), errors.New("db error"))

		req := httptest.NewRequest("GET", "/groups/search?q=Bible", nil)
		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.SearchGroups(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGroupHandler_JoinGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("JoinGroup", mock.Anything, "group-1", mock.AnythingOfType("string")).Return(nil)

		req := httptest.NewRequest("POST", "/groups/group-1/join", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.JoinGroup(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Already Member", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("JoinGroup", mock.Anything, "group-1", mock.AnythingOfType("string")).Return(models.ErrAlreadyExists)

		req := httptest.NewRequest("POST", "/groups/group-1/join", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.JoinGroup(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("GenericError", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("JoinGroup", mock.Anything, "group-1", mock.AnythingOfType("string")).Return(errors.New("generic error"))

		req := httptest.NewRequest("POST", "/groups/group-1/join", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.JoinGroup(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGroupHandler_LeaveGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("LeaveGroup", mock.Anything, "group-1", mock.AnythingOfType("string")).Return(nil)

		req := httptest.NewRequest("POST", "/groups/group-1/leave", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.LeaveGroup(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Not Member", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("LeaveGroup", mock.Anything, "group-1", mock.AnythingOfType("string")).Return(models.ErrNotFound)

		req := httptest.NewRequest("POST", "/groups/group-1/leave", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.LeaveGroup(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("GenericError", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("LeaveGroup", mock.Anything, "group-1", mock.AnythingOfType("string")).Return(errors.New("generic error"))

		req := httptest.NewRequest("POST", "/groups/group-1/leave", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.LeaveGroup(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGroupHandler_GetGroupMembers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("GetGroupMembers", mock.Anything, "group-1", mock.AnythingOfType("string")).
			Return([]services.GroupMember{{UserID: "u2", Role: "member", JoinedAt: time.Now()}}, nil)

		req := httptest.NewRequest("GET", "/groups/group-1/members", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetGroupMembers(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("AccessDenied", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("GetGroupMembers", mock.Anything, "group-1", mock.AnythingOfType("string")).
			Return(nil, errors.New("access denied"))

		req := httptest.NewRequest("GET", "/groups/group-1/members", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetGroupMembers(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("DBError", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("GetGroupMembers", mock.Anything, "group-1", mock.AnythingOfType("string")).
			Return(nil, errors.New("db error"))

		req := httptest.NewRequest("GET", "/groups/group-1/members", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetGroupMembers(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGroupHandler_AddGroupMember(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("AddGroupMember", mock.Anything, mock.AnythingOfType("string"), "group-1", "user-2").Return(nil)

		reqBody := `{"user_id": "user-2"}`
		req := httptest.NewRequest("POST", "/groups/group-1/members", strings.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.AddGroupMember(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("AddGroupMember", mock.Anything, mock.AnythingOfType("string"), "group-1", "user-2").
			Return(errors.New("admin rights required"))

		reqBody := `{"user_id": "user-2"}`
		req := httptest.NewRequest("POST", "/groups/group-1/members", strings.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.AddGroupMember(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("NotConnected", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("AddGroupMember", mock.Anything, mock.AnythingOfType("string"), "group-1", "user-2").
			Return(errors.New("user is not in your connections"))

		reqBody := `{"user_id": "user-2"}`
		req := httptest.NewRequest("POST", "/groups/group-1/members", strings.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.AddGroupMember(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("GenericError", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("AddGroupMember", mock.Anything, mock.AnythingOfType("string"), "group-1", "user-2").
			Return(errors.New("db error"))

		reqBody := `{"user_id": "user-2"}`
		req := httptest.NewRequest("POST", "/groups/group-1/members", strings.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.AddGroupMember(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGroupHandler_GetOrCreateDirectGroup(t *testing.T) {
	testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	t.Run("Success - Created", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("GetOrCreateDirectGroup", mock.Anything, mock.AnythingOfType("string"), "partner-1").
			Return(&services.Group{ID: "group-direct"}, true, nil)

		reqBody := `{"partner_id": "partner-1"}`
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Success - Existing", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("GetOrCreateDirectGroup", mock.Anything, mock.AnythingOfType("string"), "partner-1").
			Return(&services.Group{ID: "group-direct"}, false, nil)

		reqBody := `{"partner_id": "partner-1"}`
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Request", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		reqBody := `{"partner_id": ""}` // invalid
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("GetOrCreateDirectGroup", mock.Anything, mock.AnythingOfType("string"), "partner-1").
			Return((*services.Group)(nil), false, errors.New("db error"))

		reqBody := `{"partner_id": "partner-1"}`
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("SelfError", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("GetOrCreateDirectGroup", mock.Anything, mock.AnythingOfType("string"), "partner-1").
			Return((*services.Group)(nil), false, errors.New("cannot create direct group with yourself"))

		reqBody := `{"partner_id": "partner-1"}`
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("NotConnected", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("GetOrCreateDirectGroup", mock.Anything, mock.AnythingOfType("string"), "partner-1").
			Return((*services.Group)(nil), false, errors.New("user is not in your connections"))

		reqBody := `{"partner_id": "partner-1"}`
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("PartnerNotFound", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("GetOrCreateDirectGroup", mock.Anything, mock.AnythingOfType("string"), "partner-1").
			Return((*services.Group)(nil), false, models.ErrNotFound)

		reqBody := `{"partner_id": "partner-1"}`
		req := httptest.NewRequest("POST", "/groups/direct", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetOrCreateDirectGroup(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGroupHandler_RemoveGroupMember(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("RemoveGroupMember", mock.Anything, mock.AnythingOfType("string"), "group-1", "user-2").Return(nil)

		req := httptest.NewRequest("DELETE", "/groups/group-1/members/user-2", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		rctx.URLParams.Add("userId", "user-2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		dummyToken := &auth.Token{UID: "firebase-uid-123"}
		ctx = context.WithValue(ctx, middleware.UserContextKey, dummyToken)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.RemoveGroupMember(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("AdminRequired", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("RemoveGroupMember", mock.Anything, mock.AnythingOfType("string"), "group-1", "user-2").
			Return(errors.New("admin rights required"))

		req := httptest.NewRequest("DELETE", "/groups/group-1/members/user-2", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		rctx.URLParams.Add("userId", "user-2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.RemoveGroupMember(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("GenericError", func(t *testing.T) {
		mockService := new(MockGroupService)
		h := NewGroupHandler(mockService)

		mockService.On("RemoveGroupMember", mock.Anything, mock.AnythingOfType("string"), "group-1", "user-2").
			Return(errors.New("db error"))

		req := httptest.NewRequest("DELETE", "/groups/group-1/members/user-2", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "group-1")
		rctx.URLParams.Add("userId", "user-2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		testUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx := context.WithValue(req.Context(), TestUserKey, testUUID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.RemoveGroupMember(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
