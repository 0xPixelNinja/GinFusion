package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/0xPixelNinja/GinFusion/internal/models"
    "github.com/0xPixelNinja/GinFusion/internal/repository"
    "github.com/google/uuid"
)

// ListAPIKeys returns a list of all API keys.
func ListAPIKeys(c *gin.Context) {
	apiKeys, err := repository.ListAPIKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve API keys"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"apikeys": apiKeys})
}

// ListUsers returns all registered users.
func ListUsers(c *gin.Context) {
    users, err := repository.ListUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"users": users})
}

// Stats returns usage statistics.
func Stats(c *gin.Context) {
    stats, err := repository.GetUsageStats()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stats"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// APIKeyUpdateRequest represents the payload for updating an API key.
type APIKeyUpdateRequest struct {
    UserID      string `json:"user_id" binding:"required"`
    RateLimit   int    `json:"rate_limit" binding:"required"`
    Concurrency int    `json:"concurrency" binding:"required"`
}

// UpdateAPIKey updates an API key's settings.
func UpdateAPIKey(c *gin.Context) {
    key := c.Param("key")
    var req APIKeyUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    apiKey, err := repository.GetAPIKey(key)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
        return
    }

    apiKey.UserID = req.UserID
    apiKey.RateLimit = req.RateLimit
    apiKey.Concurrency = req.Concurrency

    if err := repository.UpdateAPIKey(apiKey); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API key"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "API key updated", "apikey": apiKey})
}

// CreateAPIKeyRequest represents the payload for creating an API key.
type CreateAPIKeyRequest struct {
    UserID      string `json:"user_id" binding:"required"`
    RateLimit   int    `json:"rate_limit" binding:"required"`
    Concurrency int    `json:"concurrency" binding:"required"`
}

// CreateAPIKey creates a new API key.
func CreateAPIKey(c *gin.Context) {
    var req CreateAPIKeyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    key := uuid.New().String()
    apiKey := models.APIKey{
        Key:         key,
        UserID:      req.UserID,
        RateLimit:   req.RateLimit,
        Concurrency: req.Concurrency,
        Created:     time.Now().Unix(),
    }
    if err := repository.CreateAPIKey(apiKey); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "API key created", "apikey": apiKey})
}

// DeleteAPIKey deletes an API key.
func DeleteAPIKey(c *gin.Context) {
    key := c.Param("key")
    if err := repository.DeleteAPIKey(key); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "API key deleted"})
}
