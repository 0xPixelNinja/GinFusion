package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/0xPixelNinja/GinFusion/internal/routes"
	"github.com/0xPixelNinja/GinFusion/internal/config"

)

func main() {

	cfg := config.LoadConfig()

    router := gin.Default()
    routes.RegisterAdminRoutes(router)

    // Start the admin server on a different port (adjust as needed)
    if err := router.Run(cfg.Admin.Address); err != nil {
        log.Fatalf("Failed to run admin server: %v", err)
    }
}
