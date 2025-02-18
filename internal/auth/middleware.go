package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/0xPixelNinja/GinFusion/internal/config"
	"github.com/0xPixelNinja/GinFusion/internal/repository"
)

// JWTAuthMiddleware enforces JWT validation on protected endpoints.
func JWTAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		tokenStr := parts[1]

		if repository.IsTokenBlacklisted(tokenStr) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is invalidated. Please log in again."})
			return
		}

		claims, err := ValidateToken(cfg, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
