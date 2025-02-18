package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/0xPixelNinja/GinFusion/internal/repository"
)

// RecordStatsMiddleware records request statistics for endpoints that require an API key.
func RecordStatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start)

		apiKey := c.Query("apikey")
		if apiKey == "" {
			return
		}
		ctx := context.Background()

		repository.GetRedisClient().Incr(ctx, "global:requests")
		hashKey := "stats:" + apiKey
		repository.GetRedisClient().HIncrBy(ctx, hashKey, "total_requests", 1)
		repository.GetRedisClient().HIncrByFloat(ctx, hashKey, "total_response_time", float64(elapsed.Milliseconds()))
		repository.GetRedisClient().HSet(ctx, hashKey, "last_request", time.Now().Unix())
	}
}
