package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/0xPixelNinja/GinFusion/internal/auth"
	"github.com/0xPixelNinja/GinFusion/internal/config"
	"github.com/0xPixelNinja/GinFusion/internal/handlers"
)

// RegisterAPIRoutes sets up routes for the user-facing API.
func RegisterAPIRoutes(r *gin.Engine) {
	// Public endpoints.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "API is healthy"})
	})
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	protected := r.Group("/")
	cfg := config.LoadConfig()
	protected.Use(auth.JWTAuthMiddleware(cfg))
	protected.POST("/apikey/generate", handlers.GenerateAPIKey)
	protected.GET("/apikey", handlers.GetUserAPIKey)
	protected.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.JSON(200, gin.H{"message": "Protected route", "user": userID})
	})
	protected.GET("/profile", handlers.GetProfile)
	protected.PUT("/profile", handlers.UpdateProfile)
	protected.POST("/refresh", handlers.RefreshToken)
	protected.GET("/activity", handlers.ActivityLog)
	protected.POST("/logout", handlers.Logout)
	protected.GET("/stats", handlers.GetUserStats)
}
