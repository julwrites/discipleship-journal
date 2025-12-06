package middleware

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type AuthMiddleware struct {
	AuthClient *auth.Client
}

func NewAuthMiddleware(ctx context.Context) (*AuthMiddleware, error) {
	// If SERVICE_ACCOUNT_KEY env is set, use it (Local Dev).
	// Otherwise, use Application Default Credentials (Cloud Run).
	var app *firebase.App
	var err error

	// For simplicity in this sandbox, we might not have a service account key.
	// In a real environment, we'd initialize with options or default creds.
	// For now, let's assume default credentials or a mock for testing if needed.

	// Check if we are in a test/sandbox environment without creds
	// This is a placeholder. In production, use os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") or similar.

	config := &firebase.Config{ProjectID: "discipleship-journal-pwa"} // Replace with actual project ID if known, or leave empty to auto-detect

	app, err = firebase.NewApp(ctx, config)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{AuthClient: client}, nil
}

// InitAuthMiddleware initializes the middleware with options.
// This allows passing a service account JSON if available locally.
func InitAuthMiddleware(ctx context.Context, saKey string) (*AuthMiddleware, error) {
	var opts []option.ClientOption
	if saKey != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(saKey)))
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
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
