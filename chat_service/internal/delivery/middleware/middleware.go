package middleware

import (
	"chat_service/internal/domain"
	"strings"

	"github.com/gin-gonic/gin"
)

type MiddlewareRepo struct {
	repoToken domain.TokenManager
}

func NewMiddleware(token domain.TokenManager) *MiddlewareRepo {
	return &MiddlewareRepo{repoToken: token}
}

func (m *MiddlewareRepo) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			c.JSON(401, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}
		userID, name, role, err := m.repoToken.VerifyToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("name", name)

		c.Next()
	}

}

func (m *MiddlewareRepo) RequireAuthWS() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")

		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		userID, name, role, err := m.repoToken.VerifyToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("name", name)
		c.Set("role", role)

		c.Next()
	}
}
