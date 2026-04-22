# Architecture Guide

## Overview

The Ludo Tournament Management System follows **hexagonal architecture** (also known as ports and adapters or clean architecture). The core business logic is isolated from external dependencies, making the system testable, maintainable, and flexible.

## Core Principle: Dependencies Point Inward

```
Infrastructure → Adapters → Ports → Application → Domain
```

The domain layer has zero dependencies on external systems. All I/O happens through interfaces defined in the ports layer.

## Layer Breakdown

### 1. Domain Layer (`core/domain/`)

Contains pure business logic with no external dependencies.

**Models** (`core/domain/models/`):
- `User` - authentication, roles, invite tracking
- `Player` - display name, tournament history, league stats
- `Tournament` - knockout tournament with configurable advancement
- `League` - round-robin league with scoring rules
- `Match` - individual game with seat assignments
- `MatchAssignment` - player seat and result tracking
- `Invitation` - tournament/league invite tracking
- `KnockoutBracket` - bracket structure with round tracking

**Errors** (`core/domain/errors.go`):
```go
var (
    ErrNotFound          = errors.New("entity not found")
    ErrInvalidInput      = errors.New("invalid input")
    ErrUnauthorized      = errors.New("unauthorized")
    ErrForbidden         = errors.New("forbidden")
    ErrTournamentActive = errors.New("tournament is active and cannot be modified")
    ErrGameAlreadyPlayed = errors.New("game has already been played")
    ErrInvalidAdvancement = errors.New("advancement configuration is invalid")
    ErrNoRematch         = errors.New("players from same source game cannot be seated together")
)
```

**Events** (`core/domain/events.go`):
Domain events for decoupled side effects (future: notifications, audit logging).

### 2. Ports Layer (`core/ports/`)

Interfaces that define how the application interacts with the outside world.

**Inbound Ports** (`core/ports/inbound/`):
Driving adapters call these interfaces. Define what the application can do:
- `TournamentService` - create/get/update/delete tournaments, generate brackets, report matches
- `LeagueService` - create/get/update/delete leagues, generate pairings, report results
- `AuthService` - register, login, logout, get current user
- `UserService` - list/update/delete users, send invites

**Outbound Ports** (`core/ports/outbound/`):
Driven adapters implement these interfaces. Define what the application needs:
- `UserRepository` - CRUD for users
- `PlayerRepository` - CRUD for players
- `TournamentRepository` - CRUD for tournaments
- `LeagueRepository` - CRUD for leagues
- `MatchRepository` - CRUD for matches, query by round/tournament
- `InvitationRepository` - CRUD for invitations

### 3. Application Layer (`core/application/`)

Use cases and business logic orchestration. Depends only on inbound/outbound ports.

**Tournament Logic:**
- `GenerateBracket` - random draw, handling odd-player promotion
- `ValidateAdvancementConfig` - ensures total advancing players = next round spots
- `AssignSeatsToNextRound` - distributes players ensuring no rematches
- `AssignYellowSeat` - first-place finisher who finished earliest gets yellow

**League Logic:**
- `GenerateFairPairings` - random pairing minimizing repeat matches
- `CalculateLeagueStandings` - applies scoring rules, computes ranks
- `DetectTiebreaker` - identifies tied players

### 4. Adapters Layer (`adapters/`)

**Primary Adapters** (`adapters/primary/http/`):
- HTTP handlers (Gin) - receive requests, call inbound ports
- Middleware - authentication, role authorization
- Router - endpoint definitions

**Secondary Adapters** (`adapters/secondary/persistence/`):
- GORM repositories - implement outbound ports
- Database connection management

## Data Model

### Entity Relationships

```
User (1) ──→ (1) Player
                 │
                 ├──→ (N) MatchAssignment ──→ (1) Match ←── (N) MatchAssignment
                 │                                    │
                 │                                    └──→ (1) Tournament (optional)
                 │                                    │
                 │                                    └──→ (1) League (optional)
                 │
                 └──→ (N) Invitation ──→ Tournament|Leeague

Tournament (1) ──→ (N) Match
League (1) ──→ (N) Match
```

### Tournament Advancement Config

```json
{
  "advancement": [
    {
      "round": "round_1",
      "games": 20,
      "advancement_per_game": [
        {"game_ids": [1, 4, 7, ...], "placements": [1, 2]},
        {"game_ids": [2, 5, 8, ...], "placements": [1, 2, 3]}
      ]
    }
  ]
}
```

The `placements` array defines which finishing positions advance. Example above: 10 games with 2-spot advancement, 10 games with 3-spot advancement = 50 players advancing to round 2.

### League Scoring Rules

```json
{
  "scoring_rules": [
    {"placement": 1, "points": 3},
    {"placement": 2, "points": 2},
    {"placement": 3, "points": 1},
    {"placement": 4, "points": 0}
  ]
}
```

Scoring is configurable per league. Can be positive (1st gets most points) or negative (1st gets 0, 4th gets -150).

## Key Design Decisions

### 1. No Immediate Rematches

Players advancing from the same game in round N cannot be seated together in round N+1. The `AssignSeatsToNextRound` function distributes players across tables using round-robin assignment.

**Example:**
```
Round 1: Game A (players 1,2,3,4), Game B (players 5,6,7,8)
Round 2: Table 1 gets {1, 5, 2, 6} - no two from same source
         Table 2 gets {3, 7, 4, 8} - no two from same source
```

### 2. Yellow Seat First-Come, First-Served

When 1st-place finishers from round N are seated in round N+1, the yellow seat goes to the player who finished 1st earliest. This rewards speed.

### 3. Edit Lock Rule

A game cannot be edited if any downstream game has been played. This prevents invalidating advancement chains.

```
If Game 5 in Round 2 is played:
  → Game 1, 2, 3, 4 in Round 1 become locked
  → Their advancement assignments are finalized
```

### 4. Partial Round Completion

A game in round N+1 can start as soon as all 4 spots are filled, even if other games in round N are still pending. This keeps the tournament moving.

### 5. Soft Deletes

All entities include `deleted_at` timestamp. Hard deletes are never performed, enabling audit trails and data recovery.

## Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Entry point, wires up dependencies
├── core/
│   ├── domain/
│   │   ├── models/           # All domain entities
│   │   │   ├── user.go
│   │   │   ├── player.go
│   │   │   ├── tournament.go
│   │   │   ├── league.go
│   │   │   ├── match.go
│   │   │   ├── knockout_bracket.go
│   │   │   └── invitation.go
│   │   ├── errors.go         # Domain errors
│   │   └── events.go         # Domain events
│   ├── ports/
│   │   ├── inbound/           # Service interfaces (driving)
│   │   │   ├── tournament_service.go
│   │   │   ├── league_service.go
│   │   │   ├── auth_service.go
│   │   │   └── user_service.go
│   │   └── outbound/          # Repository interfaces (driven)
│   │       ├── user_repository.go
│   │       ├── player_repository.go
│   │       ├── tournament_repository.go
│   │       ├── league_repository.go
│   │       ├── match_repository.go
│   │       └── invitation_repository.go
│   └── application/           # Business logic (orchestration)
│       ├── tournament.go
│       ├── tournament_test.go
│       ├── league.go
│       ├── league_test.go
│       └── dto/               # Request/response objects
├── adapters/
│   ├── primary/
│   │   └── http/              # HTTP handlers, middleware, router
│   │       ├── middleware/
│   │       │   ├── auth.go
│   │       │   └── role.go
│   │       ├── auth_handler.go
│   │       ├── tournament_handler.go
│   │       ├── league_handler.go
│   │       ├── user_handler.go
│   │       └── router.go
│   └── secondary/
│       └── persistence/      # GORM implementations
│           ├── postgres.go   # DB connection
│           ├── gorm_user.go
│           ├── gorm_player.go
│           ├── gorm_tournament.go
│           ├── gorm_league.go
│           ├── gorm_match.go
│           ├── gorm_match_assignment.go
│           └── gorm_invitation.go
└── Dockerfile

frontend/
├── src/
│   ├── components/
│   │   ├── tournament/        # BracketView, MatchCard, TableAssignment
│   │   ├── league/            # StandingsTable, ScheduleView
│   │   ├── ui/                # Shared Button, Input, Modal, etc.
│   │   └── layout/            # Navbar, Sidebar, Container
│   ├── pages/
│   │   ├── TournamentList.tsx
│   │   ├── TournamentDetail.tsx
│   │   ├── LeagueList.tsx
│   │   ├── LeagueDetail.tsx
│   │   ├── Dashboard.tsx
│   │   ├── Profile.tsx
│   │   └── AuthPages/
│   ├── hooks/                 # useAuth, useTournament, useLeague, etc.
│   ├── services/
│   │   └── api.ts            # Fetch wrapper with auth headers
│   ├── types/
│   │   └── index.ts          # TypeScript interfaces
│   ├── App.tsx
│   └── main.tsx
├── index.html
├── package.json
├── vite.config.ts
└── tailwind.config.js

docs/
├── ARCHITECTURE.md    # This file
├── API.md             # REST API reference
└── DEVELOPMENT.md     # Development setup guide

k8s/
├── backend-deployment.yaml
├── frontend-deployment.yaml
├── postgres-statefulset.yaml
├── ingress.yaml
├── configmap.yaml
└── secrets.yaml
```

## Testing Strategy

**Unit tests** in `core/application/`:
- Bracket generation correctness
- Advancement validation
- No-rematch seat assignment
- Yellow seat assignment
- Round-robin fairness pairing
- Scoring calculations

**Integration tests** (future):
- Full API flows
- Database operations with test containers

Tests follow TDD: write failing test first, implement to make it pass.

## Future Considerations

- WebSocket support for live tournament updates
- Email notifications for invitations
- ELO/rating system for seeding
- Advanced analytics dashboard