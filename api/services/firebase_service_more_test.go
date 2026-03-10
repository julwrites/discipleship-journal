package services

import (
	"context"
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestNewFirebaseService_WithSAKey(t *testing.T) {
	service, err := NewFirebaseService(context.Background(), "invalid json", "discipleship-journal-pwa")
	assert.Error(t, err)
	assert.Nil(t, service)
}

func TestNewFirebaseService_WithValidishSAKey(t *testing.T) {
	validJson := `{"type": "service_account", "project_id": "test"}`
	service, err := NewFirebaseService(context.Background(), validJson, "discipleship-journal-pwa")

	if err == nil {
		assert.NotNil(t, service)
	} else {
		assert.Nil(t, service)
	}
}

func TestNewFirebaseService_EnvConfig(t *testing.T) {
	os.Setenv("FIREBASE_CONFIG", `{"projectId": "test-project"}`)
	service, err := NewFirebaseService(context.Background(), "", "test-project")

	if err == nil {
		assert.NotNil(t, service)
	} else {
		assert.Nil(t, service)
	}
	os.Unsetenv("FIREBASE_CONFIG")
}
