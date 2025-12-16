package services

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FirebaseService struct {
	App             *firebase.App
	AuthClient      *auth.Client
	MessagingClient *messaging.Client
}

func NewFirebaseService(ctx context.Context, saKey string, projectID string) (*FirebaseService, error) {
	var opts []option.ClientOption
	if saKey != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(saKey)))
	}

	config := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	// Messaging client might fail if the service account doesn't have permissions or if running in a constrained env.
	// We allow it to be nil if it fails, so the app can start (e.g. in development/tests without full creds).
	// In production, this should likely be a fatal error if notifications are critical.
	messagingClient, _ := app.Messaging(ctx)

	return &FirebaseService{
		App:             app,
		AuthClient:      authClient,
		MessagingClient: messagingClient,
	}, nil
}
