package routes

import "github.com/gin-gonic/gin"

// RegisterAPIRoutes sets up routes for the user-facing API.
func RegisterAPIRoutes(r *gin.Engine) {
    // Define a simple health-check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "API is healthy"})
    })

    // More API routes (registration, login, etc.) will be added later.
}

// RegisterAdminRoutes sets up routes for the admin interface.
func RegisterAdminRoutes(r *gin.Engine) {
    // Define a simple health-check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "Admin is healthy"})
    })

    // More admin routes (user management, API key management, etc.) will be added later.
}
