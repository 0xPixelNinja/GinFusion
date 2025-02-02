package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/0xPixelNinja/GinFusion/internal/handlers"
)

// RegisterAdminRoutes sets up routes for the admin server.
// These endpoints allow administrators to manage users, API keys, and view usage statistics.
func RegisterAdminRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Admin is healthy"})
	})

	// User management endpoints.
	r.GET("/users", handlers.ListUsers)
	r.GET("/stats", handlers.Stats)
	r.GET("/apikeys", handlers.ListAPIKeys)
	r.POST("/apikeys", handlers.CreateAPIKey)
	r.PUT("/apikeys/:key", handlers.UpdateAPIKey)
	r.DELETE("/apikeys/:key", handlers.DeleteAPIKey)
}
