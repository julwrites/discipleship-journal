package services

import (
	"context"
	"errors"
	"os"
	"testing"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
)

// MockSecretManagerClient for testing
type MockSecretManagerClient struct {
	accessSecretVersionFunc func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
	closeFunc               func() error
}

func (m *MockSecretManagerClient) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	if m.accessSecretVersionFunc != nil {
		return m.accessSecretVersionFunc(ctx, req, opts...)
	}
	return nil, errors.New("mock not implemented")
}

func (m *MockSecretManagerClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

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

func TestSecretLoader_LoadSecret_FromManager(t *testing.T) {
	mockClient := &MockSecretManagerClient{
		accessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
			return &secretmanagerpb.AccessSecretVersionResponse{
				Payload: &secretmanagerpb.SecretPayload{
					Data: []byte("secret-value"),
				},
			}, nil
		},
	}

	loader := &SecretLoader{
		projectID: "test-project",
		client:    mockClient,
	}

	value, err := loader.LoadSecret(context.Background(), "MY_SECRET")
	assert.NoError(t, err)
	assert.Equal(t, "secret-value", value)
}

func TestSecretLoader_LoadSecret_ManagerFailure(t *testing.T) {
	// Test fallback to env var when manager fails
	_ = os.Setenv("MY_SECRET", "env-value")
	defer func() { _ = os.Unsetenv("MY_SECRET") }()

	mockClient := &MockSecretManagerClient{
		accessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
			return nil, errors.New("access denied")
		},
	}

	loader := &SecretLoader{
		projectID: "test-project",
		client:    mockClient,
	}

	value, err := loader.LoadSecret(context.Background(), "MY_SECRET")
	assert.NoError(t, err)
	assert.Equal(t, "env-value", value)
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
	mockClient := &MockSecretManagerClient{}

	loader := &SecretLoader{
		projectID: "",  // Empty project ID
		client:    mockClient,
	}

	// Since projectID is empty, it will fall back to env var even if client is set
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

func TestSecretLoader_Close_WithClient(t *testing.T) {
    called := false
    mockClient := &MockSecretManagerClient{
        closeFunc: func() error {
            called = true
            return nil
        },
    }

    loader := &SecretLoader{
        projectID: "test",
        client:    mockClient,
    }

    err := loader.Close()
    assert.NoError(t, err)
    assert.True(t, called)
}
