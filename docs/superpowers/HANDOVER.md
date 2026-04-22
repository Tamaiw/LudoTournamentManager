# Ludo Tournament Management System — Status

> This document is updated periodically as tasks complete. It serves as a handover checkpoint if context runs low.

## Project Overview

**Repository:** https://github.com/Tamaiw/LudoTournamentManager
**Main branch:** main
**Status:** Implementation in progress — Tasks 1-5 complete, starting Task 6

## Architecture

- **Backend:** Go 1.22+, Gin, GORM, PostgreSQL 16
- **Frontend:** React 18+, Vite, Tailwind CSS
- **Pattern:** Hexagonal architecture (ports & adapters)
- **Deployment:** Docker + Kubernetes

## Documentation Files

| File | Purpose |
|------|---------|
| `docs/superpowers/specs/2026-04-21-ludo-tournament-management-system-design.md` | Design specification |
| `docs/superpowers/plans/2026-04-21-ludo-tournament-management-system.md` | Implementation plan |

## Worktrees (Isolated Task Workspaces)

Each task has its own worktree directory with an isolated branch:

| Worktree | Branch | Task |
|----------|--------|------|
| `.worktrees/task-1-scaffold` | task-1-scaffold | Initialize Go backend structure |
| `.worktrees/task-2-ports` | task-2-ports | Define ports (interfaces) |
| `.worktrees/task-3-tournament` | task-3-tournament | Tournament bracket & advancement |
| `.worktrees/task-4-pairing` | task-4-pairing | Table assignment & no-rematch |
| `.worktrees/task-5-league` | task-5-league | League round-robin & fairness |
| `.worktrees/task-6-persistence` | task-6-persistence | GORM persistence layer |
| `.worktrees/task-7-http` | task-7-http | HTTP handlers & router |
| `.worktrees/task-8-frontend` | task-8-frontend | React + Vite + Tailwind setup |
| `.worktrees/task-9-components` | task-9-components | Tournament UI components |
| `.worktrees/task-10-docker` | task-10-docker | Docker & Kubernetes setup |

## Key Design Decisions

### Tournament Advancement
- Configurable per-game advancement: some games allow 2-spot, others 3-spot
- System validates that total advancing players = next round spots
- Yellow seat: first-come, first-served based on 1st-place finish time
- No immediate rematches: players from same source game distributed across different tables
- Partial completion: game can start when all 4 spots filled, even if other round games pending
- Edit lock: cannot edit game N if downstream games already played

### League System
- Round-robin format with customizable scoring rules
- Random pairing with fairness tracking (minimizes repeat matches)
- Organizer can manually swap players with warnings for uneven matchups
- Tiebreaker: default is most wins, organizer can add single tiebreaker game

### User Roles
- Admin: full system control
- Member: participate, organize, invite guests
- Guest: accept invitations, register results

## Pending Tasks (in order)

- [x] Task 1: Initialize Go backend structure (domain models, errors, events)
- [x] Task 2: Define ports (interfaces) for hexagonal architecture
- [x] Task 3: Tournament bracket generation & advancement logic (TDD)
- [x] Task 4: Table assignment & no-rematch logic (TDD)
- [x] Task 5: League round-robin & fairness-aware pairing (TDD)
- [ ] Task 6: GORM persistence layer
- [ ] Task 7: HTTP handlers & Gin router
- [ ] Task 8: React + Vite + Tailwind frontend scaffolding
- [ ] Task 9: Tournament components (MatchCard, TableAssignment, etc.)
- [ ] Task 10: Docker & Kubernetes setup

## Open Questions / Notes

- SSH key set up for GitHub push (passphrase-protected key)
- Node.js v18.19.1 available
- Docker v29.1.3 available
- Go v1.22.2 available

## Last Updated

Context: 2026-04-22
