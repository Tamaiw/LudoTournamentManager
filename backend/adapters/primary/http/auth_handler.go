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