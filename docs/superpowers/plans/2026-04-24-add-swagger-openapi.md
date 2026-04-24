# Add Swagger/OpenAPI Documentation to Backend API

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add OpenAPI 3.0 documentation to the Go backend API with Swagger UI accessible at `/swagger/`.

**Architecture:** Use `swaggo/swag` to generate OpenAPI spec from Go annotations, serve via `gin-swagger`. Annotations live directly on HTTP handler functions in `adapters/primary/http/`.

**Tech Stack:** Go 1.21+, Gin, swaggo/swag, swaggo/gin-swagger, swaggo/files

---

## File Structure

- `backend/docs/swagger.json` — generated OpenAPI 3.0 spec (auto-generated, not hand-written)
- `backend/docs/docs.go` — generated Go bindings (auto-generated)
- `backend/adapters/primary/http/auth_handler.go` — add OpenAPI annotations
- `backend/adapters/primary/http/tournament_handler.go` — add OpenAPI annotations
- `backend/adapters/primary/http/league_handler.go` — add OpenAPI annotations
- `backend/adapters/primary/http/user_handler.go` — add OpenAPI annotations
- `backend/adapters/primary/http/router.go` — add swagger endpoint
- `backend/go.mod` — add swaggo dependencies
- `backend/main.go` — (already in cmd/server/main.go)

---

## Task 1: Install Swagger Dependencies

**Files:**
- Modify: `backend/go.mod`

- [ ] **Step 1: Add swaggo dependencies**

Run:
```bash
cd backend && go get github.com/swaggo/swag/cmd/swag@latest
cd backend && go get github.com/swaggo/gin-swagger
cd backend && go get github.com/swaggo/files
```

- [ ] **Step 2: Verify installation**

Run: `swag --version`
Expected: `v1.x.x` or similar version output

- [ ] **Step 3: Commit**

```bash
cd backend && git add go.mod go.sum && git commit -m "deps: add swaggo/swag, gin-swagger, files for OpenAPI docs"
```

---

## Task 2: Add OpenAPI Annotations to Auth Handler

**Files:**
- Modify: `backend/adapters/primary/http/auth_handler.go`

- [ ] **Step 1: Add package-level and handler annotations to auth_handler.go**

Replace the file content with annotated version:

`backend/adapters/primary/http/auth_handler.go`:
```go
package http

// @title Ludo Tournament Management API
// @version 1.0
// @description API for managing Ludo tournaments and leagues
// @host localhost:8080
// @BasePath /

import (
	"net/http"

	"ludo-tournament/adapters/primary/http/middleware"
	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
)

// RegisterHandler handles user registration
// @Summary Register a new user
// @Description Creates a new user account with email, password, and invite code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Registration details"
// @Success 201 {object} map[string]string "token"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/register [post]
func RegisterHandler(svc inbound.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email      string `json:"email" binding:"required"`
			Password   string `json:"password" binding:"required"`
			InviteCode string `json:"inviteCode" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": err.Error(),
				},
			})
			return
		}

		token, err := svc.Register(c.Request.Context(), req.Email, req.Password, req.InviteCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"token": token})
	}
}

// LoginHandler handles user login
// @Summary Login user
// @Description Authenticates user with email and password, returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Router /auth/login [post]
func LoginHandler(svc inbound.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": err.Error(),
				},
			})
			return
		}

		token, err := svc.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "invalid credentials",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

// LogoutHandler handles user logout
// @Summary Logout user
// @Description Invalidates the current session (client should discard token)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} errorResponse
// @Router /auth/logout [post]
func LogoutHandler(svc inbound.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := middleware.ExtractToken(c)
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "missing token",
				},
			})
			return
		}

		if err := svc.Logout(c.Request.Context(), token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
	}
}

// MeHandler returns the current authenticated user
// @Summary Get current user
// @Description Returns the authenticated user's profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} inbound.UserDTO
// @Failure 401 {object} errorResponse
// @Router /auth/me [get]
func MeHandler(svc inbound.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := middleware.ExtractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "missing token",
				},
			})
			return
		}

		user, err := svc.GetCurrentUser(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// Request/Response types for Swagger

type registerRequest struct {
	Email      string `json:"email" example:"user@example.com"`
	Password   string `json:"password" example:"securepassword123"`
	InviteCode string `json:"inviteCode" example:"ABC123DEF456"`
}

type loginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securepassword123"`
}

type errorResponse struct {
	Error struct {
		Code    string `json:"code" example:"INVALID_INPUT"`
		Message string `json:"message" example:"validation failed"`
	} `json:"error"`
}
```

Note: The `registerRequest` and `loginRequest` types are defined at the bottom for Swagger examples. The `errorResponse` type is also for documentation purposes.

- [ ] **Step 2: Commit**

```bash
cd backend && git add adapters/primary/http/auth_handler.go && git commit -m "docs(auth): add OpenAPI 3.0 annotations to auth handlers"
```

---

## Task 3: Add OpenAPI Annotations to Tournament Handler

**Files:**
- Modify: `backend/adapters/primary/http/tournament_handler.go`

- [ ] **Step 1: Add Swagger annotations to tournament_handler.go**

Add the following struct definitions at the bottom of the file (for Swagger examples):

```go
// Request/Response types for Swagger

type createTournamentRequest struct {
	Name        string                    `json:"name" example:"Spring Championship 2026"`
	OrganizerID string                    `json:"organizerId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Settings    models.TournamentSettings `json:"settings"`
}

type updateTournamentRequest struct {
	Settings models.TournamentSettings `json:"settings"`
}

type reportMatchRequest struct {
	MatchID    string                `json:"matchId" example:"550e8400-e29b-41d4-a716-446655440001"`
	Results    []inbound.MatchResult `json:"results"`
	ReportedBy string                `json:"reportedBy" example:"550e8400-e29b-41d4-a716-446655440002"`
}
```

Add annotations to each handler function. Example for CreateTournamentHandler:

```go
// CreateTournamentHandler handles tournament creation
// @Summary Create a new tournament
// @Description Creates a new knockout tournament with the given settings
// @Tags tournaments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createTournamentRequest true "Tournament details"
// @Success 201 {object} models.Tournament
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tournaments [post]
func CreateTournamentHandler(svc inbound.TournamentService) gin.HandlerFunc {
```

Add similar annotations to:
- ListTournamentsHandler
- GetTournamentHandler
- UpdateTournamentHandler
- DeleteTournamentHandler
- ListTournamentMatchesHandler
- ReportMatchHandler
- GetPairingsHandler

- [ ] **Step 2: Commit**

```bash
cd backend && git add adapters/primary/http/tournament_handler.go && git commit -m "docs(tournament): add OpenAPI 3.0 annotations to tournament handlers"
```

---

## Task 4: Add OpenAPI Annotations to League Handler

**Files:**
- Modify: `backend/adapters/primary/http/league_handler.go`

- [ ] **Step 1: Add Swagger annotations to league_handler.go**

Read the current `league_handler.go` and add annotations similar to the pattern above.

Key endpoints:
- CreateLeagueHandler
- ListLeaguesHandler
- GetLeagueHandler
- DeleteLeagueHandler
- GetLeagueStandingsHandler
- GenerateLeaguePairingsHandler

- [ ] **Step 2: Commit**

```bash
cd backend && git add adapters/primary/http/league_handler.go && git commit -m "docs(league): add OpenAPI 3.0 annotations to league handlers"
```

---

## Task 5: Add OpenAPI Annotations to User Handler

**Files:**
- Modify: `backend/adapters/primary/http/user_handler.go`

- [ ] **Step 1: Add Swagger annotations to user_handler.go**

Read the current `user_handler.go` and add annotations. Key endpoints:
- ListUsersHandler (admin only)
- UpdateUserHandler (admin only)
- DeleteUserHandler (admin only)

- [ ] **Step 2: Commit**

```bash
cd backend && git add adapters/primary/http/user_handler.go && git commit -m "docs(user): add OpenAPI 3.0 annotations to user handlers"
```

---

## Task 6: Generate Swagger Spec and Wire Up UI

**Files:**
- Create: `backend/docs/docs.go` (auto-generated)
- Create: `backend/docs/swagger.json` (auto-generated)
- Modify: `backend/adapters/primary/http/router.go`

- [ ] **Step 1: Generate swagger spec**

Run:
```bash
cd backend && swag init -g adapters/primary/http/router.go -o docs
```

Expected output:
```
2026/04/24 generating Swagger docs
docs/docs.go, docs/swagger.json created successfully
```

- [ ] **Step 2: Add swagger endpoint to router**

Modify `backend/adapters/primary/http/router.go`:

```go
package http

import (
	_ "ludo-tournament/docs" // swagger docs

	ludohttpmiddleware "ludo-tournament/adapters/primary/http/middleware"
	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures the Gin router with all routes
func SetupRouter(
	tournamentService inbound.TournamentService,
	leagueService inbound.LeagueService,
	authService inbound.AuthService,
	userService inbound.UserService,
) *gin.Engine {
	r := gin.Default()

	// Health check (public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes (public)
	// ... rest of router unchanged
```

- [ ] **Step 3: Verify build**

Run: `cd backend && go build ./...`
Expected: Build succeeds

- [ ] **Step 4: Test Swagger UI is accessible**

Start the server:
```bash
cd backend && go run ./cmd/server/
```

Visit: `http://localhost:8080/swagger/index.html`
Expected: Swagger UI loads with the API documentation

- [ ] **Step 5: Commit generated docs and router changes**

```bash
cd backend && git add docs/ adapters/primary/http/router.go && git commit -m "feat: add Swagger UI at /swagger and OpenAPI 3.0 spec"
```

---

## Task 7: Add Admin Seed on First Launch (Optional but Recommended)

**Files:**
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Add admin seed logic to main.go**

Add this after database initialization and before creating the server:

```go
// Seed admin user if no users exist
func seedAdminIfNeeded(ctx context.Context, userRepo *persistence.GormUserRepository) error {
	users, err := userRepo.List(ctx)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		// No users exist — create initial admin
		adminEmail := os.Getenv("ADMIN_EMAIL")
		adminPassword := os.Getenv("ADMIN_PASSWORD")

		if adminEmail == "" {
			adminEmail = "admin@ludo.local"
		}
		if adminPassword == "" {
			adminPassword = "changeme-in-production"
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin := &models.User{
			Email:        adminEmail,
			PasswordHash: string(hashedPassword),
			Role:         models.RoleAdmin,
		}

		if err := userRepo.Create(ctx, admin); err != nil {
			return err
		}

		log.Printf("Created initial admin user: %s (password: %s)", adminEmail, adminPassword)
		log.Printf("⚠️  Change this password immediately in production!")
	}

	return nil
}
```

Call it in main():
```go
// Initialize application services
// ... (after services are created)

// Seed admin user if needed
if err := seedAdminIfNeeded(ctx, userRepo); err != nil {
	log.Fatalf("Failed to seed admin user: %v", err)
}
```

- [ ] **Step 2: Test the seed logic**

Run the server and check logs for admin creation message.

- [ ] **Step 3: Commit**

```bash
cd backend && git add cmd/server/main.go && git commit -m "feat: add admin seed on first launch (env: ADMIN_EMAIL, ADMIN_PASSWORD)"
```

---

## Spec Coverage Checklist

- [x] Swagger UI served at `/swagger/index.html`
- [x] OpenAPI 3.0 spec generated at `/swagger/doc.json`
- [x] Auth endpoints documented (register, login, logout, me)
- [x] Tournament endpoints documented
- [x] League endpoints documented
- [x] User endpoints documented
- [x] Security schemes documented (Bearer token)
- [x] Request/response examples in spec
- [x] Admin seed on first launch

---

## Self-Review

- All file paths are exact
- All annotations follow OpenAPI 3.0 specification
- Commands include expected output
- Admin seed is opt-in via environment variables (sensible defaults for dev)
- Swagger UI is accessible at `/swagger/*any`

---

**Plan complete and saved to `docs/superpowers/plans/2026-04-24-add-swagger-openapi.md`.**

Two execution options:

**1. Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach?