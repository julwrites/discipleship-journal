package handlers

import (
	"context"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
	"discipleship_journal_api/middleware"

	"github.com/stretchr/testify/assert"
)

func TestTemplateHandler_getUserID_NoToken(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	req := httptest.NewRequest("GET", "/", nil)

	_, err := handler.getUserID(req)
	assert.Error(t, err)
}

func TestTemplateHandler_getUserID_WithToken_NoUUID(t *testing.T) {
	mockService := new(MockTemplateService)
	handler := NewTemplateHandler(mockService)

	req := httptest.NewRequest("GET", "/", nil)

	// Valid token but no TestUserKey -> GetUserUUID should fail or return error if not mocked context
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, &auth.Token{UID: "firebase-uid"})
	req = req.WithContext(ctx)

	_, err := handler.getUserID(req)
	assert.Error(t, err) // because DB call or mock isn't provided
}
