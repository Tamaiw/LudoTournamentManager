# Ludo Tournament Management System

A web application for managing knockout tournaments and round-robin leagues for the board game Ludo. Supports tournaments of ~100 players and long-running league systems with customizable scoring.

## Features

### Knockout Tournaments
- **Configurable advancement** - Define how many players advance from each game per round
- **No rematches** - Players from the same source game are distributed across different tables in the next round
- **Yellow seat tracking** - First-come, first-served assignment based on 1st-place finish time
- **Partial completion** - Games can start before all games in the previous round are complete
- **Edit lock** - Prevents editing games if downstream games have already been played

### Round-Robin Leagues
- **Fairness-aware pairing** - Minimizes repeat matches across the league
- **Customizable scoring** - Define point values per placement (e.g., 1st=3, 2nd=2, 3rd=1, 4th=0)
- **Manual swap with warnings** - Organizers can swap players with fairness alerts
- **Tiebreaker support** - Single tiebreaker game when players are tied

### User Management
- **Role-based access** - Admin, Member, Guest roles with appropriate permissions
- **Invite flow** - Email-based invitation system with unique registration codes
- **Player profiles** - Tournament history and league statistics

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.21+ with Gin HTTP framework |
| Database | PostgreSQL 16 with GORM ORM |
| Frontend | React 18 with Vite and Tailwind CSS |
| State | React Query (TanStack Query) for server state |
| Architecture | Hexagonal (Ports & Adapters) |
| Deployment | Docker, Kubernetes |

## Project Structure

```
├── backend/               # Go backend (hexagonal architecture)
│   ├── cmd/server/        # Entry point
│   ├── core/              # Domain logic (no external dependencies)
│   │   ├── domain/         # Models, errors, events
│   │   ├── ports/         # Interface definitions
│   │   └── application/   # Business logic
│   └── adapters/          # Infrastructure adapters
│       ├── primary/http/   # HTTP handlers and middleware
│       └── secondary/      # Database adapters
├── frontend/             # React + Vite + Tailwind
├── docs/                  # Documentation
│   ├── ARCHITECTURE.md    # Detailed architecture guide
│   ├── API.md             # API reference
│   └── DEVELOPMENT.md     # Development setup and workflow
├── k8s/                   # Kubernetes manifests
└── docker-compose.yml     # Local development stack
```

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker 20+
- PostgreSQL 16 (or use Docker Compose)

### Local Development

**1. Start PostgreSQL and backend:**

```bash
docker compose up -d postgres
cd backend && go run cmd/server/main.go
```

**2. Start frontend:**

```bash
cd frontend && npm install && npm run dev
```

**3. Access the application:**

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

### Docker Compose (Full Stack)

```bash
docker compose up
```

This starts PostgreSQL, backend (port 8080), and frontend (port 80).

## Documentation

| Document | Description |
|----------|-------------|
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Hexagonal architecture, data model, design decisions |
| [API.md](docs/API.md) | REST API endpoints and authentication |
| [DEVELOPMENT.md](docs/DEVELOPMENT.md) | Setup, testing, workflow guidelines |

## Architecture

The system follows **hexagonal architecture** (ports and adapters):

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend (React)                      │
└─────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│              HTTP Adapters (Gin Router)                  │
│         Auth middleware, handlers, router                │
└─────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│              Inbound Ports (Interfaces)                 │
│    TournamentService, LeagueService, AuthService        │
└─────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│            Application Layer (Business Logic)            │
│       Tournament, League, Pairing, Auth logic           │
└─────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│             Outbound Ports (Interfaces)                  │
│   UserRepository, TournamentRepository, MatchRepository  │
└─────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│           Persistence Adapters (GORM)                   │
│            PostgreSQL repositories                       │
└─────────────────────────────────────────────────────────┘
```

**Key principle:** Core domain has no dependencies on infrastructure. All external I/O happens through interfaces.

## API Overview

### Authentication

```
POST /auth/register  - Create account (requires invite code)
POST /auth/login     - Login with email/password
POST /auth/logout    - Logout
GET  /auth/me        - Current user profile
```

### Tournaments

```
POST   /tournaments/:id/matches  - Report match result
GET    /tournaments/:id/pairings - Get current round pairings
```

### Leagues

```
POST   /leagues/:id/pairings/generate  - Generate fair pairings
GET    /leagues/:id/standings          - Get league standings
```

## License

MIT