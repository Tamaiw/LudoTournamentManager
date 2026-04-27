package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var initOnce = func() {
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("fallback-secret-for-dev-only")
	}
}

func init() {
	initOnce()
}

const (
	UserIDKey   = "user_id"
	UserRoleKey = "user_role"
)

// AuthMiddleware validates JWT tokens and sets user info in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "missing token",
				},
			})
			c.Abort()
			return
		}

		userID, role, err := validateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "invalid token",
				},
			})
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Set(UserRoleKey, role)
		c.Next()
	}
}

// ExtractToken extracts the JWT token from the Authorization header
func ExtractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// validateJWT validates a JWT token and returns the user ID and role
func validateJWT(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	userID, _ := claims["user_id"].(string)
	role, _ := claims["role"].(string)

	return userID, role, nil
}
