package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup code if needed
	os.Exit(m.Run())
}

func TestSimple(t *testing.T) {
	if 1+1 != 2 {
		t.Error("Math is broken")
	}
}

// TestWithRedis is commented out because Docker is not available in this sandbox environment.
// To run this test, uncomment it and ensure Docker is running.
/*
func TestWithRedis(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()
}
*/
