package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/0xPixelNinja/GinFusion/internal/routes"
)

func main() {
	
    router := gin.Default()
    routes.RegisterAdminRoutes(router)

    // Start the admin server on a different port (adjust as needed)
    if err := router.Run(":8081"); err != nil {
        log.Fatalf("Failed to run admin server: %v", err)
    }
}
