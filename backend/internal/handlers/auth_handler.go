package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/internal/auth"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/response"
	"github.com/notFil/cspotlight/internal/services"
	"github.com/notFil/cspotlight/internal/store"
	"go.uber.org/zap"
)

type AuthHandler struct {
	userService services.UserService
	cache       store.Cache
	tokenCfg    config.Token
}

func NewAuthHandler(userService services.UserService, cache store.Cache, tokenCfg config.Token) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		cache:       cache,
		tokenCfg:    tokenCfg,
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.AuthRequest  true  "Login credentials"
// @Success      200   {object}  map[string]interface{} "Login successful"
// @Failure      401   {object}  map[string]interface{} "Authentication failed"
// @Failure      500   {object}  map[string]interface{} "Internal server error"
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	var authRequest models.AuthRequest
	if err := c.BindJSON(&authRequest); err != nil {
		log.Warn("invalid login payload", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid login payload"))
		return
	}

	log.Info("attempting login", zap.String("username", authRequest.Username))

	user, err := h.userService.AuthenticateUser(ctx, &authRequest)
	if err != nil {
		log.Warn("authentication failed", zap.String("username", authRequest.Username), zap.Error(err))
		c.Error(apperrors.New(http.StatusUnauthorized, "authentication failed"))
		return
	}

	tokens, err := h.generateAndCacheTokens(ctx, user)
	if err != nil {
		log.Error("failed to generate and cache tokens", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "token retrieval failed"))
		return
	}

	log.Info("login successful", zap.String("user_id", user.ID.String()))
	response.Success(c, http.StatusOK, "login successful", tokens)
}

// Register godoc
// @Summary      Register user
// @Description  Registers a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.UserRegisterDTO  true  "User registration info"
// @Success      201   {object}  map[string]interface{} "User registered successfully"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to register user"
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	var user models.UserRegisterDTO
	if err := c.BindJSON(&user); err != nil {
		log.Warn("invalid register payload", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid registration payload"))
		return
	}

	log.Info("registering user", zap.String("email", user.Email))

	if err := h.userService.RegisterUser(ctx, &user); err != nil {
		log.Error("failed to register user", zap.String("email", user.Email), zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to register user"))
		return
	}

	log.Info("user registered successfully", zap.String("email", user.Email))
	response.Success(c, http.StatusCreated, "user registered successfully", nil)
}

// Refresh token godoc
// @Summary      Refresh token
// @Description  Refreshes the JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RefreshTokenRequest  true  "Refresh token request"
// @Success      200   {object}  map[string]interface{} "token refreshed successfully"
// @Failure      401   {object}  map[string]interface{} "invalid refresh token"
// @Failure      500   {object}  map[string]interface{} "failed to refresh token"
// @Router       /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	var req models.RefreshTokenRequest

	if err := c.BindJSON(&req); err != nil {
		log.Warn("invalid refresh token payload", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid refresh token payload"))
		return
	}

	claims, err := auth.ValidateRefreshToken(req.RefreshToken, h.tokenCfg.RefreshSecretKey)
	if err != nil {
		log.Warn("invalid refresh token", zap.Error(err))
		c.Error(apperrors.New(http.StatusUnauthorized, "invalid refresh token"))
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Error("invalid user id in claims", zap.String("id", claims.Subject), zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to get user"))
		return
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		log.Error("failed to get user for refresh token", zap.String("user_id", claims.Subject), zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to get user"))
		return
	}

	// Invalidate old refresh token
	if err := h.revokeTokens(ctx, user.ID, []string{req.RefreshToken}); err != nil {
		log.Error("failed to revoke old refresh token", zap.Error(err))
	}

	tokens, err := h.generateAndCacheTokens(ctx, user)
	if err != nil {
		log.Error("failed to generate and cache tokens", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to generate tokens"))
		return
	}

	log.Info("token refreshed", zap.String("user_id", user.ID.String()))
	response.Success(c, http.StatusOK, "token refreshed", tokens)
}

// Sign out godoc
// @Summary      Sign out
// @Description  Signs out the user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200   {object}  map[string]interface{} "signed out successfully"
// @Failure      500   {object}  map[string]interface{} "failed to sign out"
// @Router       /api/auth/signout [post]
func (h *AuthHandler) SignOut(c *gin.Context) {
	log := logger.FromContext(c)
	log.Info("user attempting to sign out")

	var req models.RefreshTokenRequest

	if err := c.BindJSON(&req); err != nil {
		log.Warn("invalid sign out payload", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid sign out payload"))
		return
	}

	authHeader := c.GetHeader("Authorization")
	tokenString, err := auth.ParseAuthHeader(authHeader)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()
	user := auth.GetUserContext(ctx)
	if user == nil {
		log.Warn("no user found in context")
		c.Error(apperrors.New(http.StatusUnauthorized, "no user found in context"))
		return
	}

	if err := h.revokeTokens(ctx, user.Subject, []string{tokenString, req.RefreshToken}); err != nil {
		log.Error("failed to revoke tokens during sign out", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to sign out"))
		return
	}

	response.Success(c, http.StatusOK, "signed out successfully", nil)
}

func (h *AuthHandler) generateAndCacheTokens(ctx context.Context, user *models.UserFetchDTO) (*auth.Token, error) {
	tokens, err := auth.GenerateJWT(user, h.tokenCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate jwt: %w", err)
	}

	accessTokenJTI, err := auth.GetJTI(tokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token jti: %w", err)
	}

	if err := h.cache.Set(ctx, fmt.Sprintf("user:%s:%s", user.ID.String(), accessTokenJTI), "", time.Duration(h.tokenCfg.ExpiryInMinutes)*time.Minute); err != nil {
		return nil, fmt.Errorf("failed to set access token jti in cache: %w", err)
	}

	refreshJTI, err := auth.GetJTI(tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh jti: %w", err)
	}

	if err := h.cache.Set(ctx, fmt.Sprintf("user:%s:%s", user.ID.String(), refreshJTI), "", time.Duration(h.tokenCfg.RefreshExpiryInMinutes)*time.Minute); err != nil {
		return nil, fmt.Errorf("failed to set refresh jti in cache: %w", err)
	}

	return tokens, nil
}

func (h *AuthHandler) revokeTokens(ctx context.Context, userID uuid.UUID, tokens []string) error {
	var errs []error
	for _, token := range tokens {
		jti, err := auth.GetJTI(token)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get jti from token: %w", err))
			continue
		}

		if err := h.cache.Del(ctx, fmt.Sprintf("user:%s:%s", userID.String(), jti)); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete jti from cache: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to revoke some tokens: %v", errs)
	}
	return nil
}
