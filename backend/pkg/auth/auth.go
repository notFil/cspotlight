package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/notFil/cspotlight/internal/models"
)

type Claims struct {
	jwt.RegisteredClaims

	Username string `json:"username"`
	TeamID   string `json:"team,omitempty"`
	Role     string `json:"role"`
}

type Token struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	CreatedAt        string `json:"createdAt"`
	ExpiresIn        string `json:"expiresIn"`
	RefreshExpiresIn string `json:"refreshExpiresIn"`
}

func GenerateJWT(user *models.UserFetchDTO, cfg config.JWTConfig) (token *Token, err error) {
	currentTime := time.Now()
	expirationTime := currentTime.Add(time.Duration(cfg.ExpiryInMinutes) * time.Minute)
	refreshExpirationTime := currentTime.Add(time.Duration(cfg.RefreshExpiryInMinutes) * time.Minute)

	accessTokenClaims := Claims{
		Username: user.Username,
		TeamID:   user.TeamID,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
			Subject:   user.ID,
			ID:        uuid.New().String(),
		},
	}

	refreshTokenClaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
			Subject:   user.ID,
			ID:        uuid.New().String(),
		},
	}

	accessToken, err := generateJWT(&accessTokenClaims, cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateJWT(&refreshTokenClaims, cfg.RefreshSecretKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		CreatedAt:        currentTime.Format(time.RFC3339),
		ExpiresIn:        expirationTime.Sub(currentTime).String(),
		RefreshExpiresIn: refreshExpirationTime.Sub(currentTime).String(),
	}, nil
}

func ValidateAccessToken(tokenString string, secretKey string) (claims *Claims, err error) {
	return validateJWT(secretKey, tokenString)
}

func ValidateRefreshToken(tokenString string, secretKey string) (claims *Claims, err error) {
	return validateJWT(secretKey, tokenString)
}

func (c *Claims) IsAdmin() bool {
	return c.Role == "admin"
}

func (c *Claims) IsSuperadmin() bool {
	return c.Role == "superadmin"
}

func GetUserClaims(c *gin.Context) *Claims {
	claims, exists := c.Get(constants.ClaimsContextKey)
	if !exists {
		return nil
	}
	if claims, ok := claims.(*Claims); ok {
		return claims
	}
	return nil
}

func generateJWT(claims *Claims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateJWT(secretKey string, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
