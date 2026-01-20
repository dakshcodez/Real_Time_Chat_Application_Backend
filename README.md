# Real-Time Chat Application Backend

A production-ready backend implementation for a real-time one-to-one chat application built with Go. This backend provides JWT-based authentication, WebSocket-powered real-time messaging, message persistence, and comprehensive rate limiting.

## Tech Stack

- **Language**: Go 1.24.5
- **HTTP Server**: net/http (standard library)
- **WebSockets**: gorilla/websocket
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Rate Limiting**: In-memory token bucket implementation
- **Testing Tools**: Postman, wscat

## Features

### Authentication & Authorization

- User registration with email and username validation
- Secure password hashing using bcrypt
- JWT-based authentication with 24-hour token expiration
- Authentication middleware using request context
- Protected endpoints requiring Bearer token authentication
- User profile management with ownership enforcement
- Users can only update their own profile details

### Real-Time Messaging

- Authenticated WebSocket connections via JWT token
- One-to-one private messaging between users
- Multi-device support per user (multiple concurrent connections)
- Automatic ping/pong heartbeats for connection health
- Safe concurrent read/write handling with goroutines
- Real-time message delivery to online recipients

### Message Persistence

- All messages stored in PostgreSQL with UUID primary keys
- Offline message support (messages delivered when user comes online)
- Accurate timestamps with timezone support
- Soft delete support (messages marked as deleted, not removed)

### Chat History API

- Fetch chat history between two authenticated users
- Authorization enforced implicitly (users can only access their own conversations)
- Cursor-based pagination using timestamp cursors
- Configurable page size (default: 20, max: 100)
- Deleted messages automatically excluded from history

### Message Edit & Delete

- Edit messages (sender only, enforced at database level)
- Delete messages ("delete for everyone" via soft delete)
- Real-time WebSocket broadcast of edit/delete events
- Ownership verification prevents unauthorized modifications
- Edit timestamps tracked for audit purposes

### Rate Limiting

- User-based rate limiting for REST APIs (60 requests per minute)
- User-based message rate limiting for WebSockets (10 messages per second)
- Protection against spam and abuse
- In-memory token bucket algorithm implementation
- Rate limit applied after authentication (user-aware)

## Project Structure

```
real_time_chat_application_backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── auth/                    # Authentication logic
│   │   ├── handler.go          # HTTP handlers for register/login
│   │   ├── jwt.go              # JWT token generation and parsing
│   │   └── service.go          # User registration and login logic
│   │
│   ├── config/                  # Configuration management
│   │   └── config.go           # Environment variable loading
│   │
│   ├── db/                      # Database connection and migrations
│   │   └── postgres.go         # GORM connection and auto-migration
│   │
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go             # JWT authentication middleware
│   │   └── rate_limit.go       # Rate limiting middleware
│   │
│   ├── models/                  # GORM data models
│   │   ├── user.go             # User model
│   │   └── message.go          # Message model
│   │
│   ├── ratelimit/               # Rate limiting implementation
│   │   └── limiter.go          # Token bucket rate limiter
│   │
│   ├── server/                  # HTTP request handlers
│   │   ├── http.go             # Route registration
│   │   ├── user_handler.go     # User profile endpoints
│   │   ├── chat_handler.go     # Chat history endpoint
│   │   └── message_handler.go  # Message edit/delete endpoints
│   │
│   └── websocket/                # WebSocket implementation
│       ├── handler.go           # WebSocket connection handler
│       ├── hub.go               # Connection hub and message routing
│       ├── client.go            # Client connection management
│       ├── protocol.go          # Message protocol definitions
│       ├── service.go           # Message persistence service
│       └── message_query.go    # Chat history query logic
│
├── go.mod                        # Go module dependencies
├── go.sum                        # Go module checksums
└── README.md                     # Project documentation
```

## API Documentation

### Authentication

#### Register User

```http
POST /auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response**: `201 Created` (empty body)

**Errors**:
- `400 Bad Request`: Invalid input or user already exists

#### Login

```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response**: `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors**:
- `401 Unauthorized`: Invalid credentials

### User Management

#### Get Current User Profile

```http
GET /users/me
Authorization: Bearer <JWT_TOKEN>
```

**Response**: `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "johndoe",
  "email": "john@example.com",
  "created": "2024-01-15T10:30:00Z"
}
```

#### Update Current User Profile

```http
PUT /users/me/update
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "username": "newusername",
  "email": "newemail@example.com"
}
```

**Note**: Both fields are optional. Only provided fields will be updated.

**Response**: `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "newusername",
  "email": "newemail@example.com",
  "updated": "2024-01-15T11:00:00Z"
}
```

### Chat & Messages

#### Get Chat History

```http
GET /chats/{userId}?limit=20&before=1705312800
Authorization: Bearer <JWT_TOKEN>
```

**Query Parameters**:
- `limit` (optional): Number of messages to return (default: 20, max: 100)
- `before` (optional): Unix timestamp cursor for pagination

**Response**: `200 OK`
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "from": "550e8400-e29b-41d4-a716-446655440000",
    "to": "770e8400-e29b-41d4-a716-446655440001",
    "content": "Hello, how are you?",
    "timestamp": "2024-01-15T10:30:00Z",
    "edited_at": null
  }
]
```

#### Edit Message

```http
PATCH /messages/{messageId}
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "content": "Updated message content"
}
```

**Response**: `200 OK`
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "sender_id": "550e8400-e29b-41d4-a716-446655440000",
  "receiver_id": "770e8400-e29b-41d4-a716-446655440001",
  "content": "Updated message content",
  "is_deleted": false,
  "edited_at": "2024-01-15T11:00:00Z",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Errors**:
- `403 Forbidden`: User is not the sender of the message

#### Delete Message

```http
DELETE /messages/{messageId}/delete
Authorization: Bearer <JWT_TOKEN>
```

**Response**: `200 OK`
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "sender_id": "550e8400-e29b-41d4-a716-446655440000",
  "receiver_id": "770e8400-e29b-41d4-a716-446655440001",
  "content": "Original message",
  "is_deleted": true,
  "edited_at": null,
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Errors**:
- `403 Forbidden`: User is not the sender of the message

### WebSocket

#### Connect to WebSocket

```
wscat -c "ws://localhost:8080/ws?token=<JWT_TOKEN>"
```

The WebSocket connection requires a valid JWT token as a query parameter. The connection will be rejected if the token is missing or invalid.

#### Sending a Message

```json
{
  "type": "direct_message",
  "to": "<RECEIVER_ID>",
  "content": "Hello, this is a test message"
}
```

#### Receiving a Message

```json
{
  "type": "direct_message",
  "from": "<SENDER_ID>",
  "content": "Hello, this is a test message",
  "timestamp": 1705312800
}
```

#### Receiving Edit Event

```json
{
  "type": "message_edited",
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "content": "Updated message content"
}
```

#### Receiving Delete Event

```json
{
  "type": "message_deleted",
  "id": "660e8400-e29b-41d4-a716-446655440000"
}
```

## Setup Instructions

### Prerequisites

- Go 1.24.5 or later
- PostgreSQL 12 or later
- Environment variables configured

### Environment Variables

Create a `.env` file in the project root:

```env
DATABASE_URL=postgres://username:password@localhost:5432/chatdb?sslmode=disable
JWT_SECRET=your-secret-key-here-minimum-32-characters
```

**Important**: Use a strong, random secret key for `JWT_SECRET` in production (minimum 32 characters).

### Database Setup

1. Create a PostgreSQL database:
```bash
createdb chatdb
```

2. The application will automatically run migrations on startup using GORM's AutoMigrate feature.

### Running the Server

1. Install dependencies:
```bash
go mod download
```

2. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`.

## Testing Instructions

### Testing REST APIs with Postman

1. **Register a new user**:
   - Method: `POST`
   - URL: `http://localhost:8080/auth/register`
   - Body (JSON):
     ```json
     {
       "username": "testuser",
       "email": "test@example.com",
       "password": "password123"
     }
     ```

2. **Login**:
   - Method: `POST`
   - URL: `http://localhost:8080/auth/login`
   - Body (JSON):
     ```json
     {
       "email": "test@example.com",
       "password": "password123"
     }
     ```
   - Copy the `token` from the response

3. **Access protected endpoints**:
   - Add header: `Authorization: Bearer <token>`
   - Test endpoints:
     - `GET /users/me`
     - `PUT /users/me/update`
     - `GET /chats/{userId}`

### Testing WebSockets with wscat

1. Install wscat (if not already installed):
```bash
npm install -g wscat
```

2. Connect to WebSocket:
```bash
wscat -c "ws://localhost:8080/ws?token=<JWT_TOKEN>"
```

3. Send a message:
```json
{"type":"direct_message","to":"<RECEIVER_USER_ID>","content":"Hello!"}
```

4. You should receive the message echoed back if the receiver is connected, or it will be delivered when they connect.

## Security & Design Notes

### Password Updates Not Included in Profile Updates

Password changes are intentionally excluded from the profile update endpoint (`PUT /users/me/update`) for security reasons:

1. **Separation of Concerns**: Password changes require additional validation (current password verification, strength requirements, confirmation) that differs from simple profile field updates.

2. **Security Best Practices**: Password changes should be a separate, more secure flow that may include:
   - Current password verification
   - Password strength validation
   - Rate limiting specific to password changes
   - Optional email notifications
   - Session invalidation considerations

3. **Audit Trail**: Password changes should be logged separately from profile updates for security auditing purposes.

### User-Based Rate Limiting

Rate limiting is implemented per-user rather than per-IP address for several reasons:

1. **Authenticated Context**: All protected endpoints require JWT authentication, providing a reliable user identity. IP-based limiting would be ineffective since users can change IPs or share networks.

2. **Fair Usage**: User-based limiting ensures fair resource allocation per account, preventing a single user from consuming excessive resources regardless of their network location.

3. **Abuse Prevention**: Malicious users cannot bypass rate limits by changing IPs. Each authenticated user has their own rate limit bucket.

4. **Multi-Device Support**: Users can connect from multiple devices/IPs simultaneously, and user-based limiting correctly aggregates their usage across all connections.

### Separation of REST and WebSocket Responsibilities

The architecture separates REST APIs and WebSocket connections:

- **REST APIs**: Handle synchronous operations like authentication, profile management, chat history retrieval, and message editing/deletion. These operations benefit from HTTP's request/response model and can leverage standard middleware.

- **WebSockets**: Handle real-time bidirectional communication for message delivery. The WebSocket hub manages connection lifecycle and message routing, while REST endpoints handle persistence and business logic.

This separation allows:
- Independent scaling of REST and WebSocket components
- Different rate limiting strategies (60 req/min for REST, 10 msg/sec for WebSocket)
- Clear separation of concerns between synchronous and asynchronous operations

## Future Improvements

- **Password Change Flow**: Implement a dedicated endpoint with current password verification and email notifications
- **Online/Offline Presence**: Track user online status and broadcast presence changes
- **Read Receipts**: Implement message read status tracking and delivery confirmations
- **Group Chats**: Extend messaging to support multi-user group conversations
- **Redis-Backed Distributed Rate Limiting**: Replace in-memory rate limiting with Redis for multi-instance deployments
- **Graceful Shutdown**: Implement graceful shutdown handling for WebSocket connections and HTTP server
- **Observability**: Add structured logging, metrics collection, and distributed tracing
- **Message Search**: Implement full-text search across message history
- **File Attachments**: Support for image, document, and media file sharing
- **Push Notifications**: Mobile push notifications for offline message delivery

## License

This project is part of a backend development task and is provided as-is for demonstration purposes.
