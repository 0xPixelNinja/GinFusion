package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// RequireAPIKey ensures that an API key is provided as a query parameter.
func RequireAPIKey() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Query("apikey") == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
            return
        }
        c.Next()
    }
}
