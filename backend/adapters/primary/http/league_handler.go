package http

import (
	"net/http"

	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
)

// CreateLeagueHandler handles league creation
//
//	@Summary		Create a new league
//	@Description	Creates a new league with the given settings
//	@Tags			leagues
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		createLeagueRequest	true	"League details"
//	@Success		201		{object}	models.League
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/leagues [post]
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
//
//	@Summary		List all leagues
//	@Description	Retrieves a list of all leagues
//	@Tags			leagues
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.League
//	@Failure		501	{object}	errorResponse
//	@Router			/leagues [get]
func ListLeaguesHandler(svc inbound.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "List leagues not yet implemented",
		}})
	}
}

// GetLeagueHandler handles getting a single league
//
//	@Summary		Get a league by ID
//	@Description	Retrieves a league by its unique identifier
//	@Tags			leagues
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"League ID"
//	@Success		200	{object}	models.League
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Router			/leagues/{id} [get]
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
//
//	@Summary		Delete a league
//	@Description	Deletes a league by its unique identifier
//	@Tags			leagues
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"League ID"
//	@Success		200	{object}	deleteLeagueResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/leagues/{id} [delete]
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
//
//	@Summary		Get league standings
//	@Description	Retrieves the standings for a league
//	@Tags			leagues
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"League ID"
//	@Success		200	{object}	leagueStandingsResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/leagues/{id}/standings [get]
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
//
//	@Summary		Generate league pairings
//	@Description	Generates pairings for a league for a specific play date
//	@Tags			leagues
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"League ID"
//	@Param			request	body		generateLeaguePairingsRequest	true	"Pairing details"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/leagues/{id}/pairings [post]
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
