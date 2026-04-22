package http

import (
	"ludo-tournament/adapters/primary/http/middleware"
	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
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

	// Auth routes (public)
	auth := r.Group("/auth")
	{
		auth.POST("/register", RegisterHandler(authService))
		auth.POST("/login", LoginHandler(authService))
		auth.POST("/logout", LogoutHandler(authService))
		auth.GET("/me", AuthMiddleware(), MeHandler(authService))
	}

	// Protected routes
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		// Tournaments
		api.POST("/tournaments", CreateTournamentHandler(tournamentService))
		api.GET("/tournaments", ListTournamentsHandler(tournamentService))
		api.GET("/tournaments/:id", GetTournamentHandler(tournamentService))
		api.PATCH("/tournaments/:id", UpdateTournamentHandler(tournamentService))
		api.DELETE("/tournaments/:id", DeleteTournamentHandler(tournamentService))
		api.GET("/tournaments/:id/matches", ListTournamentMatchesHandler(tournamentService))
		api.POST("/tournaments/:id/matches", ReportMatchHandler(tournamentService))
		api.GET("/tournaments/:id/pairings", GetPairingsHandler(tournamentService))

		// Leagues
		api.POST("/leagues", CreateLeagueHandler(leagueService))
		api.GET("/leagues", ListLeaguesHandler(leagueService))
		api.GET("/leagues/:id", GetLeagueHandler(leagueService))
		api.DELETE("/leagues/:id", DeleteLeagueHandler(leagueService))
		api.GET("/leagues/:id/standings", GetLeagueStandingsHandler(leagueService))
		api.POST("/leagues/:id/pairings/generate", GenerateLeaguePairingsHandler(leagueService))

		// Users (admin only)
		api.GET("/users", middleware.RequireRole("admin"), ListUsersHandler(userService))
		api.PATCH("/users/:id", middleware.RequireRole("admin"), UpdateUserHandler(userService))
		api.DELETE("/users/:id", middleware.RequireRole("admin"), DeleteUserHandler(userService))
	}

	return r
}
