package http

import (
	"net/http"

	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
)

// CreateTournamentHandler handles tournament creation
//
//	@Summary		Create a new tournament
//	@Description	Creates a new knockout tournament with the given settings
//	@Tags			tournaments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		createTournamentRequest	true	"Tournament details"
//	@Success		201		{object}	models.Tournament
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/tournaments [post]
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
//
//	@Summary		List all tournaments
//	@Description	Retrieves a list of all tournaments
//	@Tags			tournaments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.Tournament
//	@Failure		500	{object}	errorResponse
//	@Router			/tournaments [get]
func ListTournamentsHandler(svc inbound.TournamentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "List tournaments not yet implemented",
		}})
	}
}

// GetTournamentHandler handles getting a single tournament
//
//	@Summary		Get a tournament by ID
//	@Description	Retrieves a tournament by its unique identifier
//	@Tags			tournaments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Tournament ID"
//	@Success		200	{object}	models.Tournament
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Router			/tournaments/{id} [get]
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
//
//	@Summary		Update a tournament
//	@Description	Updates an existing tournament's settings
//	@Tags			tournaments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"Tournament ID"
//	@Param			request	body		updateTournamentRequest	true	"Tournament settings"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/tournaments/{id} [patch]
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
//
//	@Summary		Delete a tournament
//	@Description	Deletes a tournament by its unique identifier
//	@Tags			tournaments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Tournament ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/tournaments/{id} [delete]
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
//
//	@Summary		List matches for a tournament
//	@Description	Retrieves all matches for a specific tournament
//	@Tags			tournaments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Tournament ID"
//	@Success		200	{array}		models.Match
//	@Failure		400	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/tournaments/{id}/matches [get]
func ListTournamentMatchesHandler(svc inbound.TournamentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "List tournament matches not yet implemented",
		}})
	}
}

// ReportMatchHandler handles reporting a match result
//
//	@Summary		Report a match result
//	@Description	Reports the result of a match in a tournament
//	@Tags			tournaments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Tournament ID"
//	@Param			request	body		reportMatchRequest	true	"Match result details"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/tournaments/{id}/matches [post]
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
//
//	@Summary		Get current round pairings
//	@Description	Retrieves the pairings for the current round of a tournament
//	@Tags			tournaments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Tournament ID"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/tournaments/{id}/pairings [get]
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
