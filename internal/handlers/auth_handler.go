package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "golang.org/x/crypto/bcrypt"

    "github.com/0xPixelNinja/GinFusion/internal/auth"
    "github.com/0xPixelNinja/GinFusion/internal/config"
    "github.com/0xPixelNinja/GinFusion/internal/models"
    "github.com/0xPixelNinja/GinFusion/internal/repository"
)

var validate = validator.New()

// RegistrationRequest represents the payload for user registration.
type RegistrationRequest struct {
    Username string `json:"username" binding:"required,min=3"`
    Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the payload for user login.
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// Register creates a new user.
func Register(c *gin.Context) {
    var req RegistrationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := validate.Struct(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	if _, err := repository.GetUser(req.Username); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
        return
    }

    user := models.User{
        ID:       req.Username,
        Username: req.Username,
        Password: string(hashedPassword),
    }
    if err := repository.CreateUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login authenticates a user and returns a JWT token.
func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := validate.Struct(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := repository.GetUser(req.Username)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    cfg := config.LoadConfig()
    token, err := auth.GenerateToken(cfg, user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
    repository.LogActivity(user.ID, "User Logged In")
    c.JSON(http.StatusOK, gin.H{"token": token})
}
