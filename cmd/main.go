package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xPixelNinja/GinFusion/internal/config"
	"github.com/0xPixelNinja/GinFusion/internal/repository"
	"github.com/0xPixelNinja/GinFusion/internal/routes"
	"github.com/0xPixelNinja/GinFusion/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {

	
	logger.InitLogger()

	cfg := config.LoadConfig()
	repository.InitRedis(cfg)
	repository.InitSQLiteConnections()

	gin.SetMode(gin.ReleaseMode)
	apiRouter := gin.New()

	cors_ := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	// Create the API router.
	apiRouter.Use(gin.Recovery(), logger.GinLogger())
	apiRouter.Use(cors_)
	routes.RegisterAPIRoutes(apiRouter)
	routes.RegisterModuleRoutes(apiRouter)

	// Create the Admin router.
	adminRouter := gin.New()
	adminRouter.Use(gin.Recovery(), logger.GinLogger())
	adminRouter.Use(cors_)
	routes.RegisterAdminRoutes(adminRouter)

	// Define HTTP servers for API and Admin.
	apiServer := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: apiRouter,
	}
	adminServer := &http.Server{
		Addr:    cfg.Admin.Address,
		Handler: adminRouter,
	}

	// Start the API server in a new goroutine.
	go func() {
		log.Printf("Starting API server on %s", cfg.Server.Address)
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API server failed: %v", err)
		}
	}()

	go func() {
		log.Printf("Starting Admin server on %s", cfg.Admin.Address)
		if err := adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Admin server failed: %v", err)
		}
	}()

	// Wait for OS interrupt signal to gracefully shutdown servers.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		log.Fatalf("API Server forced to shutdown: %v", err)
	}
	if err := adminServer.Shutdown(ctx); err != nil {
		log.Fatalf("Admin Server forced to shutdown: %v", err)
	}

	log.Println("Servers exited gracefully")
}