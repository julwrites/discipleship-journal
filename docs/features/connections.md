# User Connections

## Overview
The User Connections feature allows users to establish direct relationships with other users on the platform. This serves as a foundational layer for social interactions, enabling users to find and connect with friends before inviting them to groups or sharing content.

## User Stories
- As a user, I want to search for other users by email so I can send them a connection request.
- As a user, I want to see a list of my pending connection requests (both sent and received).
- As a user, I want to accept or reject connection requests from others.
- As a user, I want to see a list of my accepted connections.
- As a user, I want to remove a connection if I no longer wish to be connected.

## API Endpoints

### Connection Management
- `POST /api/connections/request`: Send a connection request to another user by email.
    - Payload: `{ "email": "user@example.com" }`
- `GET /api/connections`: List all connections for the current user (requests and accepted).
- `PUT /api/connections/{id}`: Accept a pending connection request.
    - Returns: `{"success": true}`
- `DELETE /api/connections/{id}`: Reject a pending request or remove an existing connection.
    - Returns: `{"success": true}`

## Data Model

### Connection
- `ID`: UUID
- `RequesterID`: UUID (User ID)
- `ReceiverID`: UUID (User ID)
- `Status`: Enum (`pending`, `accepted`)
- `CreatedAt`: Timestamp
- `UpdatedAt`: Timestamp

## Frontend Implementation
- **ConnectionsPage**: The main interface for managing connections.
    - **Tabs**: "My Connections" (list accepted) and "Find People" (search and request).
    - **Pending Requests**: Displayed prominently if there are incoming requests.
- **Search**: Debounced search by email to find users to connect with.

## Interaction Flow
1.  **Request**: User A searches for User B by email and sends a request.
2.  **Pending**: User B sees a pending request from User A.
3.  **Accept**: User B accepts the request. The status becomes `accepted`.
4.  **Connect**: Both users now appear in each other's connection list.

## Future Enhancements
- Notification when a connection request is received.
- Block users.
- Profile visibility settings.
