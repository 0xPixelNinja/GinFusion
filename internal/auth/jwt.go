package auth

import (
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/0xPixelNinja/GinFusion/internal/config"
)

// Claims defines the structure for JWT claims.
type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for a given user ID.
func GenerateToken(cfg *config.Config, userID string) (string, error) {
    expirationTime := time.Now().Add(time.Second * time.Duration(cfg.JWT.Expiration))
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.JWT.Secret))
}

// ValidateToken validates a JWT token string.
func ValidateToken(cfg *config.Config, tokenStr string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(cfg.JWT.Secret), nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    return claims, nil
}
