package http

import (
	"net/http"

	"ludo-tournament/adapters/primary/http/middleware"
	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
)

// CreateTournamentHandler handles tournament creation
func CreateTournamentHandler(svc inbound.TournamentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string                    `json:"name" binding:"required"`
			OrganizerID string                    `json:"organizerId" binding:"required"`
			Settings    models.TournamentSettings `json:"settings"`
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

		tournament, err := svc.CreateTournament(c.Request.Context(), req.Name, req.OrganizerID, req.Settings)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusCreated, tournament)
	}
}

// ListTournamentsHandler handles listing tournaments
func ListTournamentsHandler(svc inbound.TournamentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "List tournaments not yet implemented",
		}})
	}
}

// GetTournamentHandler handles getting a single tournament
func GetTournamentHandler(svc inbound.TournamentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "tournament id is required",
				},
			})
			return
		}

		tournament, err := svc.GetTournament(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "TOURNAMENT_NOT_FOUND",
					"message": "Tournament with ID " + id + " not found",
				},
			})
			return
		}

		c.JSON(http.StatusOK, tournament)
	}
}

// UpdateTournamentHandler handles updating a tournament
func UpdateTournamentHandler(svc inbound.TournamentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "tournament id is required",
				},
			})
			return
		}

		var req struct {
			Settings models.TournamentSettings `json:"settings"`
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

		if err := svc.UpdateTournament(c.Request.Context(), id, req.Settings); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "tournament updated"})
	}
}

// DeleteTournamentHandler handles deleting a tournament
func DeleteTournamentHandler(svc inbound.TournamentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "tournament id is required",
				},
			})
			return
		}

		if err := svc.DeleteTournament(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "tournament deleted"})
	}
}

// ListTournamentMatchesHandler handles listing matches for a tournament
func ListTournamentMatchesHandler(svc inbound.TournamentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "List tournament matches not yet implemented",
		}})
	}
}

// ReportMatchHandler handles reporting a match result
func ReportMatchHandler(svc inbound.TournamentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tournamentID := c.Param("id")
		if tournamentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "tournament id is required",
				},
			})
			return
		}

		var req struct {
			MatchID    string                `json:"matchId" binding:"required"`
			Results    []inbound.MatchResult `json:"results" binding:"required"`
			ReportedBy string                `json:"reportedBy" binding:"required"`
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

		if err := svc.ReportMatch(c.Request.Context(), req.MatchID, req.Results, req.ReportedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "match reported"})
	}
}

// GetPairingsHandler handles getting current round pairings
func GetPairingsHandler(svc inbound.TournamentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tournamentID := c.Param("id")
		if tournamentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "tournament id is required",
				},
			})
			return
		}

		pairings, err := svc.GetCurrentRoundPairings(c.Request.Context(), tournamentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"pairings": pairings})
	}
}
