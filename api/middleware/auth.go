package middleware

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type AuthMiddleware struct {
	AuthClient *auth.Client
}

// NewAuthMiddlewareFromClient creates a new AuthMiddleware with an existing Auth Client.
func NewAuthMiddlewareFromClient(client *auth.Client) *AuthMiddleware {
	return &AuthMiddleware{AuthClient: client}
}

func NewAuthMiddleware(ctx context.Context) (*AuthMiddleware, error) {
	// Legacy or for testing simple init
	config := &firebase.Config{ProjectID: "discipleship-journal-pwa"}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return &AuthMiddleware{AuthClient: client}, nil
}

func (am *AuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := am.AuthClient.VerifyIDToken(r.Context(), tokenString)
		if err != nil {
			// In a real scenario, handle expired tokens, etc.
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
