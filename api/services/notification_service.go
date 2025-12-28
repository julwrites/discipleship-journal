package services

import (
	"context"
	"log/slog"

	"discipleship_journal_api/database"

	"firebase.google.com/go/v4/messaging"
)

type NotificationService interface {
	RegisterDevice(ctx context.Context, userID, token, deviceType string) error
	SendNotification(ctx context.Context, userID, title, body string, data map[string]string) error
}

type MessagingClientInterface interface {
	SendEachForMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)
}

type notificationService struct {
	db              database.DBInterface
	messagingClient MessagingClientInterface
}

func NewNotificationService(db database.DBInterface, msgClient MessagingClientInterface) NotificationService {
	return &notificationService{
		db:              db,
		messagingClient: msgClient,
	}
}

func (s *notificationService) RegisterDevice(ctx context.Context, userID, token, deviceType string) error {
	query := `
		INSERT INTO user_devices (user_id, fcm_token, device_type, last_used_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, fcm_token)
		DO UPDATE SET last_used_at = NOW(), device_type = EXCLUDED.device_type
	`
	_, err := s.db.Exec(ctx, query, userID, token, deviceType)
	return err
}

func (s *notificationService) SendNotification(ctx context.Context, userID, title, body string, data map[string]string) error {
	// 1. Get tokens for user
	tokens, err := s.getUserTokens(ctx, userID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil // No devices to notify
	}

	// 2. Create message
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	// 3. Send message
	br, err := s.messagingClient.SendEachForMulticast(ctx, message)
	if err != nil {
		return err
	}

	// 4. Handle failed tokens (cleanup)
	if br.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range br.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens = append(failedTokens, tokens[idx])
				slog.Warn("Failed to send notification", "token", tokens[idx], "error", resp.Error)
			}
		}
		if len(failedTokens) > 0 {
			if err := s.removeInvalidTokens(ctx, userID, failedTokens); err != nil {
				slog.Error("Failed to remove invalid tokens", "error", err)
			}
		}
	}

	return nil
}

func (s *notificationService) getUserTokens(ctx context.Context, userID string) ([]string, error) {
	rows, err := s.db.Query(ctx, "SELECT fcm_token FROM user_devices WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, rows.Err()
}

func (s *notificationService) removeInvalidTokens(ctx context.Context, userID string, tokens []string) error {
	// Bulk delete could be better, but loop is simple for now or use ANY
	query := "DELETE FROM user_devices WHERE user_id = $1 AND fcm_token = ANY($2)"
	_, err := s.db.Exec(ctx, query, userID, tokens)
	return err
}
