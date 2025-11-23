package common

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/notFil/cspotlight/internal/models"
)

type AuthClaims struct {
	UserID string
	TeamID *string
	Role   string
}

type claims struct {
	Username string `json:"username"`
	TeamID   string `json:"team,omitempty"`
	Role     string `json:"role"`

	jwt.RegisteredClaims
}

func GenerateJWT(secretKey string, user *models.UserFetchDTO, expiryMinutes int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)
	claims := &claims{
		Username: user.Username,
		TeamID:   user.TeamID,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(secretKey string, tokenString string) (*claims, error) {
	claims := &claims{}
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

func (c *AuthClaims) IsAdmin() bool {
	if c.Role == "admin" || c.Role == "superadmin" {
		return true
	}
	return false
}

func (c *AuthClaims) IsSuperadmin() bool {
	return c.Role == "superadmin"
}

func GetUserClaims(c *gin.Context) AuthClaims {
	return c.MustGet(ClaimsContextKey).(AuthClaims)
}
