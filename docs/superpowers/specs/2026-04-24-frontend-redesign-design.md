# Frontend Redesign — Design Specification

## Overview

Redesign the Ludo Tournament Management System frontend to provide a proper authenticated experience: login flow, role-aware dashboard, tournament/league management, profile, and admin panel.

---

## 1. Page Structure & Routing

### Routes

| Path | Component | Access |
|------|-----------|--------|
| `/login` | `LoginPage` | public |
| `/register` | `RegisterPage` | public (requires invite code) |
| `/` | `DashboardPage` | authenticated |
| `/tournaments` | `TournamentListPage` | authenticated |
| `/tournaments/new` | `CreateTournamentPage` | member+ |
| `/tournaments/:id` | `TournamentDetailPage` | authenticated |
| `/leagues` | `LeagueListPage` | authenticated |
| `/leagues/new` | `CreateLeaguePage` | member+ |
| `/leagues/:id` | `LeagueDetailPage` | authenticated |
| `/profile` | `ProfilePage` | authenticated |
| `/admin` | `AdminPage` | admin |
| `/admin/users` | `UserManagementPage` | admin |

### Layout Hierarchy

```
<App>
  <AuthLayout>          // login, register
    <Outlet />
  </AuthLayout>
  <AppLayout>           // authenticated shell
    <Navbar />
    <Outlet />          // page content
  </AppLayout>
</App>
```

---

## 2. Auth Flow

### Login Page (`/login`)
- Email + password form
- "Sign in" button
- Inline error display for invalid credentials
- Link to register page

### Register Page (`/register`)
- Email + password + invite code form
- "Create account" button
- Inline error for invalid invite code / email conflict

### Logout
- `POST /auth/logout` on navbar logout click
- Clear token, redirect to `/login`

### Session Persistence
- JWT stored in memory (component state) — not localStorage
- On page refresh: `GET /auth/me` to rehydrate user
- On app mount: attempt `GET /auth/me` to restore session

### Auth Context (`useAuth`)
```typescript
interface AuthContext {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, inviteCode: string) => Promise<void>;
  logout: () => Promise<void>;
}
```

### Protected Route Wrapper
- Redirects to `/login` if `user === null` and not loading
- Admin route wrapper: redirects to `/` if user is not admin

---

## 3. Dashboard

**Route:** `/` (authenticated)

### Layout: Single column, max-width container

### Sections (top to bottom)

**1. Welcome Header**
- "Welcome back, {email}" or first name if available
- Current date

**2. Next Event Card (prominent)**
- Shows nearest upcoming `live` tournament OR league
- Card: name, type badge (Tournament/League), next play date, "Enter" button
- If no live events: empty state with "No active events" + link to create

**3. Quick Actions**
Role-aware buttons:

| Role | Actions |
|------|---------|
| `guest` | Browse Tournaments, Browse Leagues |
| `member` | Create Tournament, Create League, Browse Tournaments, Browse Leagues |
| `admin` | Same as member + Admin Panel link |

**4. Recent Activity**
- Last 3 tournaments and last 3 leagues user is associated with
- Simple list: name, status badge, date

**5. Admin Panel Link**
- Visible only to `admin` role
- Card linking to `/admin/users`

### API Calls on Load
- `GET /tournaments?status=live`
- `GET /leagues?status=live`
- `GET /tournaments` (paginated, recent)
- `GET /leagues` (paginated, recent)

---

## 4. Tournament List & Detail

### Tournament List (`/tournaments`)
- Page header: "Tournaments" + "Create Tournament" button (member+ only)
- Filter tabs: All | Live | Completed | Draft
- Card grid: name, status badge, created date, organizer
- Click → navigate to `/tournaments/:id`

### Create Tournament (`/tournaments/new`)
- Form: name, tables count, advancement config (simplified for v1)
- `POST /tournaments`
- Success → redirect to `/tournaments/:id`

### Tournament Detail (`/tournaments/:id`)
- Header: name, status badge, Edit button (organizer+ only)
- Tabs:
  - **Bracket** — bracket view (use existing `BracketView` component)
  - **Matches** — list of all matches with status, report action
  - **Players** — participant list
- Live tournament: show current round pairings via `GET /tournaments/:id/pairings`

### API Calls
| Action | Endpoint |
|--------|----------|
| List | `GET /tournaments` |
| Create | `POST /tournaments` |
| Detail | `GET /tournaments/:id` |
| Pairings | `GET /tournaments/:id/pairings` |
| Matches | `GET /tournaments/:id/matches` |
| Report | `POST /tournaments/:id/matches` |

---

## 5. League List & Detail

### League List (`/leagues`)
- Page header: "Leagues" + "Create League" button (member+ only)
- Filter tabs: All | Live | Completed | Draft
- Card grid: name, status badge, scoring rules summary, created date
- Click → navigate to `/leagues/:id`

### Create League (`/leagues/new`)
- Form: name, scoring rules editor, games per player, tables count
- `POST /leagues`
- Success → redirect to `/leagues/:id`

### League Detail (`/leagues/:id`)
- Header: name, status badge, scoring rules display
- Tabs:
  - **Standings** — table via `GET /leagues/:id/standings` (rank, player, points, wins, games played)
  - **Generate Pairings** — button to trigger `POST /leagues/:id/pairings/generate`
- Play dates/schedule section (visible but note: `POST /leagues/:id/play-dates` not yet implemented)

### API Calls
| Action | Endpoint |
|--------|----------|
| List | `GET /leagues` |
| Create | `POST /leagues` |
| Detail | `GET /leagues/:id` |
| Standings | `GET /leagues/:id/standings` |
| Generate Pairings | `POST /leagues/:id/pairings/generate` |

---

## 6. Profile & Admin

### Profile Page (`/profile`)
- User info: email, role badge, member since date
- Section: "Your Tournaments" — list of user's tournaments
- Section: "Your Leagues" — list of user's leagues
- Logout button at bottom
- API: `GET /auth/me`

### Admin Page (`/admin`)
- Tab navigation: Users | Invitations (future)
- **Users tab:**
  - `GET /users` → table: email, role badge, last active, actions
  - Role change: `PATCH /users/:id` with `{ role }`
  - Delete: `DELETE /users/:id` (soft delete)
- **Invitations tab:** placeholder, note: backend not yet implemented

### API Calls (Admin)
| Action | Endpoint |
|--------|----------|
| List users | `GET /users` |
| Update user | `PATCH /users/:id` |
| Delete user | `DELETE /users/:id` |

---

## 7. Component Inventory

### New Components Needed

**Layout:**
- `AppLayout` — authenticated shell with navbar
- `Navbar` — logo, nav links, user menu (profile, logout), admin link
- `AuthLayout` — centered card for login/register

**UI (shared):**
- `Button` — variants: primary, secondary, danger, ghost
- `Input` — text input with label and error state
- `Badge` — status badges (live, draft, completed)
- `Card` — container with optional header
- `Modal` — dialog overlay
- `Tabs` — tab navigation

**Auth:**
- `LoginForm`
- `RegisterForm`
- `ProtectedRoute`
- `AdminRoute`

**Dashboard:**
- `NextEventCard`
- `QuickActions`
- `RecentActivityList`

**Tournament:**
- `TournamentCard`
- `TournamentForm`
- `BracketView` (exists, integrate)

**League:**
- `LeagueCard`
- `LeagueForm`
- `StandingsTable`

**Admin:**
- `UserTable`
- `RoleSelect`

---

## 8. Technical Approach

### Stack
- React 18 + Vite (existing)
- React Router v6 for routing
- TanStack Query (existing) for server state
- Tailwind CSS (existing)

### State Management
- Auth state: React Context (`AuthContext`)
- Server state: TanStack Query hooks (`useTournaments`, `useLeague`, etc.)
- No additional state library needed

### API Layer
- Extend existing `services/api.ts` with:
  - `getMe()`
  - `listUsers()`, `updateUser()`, `deleteUser()`
  - `createLeague()`, `getLeague()`, `getLeagueStandings()`
  - `generateLeaguePairings()`

### File Structure
```
frontend/src/
├── components/
│   ├── layout/          # AppLayout, Navbar, AuthLayout
│   ├── ui/              # Button, Input, Badge, Card, Modal, Tabs
│   ├── auth/            # LoginForm, RegisterForm
│   ├── dashboard/       # NextEventCard, QuickActions
│   ├── tournament/       # TournamentCard, TournamentForm
│   ├── league/          # LeagueCard, LeagueForm, StandingsTable
│   └── admin/            # UserTable, RoleSelect
├── pages/
│   ├── LoginPage.tsx
│   ├── RegisterPage.tsx
│   ├── DashboardPage.tsx
│   ├── TournamentListPage.tsx
│   ├── CreateTournamentPage.tsx
│   ├── TournamentDetailPage.tsx
│   ├── LeagueListPage.tsx
│   ├── CreateLeaguePage.tsx
│   ├── LeagueDetailPage.tsx
│   ├── ProfilePage.tsx
│   ├── AdminPage.tsx
│   └── UserManagementPage.tsx
├── hooks/
│   ├── useAuth.tsx      # Auth context provider
│   ├── useTournaments.ts
│   ├── useTournament.ts
│   ├── useLeagues.ts
│   ├── useLeague.ts
│   └── useUsers.ts
├── context/
│   └── AuthContext.tsx
└── types/
    └── index.ts         # extend with MatchResult, ScoringRule, etc.
```

---

## 9. Dependencies

No new runtime dependencies needed — all required packages likely already present or simple to add:
- `react-router-dom` — routing (check if v6 is installed)
- `@tanstack/react-query` — already in use

---

## 10. Out of Scope (v1)

- Invitation system (backend missing)
- Play dates management (endpoint missing)
- League status update (endpoint missing)
- Email notifications
- WebSocket live updates
- ELO/rating system
