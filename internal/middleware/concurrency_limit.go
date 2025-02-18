package middleware

import (
    "net/http"
    "sync"

    "github.com/gin-gonic/gin"
    "github.com/0xPixelNinja/GinFusion/internal/repository"
)

var (
    apiKeySemaphores = make(map[string]chan struct{})
    semMutex         sync.Mutex
)

func getSemaphore(identifier string, concurrency int) chan struct{} {
    semMutex.Lock()
    defer semMutex.Unlock()
    sem, exists := apiKeySemaphores[identifier]
    if !exists || cap(sem) != concurrency {
        sem = make(chan struct{}, concurrency)
        apiKeySemaphores[identifier] = sem
    }
    return sem
}

// ConcurrencyLimitMiddlewareWithAPIKey limits concurrent requests per API key.
// The API key is expected as a query parameter "apikey".
func ConcurrencyLimitMiddlewareWithAPIKey(defaultConcurrency int) gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.Query("apikey")
        if apiKey == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
            return
        }
        var concurrency int
        keyData, err := repository.GetAPIKey(apiKey)
        if err == nil && keyData.Concurrency > 0 {
            concurrency = keyData.Concurrency
        } else {
            concurrency = defaultConcurrency
        }
        identifier := apiKey
        sem := getSemaphore(identifier, concurrency)
        select {
        case sem <- struct{}{}:
            defer func() { <-sem }()
            c.Next()
        default:
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "Too many concurrent requests, please try again later.",
            })
        }
    }
}
