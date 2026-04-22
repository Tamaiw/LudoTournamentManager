# Development Guide

## Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| Go | 1.21+ | Backend runtime |
| Node.js | 18+ | Frontend build |
| npm | 9+ | Package management |
| Docker | 20+ | Local development |
| PostgreSQL | 16+ | Database (or use Docker) |

## Project Setup

### 1. Clone and Install Dependencies

```bash
# Backend
cd backend && go mod download

# Frontend
cd frontend && npm install
```

### 2. Start PostgreSQL

**Option A: Docker (recommended)**
```bash
docker run -d \
  --name ludo-postgres \
  -e POSTGRES_DB=ludo_tournament \
  -e POSTGRES_USER=ludo \
  -e POSTGRES_PASSWORD=changeme \
  -p 5432:5432 \
  postgres:16-alpine
```

**Option B: docker-compose**
```bash
docker-compose up -d postgres
```

### 3. Configure Environment

Set `DATABASE_URL` environment variable:

```bash
export DATABASE_URL=postgres://ludo:changeme@localhost:5432/ludo_tournament
```

### 4. Run Backend

```bash
cd backend && go run cmd/server/main.go
```

Backend starts on http://localhost:8080

### 5. Run Frontend

```bash
cd frontend && npm run dev
```

Frontend starts on http://localhost:5173

## Running Tests

### Backend Tests

```bash
cd backend && go test ./...
```

Run with verbose output:
```bash
go test -v ./core/application/...
```

Run with coverage:
```bash
go test -cover ./...
```

### Frontend Tests

```bash
cd frontend && npm test
```

### Full Stack (Docker Compose)

```bash
docker-compose up --build
```

Access:
- Frontend: http://localhost
- Backend: http://localhost:8080
- Database: localhost:5432

## Project Structure

```
.
├── backend/               # Go backend
│   ├── cmd/server/        # Entry point
│   ├── core/              # Business logic (no external deps)
│   │   ├── domain/         # Models, errors, events
│   │   ├── ports/         # Interfaces
│   │   └── application/   # Use cases, TDD tests
│   └── adapters/          # Infrastructure
│       ├── primary/http/   # HTTP handlers
│       └── secondary/      # GORM persistence
├── frontend/             # React + Vite + Tailwind
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── pages/        # Route pages
│   │   ├── hooks/        # Custom hooks
│   │   ├── services/     # API client
│   │   └── types/        # TypeScript types
│   └── ...
├── docs/                 # Documentation
└── k8s/                  # Kubernetes manifests
```

## Development Workflow

### 1. Make Changes

- Edit code in `backend/` or `frontend/`
- Backend uses hexagonal architecture: business logic in `core/application/`
- Frontend follows React conventions with hooks and React Query

### 2. Run Tests

```bash
# Backend unit tests (TDD)
cd backend && go test ./core/application/...

# Frontend
cd frontend && npm test
```

### 3. Rebuild

```bash
# Backend (auto-rebuilds with go run)
go build ./...

# Frontend
cd frontend && npm run build
```

## Key Conventions

### Go (Backend)

1. **Hexagonal architecture**: Core has no external dependencies
2. **Ports define boundaries**: Interfaces in `core/ports/` for all I/O
3. **TDD for business logic**: Tests in `*_test.go` files next to implementation
4. **Soft deletes**: All entities have `deleted_at` - no hard deletes
5. **UUIDs**: Use UUID for all entity IDs

### React (Frontend)

1. **TypeScript**: All components use TypeScript
2. **React Query**: Server state via TanStack Query
3. **Tailwind CSS**: Utility-first styling
4. **Hooks**: Custom hooks for reusable logic
5. **Components**: Functional components with props

## Database

### Schema Management

The system uses GORM's auto-migration. On startup, GORM automatically creates tables based on model definitions.

### Manual Database Operations

Connect to PostgreSQL:
```bash
docker exec -it ludo-postgres psql -U ludo -d ludo_tournament
```

Run migrations:
```bash
# Not implemented yet - using GORM auto-migrate
```

### Seeding (Future)

Admin endpoint `POST /admin/seed` for odd-player promotion from prior tournaments.

## Troubleshooting

### Backend Won't Start

1. Check PostgreSQL is running: `docker ps | grep postgres`
2. Verify DATABASE_URL: `echo $DATABASE_URL`
3. Check port 8080 is free: `lsof -i :8080`

### Frontend Build Fails

1. Clear node_modules: `rm -rf node_modules && npm install`
2. Check TypeScript errors: `npx tsc --noEmit`

### Database Connection Issues

1. Verify PostgreSQL credentials match DATABASE_URL
2. Check PostgreSQL logs: `docker logs ludo-postgres`
3. Test connection: `psql $DATABASE_URL -c "SELECT 1"`

## Code Style

### Go

- Run `gofmt` before committing: `gofmt -w .`
- Follow Go idioms: error handling, context propagation
- Comments on public functions/types

### TypeScript/React

- Run `npm run lint` (if configured)
- Prefer functional components and hooks
- Use TypeScript types - avoid `any`

## Common Tasks

### Add a New Tournament Feature

1. Add domain model in `backend/core/domain/models/`
2. Add outbound port in `backend/core/ports/outbound/`
3. Implement in `backend/core/application/`
4. Add HTTP handler in `backend/adapters/primary/http/`
5. Update frontend types and API service

### Add a New Frontend Component

1. Create in `frontend/src/components/` (tournament/, league/, ui/)
2. Add TypeScript props interface
3. Use existing API service in `frontend/src/services/api.ts`
4. Add route in `App.tsx`

## Debugging

### Backend

```bash
# Enable debug logging
DEBUG=1 go run cmd/server/main.go

# Attach Delve
dlv debug cmd/server/main.go
```

### Frontend

```bash
# Vite debug
npm run dev -- --debug

# React DevTools
# Install browser extension
```

## Deployment

### Docker

```bash
# Build images
docker-compose build

# Push to registry (production)
docker push registry.example.com/ludo-backend:latest
docker push registry.example.com/ludo-frontend:latest
```

### Kubernetes

```bash
# Apply manifests
kubectl apply -f k8s/

# Check status
kubectl get pods
kubectl get services
```

See `docs/ARCHITECTURE.md` for deployment architecture details.