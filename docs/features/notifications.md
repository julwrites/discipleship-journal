# Push Notifications

## Overview
The Push Notifications feature enables the application to send real-time alerts to users on their devices (Web/PWA). It leverages Firebase Cloud Messaging (FCM).

## User Stories
- As a user, I want to receive notifications when I am invited to a group.
- As a user, I want to receive notifications when someone shares a note with me.
- As a user, I want to be able to enable/disable notifications on my current device.

## API Endpoints

### Device Management
- `POST /api/notifications/register`: Register the current device's FCM token.
  - Payload: `{ "token": "string", "device_type": "web" }`

## Data Model

### UserDevice
- `ID`: UUID (Implicit/PK)
- `UserID`: UUID
- `FCMToken`: String (Unique per user)
- `DeviceType`: String (e.g., "web", "android", "ios")
- `LastUsedAt`: Timestamp

## Backend Implementation
- **NotificationService**: Handles logic for registering devices and sending multicast messages via FCM.
- **Cleanup**: Automatically removes invalid tokens if FCM reports failure (e.g., `UNREGISTERED`).
- **Triggers**: Notifications are triggered by other services (e.g., `GroupHandler` calls `NotificationService.SendNotification` upon member addition).

## Frontend Implementation
- **Service Worker (`web/src/sw.ts`)**: Handles background messages using `firebase/messaging/sw`.
- **Firebase Config (`web/src/lib/firebase.ts`)**: Initializes Messaging (if supported).
- **UI**: Components request permission (`Notification.requestPermission()`) and call the registration API upon success.

## Notification Types
Currently supported notification types (sent in `data` payload):
- `group_invite`: Sent when added to a group.
- `note_share`: (Planned) Sent when a note is shared.

## Future Enhancements
- User preferences (mute specific groups/types).
- In-app notification center (history).
- Scheduled reminders (e.g., "Time to read your Bible").
