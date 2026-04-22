package http

import (
	"net/http"

	"ludo-tournament/adapters/primary/http/middleware"
	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
)

// RegisterHandler handles user registration
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
func LogoutHandler(svc inbound.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := middleware.extractToken(c)
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
func MeHandler(svc inbound.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := middleware.extractToken(c)
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
