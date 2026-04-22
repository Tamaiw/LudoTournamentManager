package http

import (
	"net/http"

	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
)

// ListUsersHandler handles listing all users (admin only)
func ListUsersHandler(svc inbound.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := svc.ListUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

// UpdateUserHandler handles updating a user (admin only)
func UpdateUserHandler(svc inbound.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "user id is required",
				},
			})
			return
		}

		var req struct {
			Role string `json:"role" binding:"required"`
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

		if err := svc.UpdateUser(c.Request.Context(), id, req.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user updated"})
	}
}

// DeleteUserHandler handles soft-deleting a user (admin only)
func DeleteUserHandler(svc inbound.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_INPUT",
					"message": "user id is required",
				},
			})
			return
		}

		if err := svc.DeleteUser(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	}
}
