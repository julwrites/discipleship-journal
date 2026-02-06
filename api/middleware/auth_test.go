package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"firebase.google.com/go/v4/auth"
)

// MockFirebaseAuthClient implements FirebaseAuthClient interface
type MockFirebaseAuthClient struct {
	ShouldError bool
	Token       *auth.Token
}

func (m *MockFirebaseAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if m.ShouldError {
		return nil, errors.New("mock error")
	}
	if idToken == "invalid-token" {
		return nil, errors.New("invalid token")
	}
	return m.Token, nil
}

func TestAuthMiddleware_NoHeader(t *testing.T) {
	am := &AuthMiddleware{AuthClient: &MockFirebaseAuthClient{}}

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := am.VerifyToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	expected := "Authorization header required\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned wrong body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	mockToken := &auth.Token{
		UID: "test-user-id",
	}
	am := &AuthMiddleware{
		AuthClient: &MockFirebaseAuthClient{
			Token: mockToken,
		},
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rr := httptest.NewRecorder()

	handler := am.VerifyToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if token is in context
		token := r.Context().Value(UserContextKey).(*auth.Token)
		if token.UID != "test-user-id" {
			t.Errorf("Token not correctly set in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	am := &AuthMiddleware{
		AuthClient: &MockFirebaseAuthClient{
			ShouldError: true,
		},
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()

	handler := am.VerifyToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	if !strings.Contains(rr.Body.String(), "Invalid token") {
		t.Errorf("handler returned unexpected body: %v", rr.Body.String())
	}
}

func TestNewAuthMiddlewareFromClient(t *testing.T) {
	mockClient := &MockFirebaseAuthClient{}
	am := NewAuthMiddlewareFromClient(mockClient)
	if am == nil {
		t.Fatal("NewAuthMiddlewareFromClient returned nil")
	}
	if am.AuthClient != mockClient {
		t.Error("AuthClient not set correctly")
	}
}
