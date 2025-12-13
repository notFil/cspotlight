package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/internal/constants"
	apperrors "github.com/notFil/cspotlight/internal/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/notFil/cspotlight/internal/models"
)

type contextKey struct{}

type AuthContext struct {
	Subject uuid.UUID
	TeamID  uuid.UUID
	Role    string
}

type claims struct {
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

func GenerateJWT(user *models.UserFetchDTO, cfg config.Token) (token *Token, err error) {
	currentTime := time.Now()
	expirationTime := currentTime.Add(time.Duration(cfg.ExpiryInMinutes) * time.Minute)
	refreshExpirationTime := currentTime.Add(time.Duration(cfg.RefreshExpiryInMinutes) * time.Minute)

	accessTokenClaims := claims{
		Username: user.Username,
		TeamID:   user.TeamID.String(),
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
			Subject:   user.ID.String(),
			ID:        uuid.New().String(),
		},
	}

	refreshTokenClaims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
			Subject:   user.ID.String(),
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

func ValidateAccessToken(tokenString string, secretKey string) (claims *claims, err error) {
	return validateJWT(secretKey, tokenString)
}

func ValidateRefreshToken(tokenString string, secretKey string) (claims *claims, err error) {
	return validateJWT(secretKey, tokenString)
}

func ParseAuthHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", apperrors.New(http.StatusUnauthorized, "missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", apperrors.New(http.StatusUnauthorized, "invalid authorization header")
	}

	tokenString := parts[1]

	return tokenString, nil
}

// Used to retrieve the JTI from a token without validating it. Validation
// is handled in the middleware
func GetJTI(tokenString string) (string, error) {
	claims := &claims{}
	_, _, err := jwt.NewParser().ParseUnverified(tokenString, claims)
	if err != nil {
		return "", err
	}
	return claims.ID, nil
}

func (c *AuthContext) IsAdmin() bool {
	return c.Role == "admin"
}

func (c *AuthContext) IsSuperadmin() bool {
	return c.Role == "superadmin"
}

func (c *AuthContext) SameUser(id uuid.UUID) bool {
	return c.Subject == id
}

func (c *AuthContext) SameTeam(id uuid.UUID) bool {
	return c.TeamID == id
}

func ContextWithUser(ctx context.Context, claims *AuthContext) context.Context {
	return context.WithValue(ctx, contextKey{}, claims)
}

func GetUserContext(ctx context.Context) *AuthContext {
	if claims, ok := ctx.Value(contextKey{}).(*AuthContext); ok {
		return claims
	}

	if val := ctx.Value(constants.ClaimsContextKey); val != nil {
		if claims, ok := val.(*AuthContext); ok {
			return claims
		}
	}

	return nil
}

func generateJWT(claims *claims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateJWT(secretKey string, tokenString string) (*claims, error) {
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
