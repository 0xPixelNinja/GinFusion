package logger

import (
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

var Log = logrus.New()

// InitLogger configures Logrus.
func InitLogger() {
    Log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })
    Log.SetOutput(os.Stdout)
    Log.SetLevel(logrus.InfoLevel)
}

// GinLogger returns a Gin middleware that logs HTTP requests.
func GinLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        c.Next()
        duration := time.Since(start)
        Log.WithFields(logrus.Fields{
            "status":   c.Writer.Status(),
            "method":   c.Request.Method,
            "path":     path,
            "duration": duration,
            "client":   c.ClientIP(),
        }).Info("HTTP request")
    }
}
