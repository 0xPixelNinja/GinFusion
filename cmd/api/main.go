package main

import (
	"log"

	"github.com/0xPixelNinja/GinFusion/internal/config"
	"github.com/0xPixelNinja/GinFusion/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()
	
    router := gin.Default()
    routes.RegisterAPIRoutes(router)

    // Start the server on a default port (adjust as needed)
    if err := router.Run(cfg.Server.Address); err != nil {
        log.Fatalf("Failed to run API server: %v", err)
    }
}
