package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
)

// SecretManagerClient defines the interface for interacting with Google Secret Manager.
// This allows for mocking the client in tests.
type SecretManagerClient interface {
	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
	Close() error
}

// SecretLoader handles loading secrets from Google Secret Manager with fallback to environment variables
type SecretLoader struct {
	projectID string
	client    SecretManagerClient
}

// NewSecretLoader creates a new SecretLoader
// If projectID is empty, it will try to read from GOOGLE_CLOUD_PROJECT environment variable
// If client is nil, Secret Manager operations will be skipped (fallback to env vars only)
func NewSecretLoader(ctx context.Context, projectID string) (*SecretLoader, error) {
	if projectID == "" {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
		if projectID == "" {
			// If no project ID, we can't use Secret Manager
			return &SecretLoader{projectID: "", client: nil}, nil
		}
	}

	// Try to create Secret Manager client
	// In production (Cloud Run), default credentials should work
	// For local development, you might need GOOGLE_APPLICATION_CREDENTIALS
	client, err := secretmanager.NewClient(ctx, option.WithUserAgent("discipleship-journal-api"))
	if err != nil {
		// If we can't create client, fall back to env vars only
		// This is expected in local development without credentials
		return &SecretLoader{projectID: projectID, client: nil}, nil
	}

	return &SecretLoader{
		projectID: projectID,
		client:    client,
	}, nil
}

// Close closes the Secret Manager client if it exists
func (sl *SecretLoader) Close() error {
	if sl.client != nil {
		return sl.client.Close()
	}
	return nil
}

// LoadSecret tries to load a secret from Google Secret Manager first,
// then falls back to environment variable if Secret Manager fails or is unavailable
func (sl *SecretLoader) LoadSecret(ctx context.Context, secretName string) (string, error) {
	// First try Secret Manager if client is available
	if sl.client != nil && sl.projectID != "" {
		secret, err := sl.loadFromSecretManager(ctx, secretName)
		if err == nil && secret != "" {
			return secret, nil
		}
		// If Secret Manager fails, fall through to env var
	}

	// Fallback to environment variable
	envVar := strings.ToUpper(secretName)
	value := os.Getenv(envVar)
	if value == "" {
		return "", fmt.Errorf("secret %q not found in Secret Manager or environment variable %q", secretName, envVar)
	}

	return value, nil
}

// loadFromSecretManager loads a secret from Google Secret Manager
func (sl *SecretLoader) loadFromSecretManager(ctx context.Context, secretName string) (string, error) {
	if sl.client == nil || sl.projectID == "" {
		return "", fmt.Errorf("secret manager client not initialized")
	}

	// Build the resource name
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", sl.projectID, secretName)

	// Access the secret version
	result, err := sl.client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to access secret %q: %w", secretName, err)
	}

	return string(result.Payload.Data), nil
}

// MustLoadSecret loads a secret or panics if it cannot be found
// Use only during application initialization for critical secrets
func (sl *SecretLoader) MustLoadSecret(ctx context.Context, secretName string) string {
	value, err := sl.LoadSecret(ctx, secretName)
	if err != nil {
		panic(fmt.Sprintf("Failed to load required secret %q: %v", secretName, err))
	}
	return value
}
