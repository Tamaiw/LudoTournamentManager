package http

import (
	"net/http"

	"ludo-tournament/core/ports/inbound"

	"github.com/gin-gonic/gin"
)

// ListUsersHandler handles listing all users (admin only)
//
//	@Summary		List all users
//	@Description	Retrieves a list of all users (admin only)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	listUsersResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/users [get]
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

		c.JSON(http.StatusOK, listUsersResponse{Users: users})
	}
}

// UpdateUserHandler handles updating a user (admin only)
//
//	@Summary		Update a user
//	@Description	Updates an existing user's role (admin only)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"User ID"
//	@Param			request	body		updateUserRequest	true	"User details"
//	@Success		200		{object}	updateUserResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/users/{id} [patch]
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

		c.JSON(http.StatusOK, updateUserResponse{Message: "user updated"})
	}
}

// DeleteUserHandler handles soft-deleting a user (admin only)
//
//	@Summary		Delete a user
//	@Description	Soft-deletes a user by their unique identifier (admin only)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	deleteUserResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/users/{id} [delete]
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

		c.JSON(http.StatusOK, deleteUserResponse{Message: "user deleted"})
	}
}

// Request/Response types for Swagger

type listUsersResponse struct {
	Users []inbound.UserDTO `json:"users"`
}

type updateUserResponse struct {
	Message string `json:"message" example:"user updated"`
}

type deleteUserResponse struct {
	Message string `json:"message" example:"user deleted"`
}
