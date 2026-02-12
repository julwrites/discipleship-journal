package services

import (
	"context"
	"testing"
)

func TestNewFirebaseService(t *testing.T) {
	// This should fail because FIREBASE_CONFIG is not set or valid
	service, err := NewFirebaseService(context.Background(), "", "discipleship-journal-pwa")
	if err == nil {
		// If it somehow succeeds
		if service == nil {
			t.Fatal("NewFirebaseService returned nil but no error")
		}
	} else {
		// Expected error
		if service != nil {
			t.Error("NewFirebaseService returned non-nil struct on error")
		}
	}
}
