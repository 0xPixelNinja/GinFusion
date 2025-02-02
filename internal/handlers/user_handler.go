package handlers

import (
    "net/http"
    "time"
    "strings"
    "github.com/golang-jwt/jwt/v4"

    "github.com/gin-gonic/gin"
    "github.com/0xPixelNinja/GinFusion/internal/auth"
    "github.com/0xPixelNinja/GinFusion/internal/config"
    "github.com/0xPixelNinja/GinFusion/internal/models"
    "github.com/0xPixelNinja/GinFusion/internal/repository"
)

// GetProfile returns the profile of the currently authenticated user.
func GetProfile(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }
    user, err := repository.GetUser(userID.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
        return
    }
    user.Password = ""
    c.JSON(http.StatusOK, gin.H{"profile": user})
}

// UpdateProfile allows a user to update their profile.
func UpdateProfile(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }
    var req models.User
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, err := repository.GetUser(userID.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
        return
    }
    user.Username = req.Username
    if err := repository.UpdateUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

// RefreshToken generates a new JWT token.
func RefreshToken(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }
    cfg := config.LoadConfig()
    token, err := auth.GenerateToken(cfg, userID.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout invalidates the current session token by blacklisting it.
func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	tokenStr := parts[1]

	cfg := config.LoadConfig()
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.ExpiresAt == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	if err := repository.AddTokenToBlacklist(tokenStr, ttl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}
	userID, _ := c.Get("userID")
	repository.LogActivity(userID.(string), "Logged out")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func ActivityLog(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	logs, err := repository.GetActivityLog(userID.(string), 50) // Retrieve up to 50 recent logs.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity log"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"activity": logs})
}