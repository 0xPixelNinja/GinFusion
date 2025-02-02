package middleware

import (
    "context"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "github.com/0xPixelNinja/GinFusion/internal/repository"
)

// RateLimitMiddlewareWithAPIKey applies a dynamic rate limit based on the API key.
// The API key is expected as a query parameter "apikey".
func RateLimitMiddlewareWithAPIKey(rdb *redis.Client, defaultLimit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKeyQuery := c.Query("apikey")
        if apiKeyQuery == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
            return
        }
        var limit int
        keyData, err := repository.GetAPIKey(apiKeyQuery)
        if err == nil && keyData.RateLimit > 0 {
            limit = keyData.RateLimit
        } else {
            limit = defaultLimit
        }
        identifier := apiKeyQuery
        ctx := context.Background()
        redisKey := "rate:" + identifier

        count, err := rdb.Incr(ctx, redisKey).Result()
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
            return
        }
        if count == 1 {
            rdb.Expire(ctx, redisKey, window)
        }
        if count > int64(limit) {
            ttl, _ := rdb.TTL(ctx, redisKey).Result()
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded. Try again in " + strconv.Itoa(int(ttl.Seconds())) + " seconds.",
            })
            return
        }
        c.Next()
    }
}
