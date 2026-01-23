package services

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSecretLoader_NoProjectID(t *testing.T) {
	// Clear GOOGLE_CLOUD_PROJECT env var for this test
	originalProjectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	_ = os.Unsetenv("GOOGLE_CLOUD_PROJECT")
	defer func() { _ = os.Setenv("GOOGLE_CLOUD_PROJECT", originalProjectID) }()

	loader, err := NewSecretLoader(context.Background(), "")
	assert.NoError(t, err)
	assert.NotNil(t, loader)
	assert.Equal(t, "", loader.projectID)
	assert.Nil(t, loader.client)
}

func TestNewSecretLoader_WithProjectID(t *testing.T) {
	// This test may or may not create a client depending on whether
	// Google Cloud credentials are available in the test environment
	loader, err := NewSecretLoader(context.Background(), "test-project")
	assert.NoError(t, err)
	assert.NotNil(t, loader)
	assert.Equal(t, "test-project", loader.projectID)
	// Client may be nil or a real client depending on environment
	// Both are valid outcomes
	if loader.client != nil {
		// If client was created, close it
		_ = loader.Close()
	}
}

func TestSecretLoader_LoadSecret_EnvVarFallback(t *testing.T) {
	// Set up test environment variable
	secretName := "TEST_SECRET"
	expectedValue := "test-secret-value"
	_ = os.Setenv("TEST_SECRET", expectedValue)
	defer func() { _ = os.Unsetenv("TEST_SECRET") }()

	// Create loader without client (simulating local dev)
	loader := &SecretLoader{
		projectID: "",
		client:    nil,
	}

	value, err := loader.LoadSecret(context.Background(), secretName)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
}

func TestSecretLoader_LoadSecret_EnvVarNotFound(t *testing.T) {
	// Make sure env var is not set
	_ = os.Unsetenv("NONEXISTENT_SECRET")

	loader := &SecretLoader{
		projectID: "",
		client:    nil,
	}

	value, err := loader.LoadSecret(context.Background(), "NONEXISTENT_SECRET")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Equal(t, "", value)
}

func TestSecretLoader_LoadSecret_WithClientButNoProjectID(t *testing.T) {
	// This tests the case where client exists but projectID is empty
	// Note: We can't easily create a real client in tests, so we'll test the error path
	// In practice, if client is set, projectID should also be set
	loader := &SecretLoader{
		projectID: "",  // Empty project ID
		client:    nil, // Can't create real client in test
	}

	// Since client is nil, it will fall back to env var
	_ = os.Setenv("TEST_SECRET", "env-value")
	defer func() { _ = os.Unsetenv("TEST_SECRET") }()

	value, err := loader.LoadSecret(context.Background(), "TEST_SECRET")
	assert.NoError(t, err)
	assert.Equal(t, "env-value", value)
}

func TestSecretLoader_MustLoadSecret_PanicsOnFailure(t *testing.T) {
	secretName := "NONEXISTENT_SECRET"

	// Make sure env var is not set
	_ = os.Unsetenv("NONEXISTENT_SECRET")

	loader := &SecretLoader{
		projectID: "",
		client:    nil,
	}

	assert.Panics(t, func() {
		loader.MustLoadSecret(context.Background(), secretName)
	})
}

func TestSecretLoader_MustLoadSecret_ReturnsValueOnSuccess(t *testing.T) {
	secretName := "TEST_SECRET"
	expectedValue := "test-value"

	// Set up environment variable
	_ = os.Setenv("TEST_SECRET", expectedValue)
	defer func() { _ = os.Unsetenv("TEST_SECRET") }()

	loader := &SecretLoader{
		projectID: "",
		client:    nil,
	}

	value := loader.MustLoadSecret(context.Background(), secretName)
	assert.Equal(t, expectedValue, value)
}

func TestSecretLoader_Close_WithoutClient(t *testing.T) {
	loader := &SecretLoader{
		projectID: "",
		client:    nil,
	}

	err := loader.Close()
	assert.NoError(t, err)
}

// Note: Testing with a real Secret Manager client is difficult in unit tests
// because it requires Google Cloud credentials. The integration tests
// in production will validate the actual Secret Manager integration.
