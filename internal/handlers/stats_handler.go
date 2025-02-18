package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/0xPixelNinja/GinFusion/internal/repository"
)

// GetUserStats returns comprehensive API statistics for the current user's API key,
// including total requests, average response time (ms), remaining credits,
// global requests, and the key details.
func GetUserStats(c *gin.Context) {
	apiKey := c.Query("apikey")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API key is required"})
		return
	}

	ctx := context.Background()
	hashKey := "stats:" + apiKey
	statsMap, err := repository.GetRedisClient().HGetAll(ctx, hashKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stats"})
		return
	}

	totalRequests, _ := strconv.ParseInt(statsMap["total_requests"], 10, 64)
	totalResponseTime, _ := strconv.ParseFloat(statsMap["total_response_time"], 64)
	avgResponseTime := 0.0
	if totalRequests > 0 {
		avgResponseTime = totalResponseTime / float64(totalRequests)
	}

	rateKey := "rate:" + apiKey
	currentCountStr, err := repository.GetRedisClient().Get(ctx, rateKey).Result()
	currentCount := int64(0)
	if err == nil {
		currentCount, _ = strconv.ParseInt(currentCountStr, 10, 64)
	}

	apiKeyData, err := repository.GetAPIKey(apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve API key details"})
		return
	}
	remainingCredits := apiKeyData.RateLimit - int(currentCount)
	if remainingCredits < 0 {
		remainingCredits = 0
	}

	globalRequestsStr, err := repository.GetRedisClient().Get(ctx, "global:requests").Result()
	globalRequests := int64(0)
	if err == nil {
		globalRequests, _ = strconv.ParseInt(globalRequestsStr, 10, 64)
	}

	c.JSON(http.StatusOK, gin.H{
		"api_key_details": apiKeyData,
		"api_key_stats": gin.H{
			"total_requests":        totalRequests,
			"avg_response_time_ms":  avgResponseTime,
			"remaining_api_credits": remainingCredits,
		},
		"global_requests": globalRequests,
	})
}
