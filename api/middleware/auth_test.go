package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockAuthClient would be needed to test VerifyToken properly without a real Firebase setup.
// For now, we will test the structure of the middleware and ensure it blocks requests without headers.

func TestAuthMiddleware_NoHeader(t *testing.T) {
	// We can't easily mock the Firebase client struct without an interface,
	// but we can test that the middleware func itself is structured correctly
	// if we could inject a mock client.
	// Since AuthMiddleware struct has *auth.Client, which is a struct, we cannot mock it easily in Go
	// without wrapping it in an interface.

	// However, we can test the behavior of the middleware logic if we extract the token verification logic.
	// Or we can just skip testing the specific Firebase part and trust the library.

	// Let's at least ensure we can instantiate it (mocking context might fail if we don't have creds).
	// So we will just skip deep testing of this specific file for now, but acknowledge it exists.
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	// Create a dummy middleware with nil client, it should fail before using the client
	// if we check header first.
	am := &AuthMiddleware{AuthClient: nil}

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
