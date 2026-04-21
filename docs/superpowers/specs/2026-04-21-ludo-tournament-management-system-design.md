# Ludo Tournament Management System — Design Specification

## Overview

A web application for managing knockout tournaments and round-robin leagues for the board game Ludo. Supports tournaments of ~100 players and long-running league systems with customizable scoring.

**Tech Stack:**
- Backend: Go + Gin + GORM + PostgreSQL
- Frontend: React + Vite + Tailwind CSS
- Deployment: Docker + Kubernetes

---

## Users & Roles

| Role | Permissions |
|------|-------------|
| **Admin** | Full system control — create/manage tournaments, leagues; ban/block users; disband any event |
| **Member** | Participate in leagues and tournaments; organize own events; invite guests |
| **Guest** | Accept invitations to register; join tournaments when invited; report game results |

**Invite Flow:**
1. Member sends invitation to an email address
2. System creates `UserInvite` record with unique code
3. Guest receives email, clicks link, creates account
4. System tracks which member invited which guest

---

## Core Features

### Knockout Tournaments

- **Bracket generation:** Random draw, seeded only when odd-player promotion applies
- **Odd-player handling:** If player count % 4 != 0, top-performing players from the 1-2 most recent prior tournaments are promoted to fill spots. Promotion is automatic; table/board placement is random to avoid seeding
- **Bracket editing:** Organizer or admin can manually adjust the bracket
- **Table assignment:** Players are assigned to tables (numbered boards) for each round. System tells each player which game and table number to play
- **Forfeit handling:** Missing player results in game played with one less participant
- **Live status:** Organizer can set tournament to `live`; results reported by winner, second-place, or organizer
- **Result reporting:** Any member who is organizing can report any game. Guests can report their own table's result

### Round-Robin Leagues

- **Format:** Every player faces every other player (or as close as possible with flexible scheduling)
- **Scheduling:** Organizer sets number of play dates and games per player per date
- **Table setup:** Configurable number of tables available per play date
- **Pairing algorithm:** Random with fairness tracking — system avoids repeat pairings when possible
- **Flexibility:** Organizer can manually swap players between tables; system warns if a swap causes players to face each other more than the average or significantly more than other players
- **Scoring:** Customizable placement points per league. Examples:
  - Positive: 1st=3, 2nd=2, 3rd=1, 4th=0
  - Negative: 1st=0, 2nd=-50, 3rd=-100, 4th=-150
- **Tiebreaker:** Default is most wins. Organizer can add a tiebreaker game
- **League standings:** Displayed with tiebreaker info, updated after each game

### Player Profiles & History

- **Profile:** Display name, join date, tournament and league history
- **Tournament history:** Tracks what round each player reached (1st round, semi-finals, etc.) for seeding purposes
- **League history:** Placement history, total points, games played, win count
- **Stats aggregation:** Both tournament and league stats viewable per player

### Invitations

- **User invites:** Sent to email → creates pending account → code-based registration
- **Tournament invites:** Organizer invites members/guests to participate
- **League invites:** Similar to tournament invites
- **Invite tracking:** Email, code, expiry, accept timestamp stored in `UserInvite`

---

## Data Model

### Entity Definitions

**User**
```
id, email, password_hash, role (admin/member/guest),
invited_by (user_id), last_active (timestamp),
created_at, modified_at, modified_by, deleted_at
```

**Player**
```
id, user_id, display_name,
tournament_history (json: [{tournament_id, round_reached, date}]),
league_stats (json: [{league_id, games_played, total_points, wins}]),
created_at, modified_at, modified_by, deleted_at
```

**Tournament**
```
id, name, type (knockout), organizer_id,
status (draft/live/completed),
settings (json: {tables_count, ...}),
created_at, modified_at, modified_by, deleted_at
```

**League**
```
id, name, organizer_id,
status (draft/live/completed),
settings (json: {scoring_rules, games_per_player, tables_count, ...}),
created_at, modified_at, modified_by, deleted_at
```

**KnockoutBracket**
```
id, tournament_id, rounds (json: [[{player_id, table}]]),
created_at, modified_at, modified_by, deleted_at
```

**LeagueSchedule**
```
id, league_id, play_dates (json: [{date, pairings}]),
created_at, modified_at, modified_by, deleted_at
```

**Match**
```
id, tournament_id/league_id, round, table_number,
status (pending/completed), placement_points (json),
created_at, modified_at, modified_by, deleted_at
```

**MatchAssignment**
```
id, match_id, player_id, seat_position,
result (1st/2nd/3rd/4th), reported_by (user_id),
created_at, modified_at, modified_by, deleted_at
```

**Invitation**
```
id, tournament_id/league_id, invitee_id,
status (pending/accepted/declined),
created_at, modified_at, modified_by, deleted_at
```

**UserInvite**
```
id, email, code (unique),
invited_by (user_id),
expires_at, accepted_at,
created_at, modified_at, modified_by, deleted_at
```

### Audit Fields

All entities include:
- `created_at` — creation timestamp
- `modified_at` — last update timestamp
- `modified_by` — user_id of last modifier
- `deleted_at` — soft delete (NULL = active)

---

## API Design

All paths prefixed with `/api` removed — clean REST paths.

### Authentication
```
POST /auth/register     — create account (requires valid invite code)
POST /auth/login        — email + password → JWT
POST /auth/logout       — invalidate session
GET  /auth/me           — current user profile
```

### Users (Admin)
```
GET    /users           — list all users
PATCH  /users/:id       — update user (role, status)
DELETE /users/:id       — soft delete user
```

### Invitations
```
POST   /invites         — send invite (email, type: user/tournament/league)
GET    /invites         — list invites (filtered by role/ownership)
POST   /invites/:code   — accept invite and register account
```

### Tournaments
```
POST   /tournaments              — create tournament
GET    /tournaments              — list tournaments (filter by status)
GET    /tournaments/:id          — tournament details + bracket/standings
PATCH  /tournaments/:id          — update settings/status
DELETE /tournaments/:id          — disband (admin/organizer)

GET    /tournaments/:id/matches  — all matches for tournament
POST   /tournaments/:id/matches  — report match result
GET    /tournaments/:id/pairings — current round pairings with table assignments
```

### Leagues
```
POST   /leagues           — create league
GET    /leagues           — list leagues
GET    /leagues/:id       — league details + standings + schedule
PATCH  /leagues/:id       — update settings/status
DELETE /leagues/:id       — disband (admin/organizer)

POST   /leagues/:id/play-dates           — add play date
GET    /leagues/:id/schedule             — full schedule with pairings
POST   /leagues/:id/pairings/generate    — generate round pairings (fairness-aware)
```

### Players
```
GET    /players/:id       — player profile + history
GET    /players/:id/stats — tournament + league statistics
```

### Admin
```
POST   /admin/seed        — promote odd-player from prior tournament
GET    /admin/audit-log  — modified_by tracking
```

### Error Format
```json
{
  "error": {
    "code": "TOURNAMENT_NOT_FOUND",
    "message": "Tournament with ID xyz not found"
  }
}
```

---

## Frontend Architecture (React)

```
src/
├── components/
│   ├── auth/           # LoginForm, RegisterForm, InviteAccept
│   ├── layout/         # Navbar, Sidebar, Footer, Container
│   ├── tournament/     # BracketView, MatchCard, TableAssignment, PairingDisplay
│   ├── league/         # StandingsTable, ScheduleView, PlayDateCard, ScorigGrid
│   ├── player/         # PlayerCard, PlayerProfile, PlayerStats
│   ├── invite/         # InviteList, InviteSender, InviteRow
│   └── ui/             # Button, Input, Modal, Card, Badge, Table (shared)
├── pages/
│   ├── Dashboard
│   ├── TournamentList, TournamentDetail
│   ├── LeagueList, LeagueDetail
│   ├── Profile
│   ├── AdminPanel
│   └── AuthPages/
├── hooks/
│   ├── useAuth, useTournament, useLeague, useToast
├── services/
│   └── api.ts          # fetch wrapper with auth
├── types/
│   └── index.ts        # TypeScript interfaces
├── App.tsx
└── main.tsx
```

**State Management:** React Query (TanStack Query) for server state. Local state for UI. Context only for auth.

**Routing:** React Router v6 with nested routes.

---

## Backend Architecture (Go — Hexagonal)

```
backend/
├── cmd/server/main.go

├── core/
│   ├── domain/
│   │   ├── models/          # User, Player, Tournament, League, Match, Invitation
│   │   ├── errors.go        # domain errors
│   │   └── events.go        # domain events
│   ├── ports/
│   │   ├── inbound/         # service interfaces (driving)
│   │   │   ├── tournament_service.go
│   │   │   ├── league_service.go
│   │   │   ├── auth_service.go
│   │   │   └── user_service.go
│   │   └── outbound/        # repository interfaces (driven)
│   │       ├── user_repository.go
│   │       ├── tournament_repository.go
│   │       ├── league_repository.go
│   │       └── match_repository.go
│   └── application/
│       ├── tournament.go     # bracket generation, seeding logic
│       ├── league.go        # round-robin scheduling, fairness tracking
│       ├── pairing.go       # table assignment
│       ├── auth.go
│       └── dto/             # request/response objects

├── adapters/
│   ├── primary/
│   │   └── http/
│   │       ├── middleware/
│   │       ├── auth_handler.go
│   │       ├── tournament_handler.go
│   │       ├── league_handler.go
│   │       ├── player_handler.go
│   │       └── router.go
│   └── secondary/
│       └── persistence/
│           ├── postgres.go
│           ├── gorm_user.go
│           ├── gorm_tournament.go
│           ├── gorm_league.go
│           └── gorm_match.go

├── Dockerfile
├── docker-compose.yml
└── go.mod
```

**Dependency Rule:** Dependencies point inward. Core knows nothing about adapters. GORM models stay in `adapters/secondary/persistence`.

---

## Deployment

**Docker:**
- Backend: Go binary, multi-stage build
- Frontend: Node build + nginx serve
- PostgreSQL 16-alpine

**docker-compose:**
- `postgres` service with persistent volume
- `backend` service on port 8080
- `frontend` service on port 80 (nginx)

**Kubernetes:**
- Deployments for backend and frontend
- StatefulSet for PostgreSQL with PVC
- ConfigMaps/Secrets for env vars
- Service + Ingress for external access

---

## Security

- Passwords hashed with bcrypt
- JWT in httpOnly cookies
- Role-based middleware (admin/member/guest checks)
- Input validation on all endpoints
- GORM parameterized queries (SQL injection prevention)

---

## Testing Strategy (TDD)

**Unit tests (core/application/):**
- Bracket generation correctness
- Odd-player promotion logic
- League pairing algorithm (fairness constraints)
- Table assignment logic
- Scoring calculations
- Tiebreaker detection

**Integration tests:**
- Full API flows
- Database operations via test containers

---

## MVP Scope vs Future Enhancements

**MVP:**
- Page refresh for live updates
- Core tournament and league management
- Invite-based registration
- Basic player profiles

**Future (real-time, etc.):**
- WebSocket updates for live standings
- Email notifications for invitations
- Advanced analytics dashboard
- Tournament seeding based on ELO/rating system
