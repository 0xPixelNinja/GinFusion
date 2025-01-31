package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/0xPixelNinja/GinFusion/internal/routes"
)

func main() {

    router := gin.Default()
    routes.RegisterAPIRoutes(router)

    // Start the server on a default port (adjust as needed)
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to run API server: %v", err)
    }
}
