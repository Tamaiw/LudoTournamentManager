package http

import (
	"net/http"

	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
)

// CreateLeagueHandler handles league creation
func CreateLeagueHandler(svc inbound.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string                `json:"name" binding:"required"`
			OrganizerID string                `json:"organizerId" binding:"required"`
			Settings    models.LeagueSettings `json:"settings"`
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

		league, err := svc.CreateLeague(c.Request.Context(), req.Name, req.OrganizerID, req.Settings)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusCreated, league)
	}
}

// ListLeaguesHandler handles listing leagues
func ListLeaguesHandler(svc inbound.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Add filtering by status query param when service supports it
		c.JSON(http.StatusOK, gin.H{
			"leagues": []models.League{},
		})
	}
}

// GetLeagueHandler handles getting a single league
func GetLeagueHandler(svc inbound.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "league id is required",
				},
			})
			return
		}

		league, err := svc.GetLeague(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "LEAGUE_NOT_FOUND",
					"message": "League with ID " + id + " not found",
				},
			})
			return
		}

		c.JSON(http.StatusOK, league)
	}
}

// DeleteLeagueHandler handles deleting a league
func DeleteLeagueHandler(svc inbound.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "league id is required",
				},
			})
			return
		}

		if err := svc.DeleteLeague(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "league deleted"})
	}
}

// GetLeagueStandingsHandler handles getting league standings
func GetLeagueStandingsHandler(svc inbound.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "league id is required",
				},
			})
			return
		}

		standings, err := svc.GetStandings(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"standings": standings})
	}
}

// GenerateLeaguePairingsHandler handles generating pairings for a league
func GenerateLeaguePairingsHandler(svc inbound.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "league id is required",
				},
			})
			return
		}

		var req struct {
			PlayDate string `json:"playDate" binding:"required"`
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

		pairings, err := svc.GeneratePairings(c.Request.Context(), id, req.PlayDate)
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
