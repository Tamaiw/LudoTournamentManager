# API Reference

## Base URL

```
Development: http://localhost:8080
Production:  https://api.ludotournament.example.com
```

## Authentication

Most endpoints require JWT authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

### Endpoints

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | /auth/register | Create account with invite code | No |
| POST | /auth/login | Login with email/password | No |
| POST | /auth/logout | Invalidate session | Yes |
| GET | /auth/me | Get current user profile | Yes |

### POST /auth/register

Create a new account. Requires a valid invite code.

**Request:**
```json
{
  "email": "player@example.com",
  "password": "securePassword123",
  "invite_code": "ABC123XYZ"
}
```

**Response (201):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "player@example.com",
    "role": "member"
  }
}
```

**Errors:**
- `400` - Invalid input or missing fields
- `401` - Invalid invite code
- `409` - Email already registered

### POST /auth/login

Login with email and password.

**Request:**
```json
{
  "email": "player@example.com",
  "password": "securePassword123"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "player@example.com",
    "role": "member"
  }
}
```

**Errors:**
- `401` - Invalid credentials

### POST /auth/logout

Invalidate the current session.

**Response (204):** No content

### GET /auth/me

Get the currently authenticated user's profile.

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "player@example.com",
  "role": "member",
  "last_active": "2024-04-22T10:30:00Z",
  "created_at": "2024-04-01T08:00:00Z"
}
```

---

## Tournaments

### Endpoints

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | /tournaments | Create tournament | Yes (Member+) |
| GET | /tournaments | List tournaments | Yes |
| GET | /tournaments/:id | Get tournament details | Yes |
| PATCH | /tournaments/:id | Update tournament | Yes (Organizer+) |
| DELETE | /tournaments/:id | Delete tournament | Yes (Admin) |
| GET | /tournaments/:id/matches | List all matches | Yes |
| POST | /tournaments/:id/matches | Report match result | Yes |
| GET | /tournaments/:id/pairings | Get current round pairings | Yes |

### POST /tournaments

Create a new knockout tournament.

**Request:**
```json
{
  "name": "Spring Championship 2024",
  "settings": {
    "tables_count": 20,
    "advancement": [
      {
        "round": "round_1",
        "games": 20,
        "advancement_per_game": [
          {"game_ids": [1, 4, 7, 10, 13, 16, 19], "placements": [1, 2]},
          {"game_ids": [2, 5, 8, 11, 14, 17, 20], "placements": [1, 2, 3]}
        ]
      }
    ],
    "default_reporter": "lowest_advancing"
  }
}
```

**Response (201):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "name": "Spring Championship 2024",
  "type": "knockout",
  "organizer_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "draft",
  "settings": { ... },
  "created_at": "2024-04-22T10:00:00Z"
}
```

### GET /tournaments

List all tournaments. Supports filtering by status.

**Query parameters:**
- `status=draft|live|completed` - Filter by status

**Response (200):**
```json
{
  "tournaments": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "name": "Spring Championship 2024",
      "type": "knockout",
      "organizer_id": "550e8400-e29b-41d4-a716-446655440000",
      "status": "live",
      "created_at": "2024-04-22T10:00:00Z"
    }
  ]
}
```

### GET /tournaments/:id

Get detailed tournament information including bracket.

**Response (200):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "name": "Spring Championship 2024",
  "type": "knockout",
  "organizer_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "live",
  "settings": {
    "tables_count": 20,
    "advancement": [ ... ]
  },
  "created_at": "2024-04-22T10:00:00Z"
}
```

### PATCH /tournaments/:id

Update tournament settings or status.

**Request:**
```json
{
  "status": "live"
}
```

**Response (200):** Updated tournament object

### DELETE /tournaments/:id

Delete (soft delete) a tournament. Admin only.

**Response (204):** No content

### GET /tournaments/:id/matches

List all matches in a tournament.

**Query parameters:**
- `round=N` - Filter by round number

**Response (200):**
```json
{
  "matches": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "round": 1,
      "table_number": 1,
      "status": "completed",
      "placement_points": [3, 2, 1, 0],
      "completed_at": "2024-04-22T11:30:00Z"
    }
  ]
}
```

### POST /tournaments/:id/matches

Report a match result.

**Request:**
```json
{
  "match_id": "770e8400-e29b-41d4-a716-446655440002",
  "results": [
    {"player_id": "880e8400-e29b-41d4-a716-446655440003", "seat_color": "yellow", "placement": 1},
    {"player_id": "880e8400-e29b-41d4-a716-446655440004", "seat_color": "green", "placement": 2},
    {"player_id": "880e8400-e29b-41d4-a716-446655440005", "seat_color": "blue", "placement": 3},
    {"player_id": "880e8400-e29b-41d4-a716-446655440006", "seat_color": "red", "placement": 4}
  ]
}
```

**Response (200):**
```json
{
  "match_id": "770e8400-e29b-41d4-a716-446655440002",
  "status": "completed",
  "advancing_players": [
    {"player_id": "880e8400-e29b-41d4-a716-446655440003", "placement": 1},
    {"player_id": "880e8400-e29b-41d4-a716-446655440004", "placement": 2}
  ]
}
```

**Errors:**
- `400` - Invalid results (wrong number of players, invalid placements)
- `409` - Match already completed or downstream game played (edit lock)

### GET /tournaments/:id/pairings

Get current round pairings with table assignments.

**Response (200):**
```json
{
  "round": 2,
  "pairings": [
    {
      "game_id": "770e8400-e29b-41d4-a716-446655440010",
      "round": 2,
      "table_number": 1,
      "player_ids": [
        "880e8400-e29b-41d4-a716-446655440003",
        "880e8400-e29b-41d4-a716-446655440007",
        "880e8400-e29b-41d4-a716-446655440004",
        "880e8400-e29b-41d4-a716-446655440008"
      ],
      "seat_colors": ["yellow", "green", "blue", "red"],
      "status": "pending"
    }
  ]
}
```

---

## Leagues

### Endpoints

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | /leagues | Create league | Yes (Member+) |
| GET | /leagues | List leagues | Yes |
| GET | /leagues/:id | Get league details | Yes |
| PATCH | /leagues/:id | Update league | Yes (Organizer+) |
| DELETE | /leagues/:id | Delete league | Yes (Admin) |
| POST | /leagues/:id/play-dates | Add play date | Yes (Organizer+) |
| GET | /leagues/:id/schedule | Get full schedule | Yes |
| POST | /leagues/:id/pairings/generate | Generate pairings | Yes |
| GET | /leagues/:id/standings | Get standings | Yes |

### POST /leagues

Create a new round-robin league.

**Request:**
```json
{
  "name": "Spring League 2024",
  "settings": {
    "scoring_rules": [
      {"placement": 1, "points": 3},
      {"placement": 2, "points": 2},
      {"placement": 3, "points": 1},
      {"placement": 4, "points": 0}
    ],
    "games_per_player": 3,
    "tables_count": 10
  }
}
```

**Response (201):**
```json
{
  "id": "990e8400-e29b-41d4-a716-446655440010",
  "name": "Spring League 2024",
  "organizer_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "draft",
  "settings": { ... },
  "created_at": "2024-04-22T10:00:00Z"
}
```

### GET /leagues/:id/standings

Get current league standings.

**Response (200):**
```json
{
  "standings": [
    {
      "player_id": "880e8400-e29b-41d4-a716-446655440003",
      "display_name": "Alice",
      "games_played": 6,
      "total_points": 15,
      "wins": 3,
      "rank": 1
    },
    {
      "player_id": "880e8400-e29b-41d4-a716-446655440004",
      "display_name": "Bob",
      "games_played": 6,
      "total_points": 12,
      "wins": 2,
      "rank": 2
    }
  ]
}
```

### POST /leagues/:id/pairings/generate

Generate fair pairings for the next play date.

**Request:**
```json
{
  "play_date": "2024-04-25"
}
```

**Response (200):**
```json
{
  "play_date": "2024-04-25",
  "pairings": [
    {
      "match_id": "aa0e8400-e29b-41d4-a716-446655440020",
      "table_number": 1,
      "player_ids": [
        "880e8400-e29b-41d4-a716-446655440003",
        "880e8400-e29b-41d4-a716-446655440004",
        "880e8400-e29b-41d4-a716-446655440005",
        "880e8400-e29b-41d4-a716-446655440006"
      ]
    }
  ]
}
```

---

## Users (Admin Only)

### Endpoints

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| GET | /users | List all users | Yes (Admin) |
| PATCH | /users/:id | Update user role/status | Yes (Admin) |
| DELETE | /users/:id | Soft delete user | Yes (Admin) |

### GET /users

List all users.

**Response (200):**
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "admin@example.com",
      "role": "admin",
      "last_active": "2024-04-22T10:30:00Z",
      "created_at": "2024-04-01T08:00:00Z"
    }
  ]
}
```

### PATCH /users/:id

Update user (role or status).

**Request:**
```json
{
  "role": "member"
}
```

**Response (200):** Updated user object

---

## Error Format

All errors follow a consistent format:

```json
{
  "error": {
    "code": "TOURNAMENT_NOT_FOUND",
    "message": "Tournament with ID xyz not found"
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| VALIDATION_ERROR | 400 | Invalid input data |
| UNAUTHORIZED | 401 | Missing or invalid auth |
| FORBIDDEN | 403 | Insufficient permissions |
| NOT_FOUND | 404 | Resource not found |
| CONFLICT | 409 | Resource conflict (e.g., already completed) |
| INTERNAL_ERROR | 500 | Server error |

---

## Rate Limits

Development: No rate limits
Production: 100 requests/minute per user

---

## Versioning

API versioning via URL path prefix when needed:
```
/api/v1/auth/register
```

Current API uses no prefix (v0 implicit).