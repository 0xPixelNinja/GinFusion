package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/0xPixelNinja/GinFusion/internal/handlers"
)

// RegisterAdminRoutes sets up routes for the admin server.
// These endpoints allow administrators to manage users, API keys, and view usage statistics.
func RegisterAdminRoutes(r *gin.Engine) {

	admin_ := r.Group("/admin")

	admin_.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Admin is healthy"})
	})

	// User management endpoints.
	admin_.GET("/users", handlers.ListUsers)
	admin_.GET("/stats", handlers.Stats)
	admin_.GET("/apikeys", handlers.ListAPIKeys)
	admin_.POST("/apikeys", handlers.CreateAPIKey)
	admin_.PUT("/apikeys/:key", handlers.UpdateAPIKey)
	admin_.DELETE("/apikeys/:key", handlers.DeleteAPIKey)
}
