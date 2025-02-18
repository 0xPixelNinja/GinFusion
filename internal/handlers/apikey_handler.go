package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/0xPixelNinja/GinFusion/internal/models"
	"github.com/0xPixelNinja/GinFusion/internal/repository"
	"github.com/google/uuid"
)

// GenerateAPIKey generates a new API key for the authenticated user with default limits.
func GenerateAPIKey(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	uid := userID.(string)

	if existingKey, err := repository.GetAPIKeyByUser(uid); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API key already generated", "apikey": existingKey})
		return
	}

	const defaultRateLimit = 20
	const defaultConcurrency = 10

	key := uuid.New().String()
	apiKey := models.APIKey{
		Key:         key,
		UserID:      uid,
		RateLimit:   defaultRateLimit,
		Concurrency: defaultConcurrency,
		Created:     time.Now().Unix(),
	}

	if err := repository.CreateAPIKey(apiKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	repository.LogActivity(uid, "API key generated")
	c.JSON(http.StatusOK, gin.H{"message": "API key generated", "apikey": apiKey})
}

// GetUserAPIKey returns the current API key for the authenticated user.
func GetUserAPIKey(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	uid := userID.(string)
	apiKey, err := repository.GetAPIKeyByUser(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No API key generated for this user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"apikey": apiKey})
}
