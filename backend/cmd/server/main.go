package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ludohttp "ludo-tournament/adapters/primary/http"
	"ludo-tournament/adapters/secondary/persistence"
	"ludo-tournament/core/application"
	"ludo-tournament/core/domain"
	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/inbound"
	"ludo-tournament/core/ports/outbound"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	ctx := context.Background()

	// Initialize database
	db, err := persistence.NewPostgresDB(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := persistence.NewGormUserRepository(db)
	tournamentRepo := persistence.NewGormTournamentRepository(db)
	leagueRepo := persistence.NewGormLeagueRepository(db)
	matchRepo := persistence.NewGormMatchRepository(db)
	assignmentRepo := persistence.NewGormMatchAssignmentRepository(db)

	// Initialize application services
	tournamentService := application.NewTournamentService(tournamentRepo, matchRepo, assignmentRepo)
	leagueService := application.NewLeagueService(leagueRepo, matchRepo)
	authService := NewAuthService(userRepo)
	userService := NewUserService(userRepo)

	// Seed admin user if needed
	if err := seedAdminIfNeeded(ctx, userRepo); err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}

	// Setup router
	router := ludohttp.SetupRouter(
		tournamentService,
		leagueService,
		authService,
		userService,
	)

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// jwtSecret is used for signing and validating JWT tokens
var jwtSecret = []byte(getEnvOrFallback("JWT_SECRET", "fallback-secret-for-dev-only"))

func getEnvOrFallback(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// AuthService implements inbound.AuthService
type AuthService struct {
	userRepo outbound.UserRepository
}

func NewAuthService(userRepo outbound.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

var _ inbound.AuthService = (*AuthService)(nil)

func (s *AuthService) Register(ctx context.Context, email, password, inviteCode string) (string, error) {
	// TODO: Verify invite code with UserInviteRepository
	_ = inviteCode // placeholder

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Create user
	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleMember,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return "", err
	}

	return s.generateJWT(user.ID, string(user.Role))
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return "", domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrUnauthorized
	}

	return s.generateJWT(user.ID, string(user.Role))
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	// JWT is stateless; client discards token on logout
	return nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, token string) (*inbound.UserDTO, error) {
	userID, role, err := s.validateJWT(token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, domain.ErrUnauthorized
	}

	return &inbound.UserDTO{
		ID:    user.ID,
		Email: user.Email,
		Role:  role,
	}, nil
}

func (s *AuthService) generateJWT(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *AuthService) validateJWT(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", domain.ErrUnauthorized
	}

	userID, _ := claims["user_id"].(string)
	role, _ := claims["role"].(string)

	return userID, role, nil
}

// UserService implements inbound.UserService
type UserService struct {
	userRepo outbound.UserRepository
}

func NewUserService(userRepo outbound.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

var _ inbound.UserService = (*UserService)(nil)

func (s *UserService) ListUsers(ctx context.Context) ([]inbound.UserDTO, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]inbound.UserDTO, len(users))
	for i, u := range users {
		dtos[i] = inbound.UserDTO{ID: u.ID, Email: u.Email, Role: string(u.Role)}
	}
	return dtos, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, role string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		return domain.ErrNotFound
	}

	user.Role = models.Role(role)
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.SoftDelete(ctx, id)
}

func (s *UserService) SendInvite(ctx context.Context, email string, inviterID string, inviteType string) (string, error) {
	// TODO: Implement with UserInviteRepository
	return "placeholder-code-" + time.Now().Format("20060102150405"), nil
}

func (s *UserService) AcceptInvite(ctx context.Context, code string, email string, password string) (string, error) {
	// TODO: Implement with UserInviteRepository
	return "", nil
}

// seedAdminIfNeeded creates an initial admin user if no users exist
func seedAdminIfNeeded(ctx context.Context, userRepo outbound.UserRepository) error {
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

		log.Printf("Created initial admin user: %s", adminEmail)
		log.Printf("WARNING: Change this password immediately in production!")
	}

	return nil
}
